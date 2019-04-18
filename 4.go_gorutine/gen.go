package __go_gorutine

import (
	"bytes"
	"context"
	"demo/4.go_gorutine/lib"
	log "demo/5.go_log"
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

var logger = log.Dlogger()

type myGenerator struct {
	caller      lib.Caller           //调度器
	timeoutNs   time.Duration        //处理超时时间: 纳秒
	lps         uint32               //每秒载荷量
	durationNS  time.Duration        //负载持续时间，单位纳秒
	concurrency uint32               //载荷并发量
	tickets     lib.GoTickets        //Goroutine票池
	ctx         context.Context      //上下文
	cancelFunc  context.CancelFunc   //取消函数
	callCount   int64                //调用计数
	status      uint32               //状态
	resultCh    chan *lib.CallResult //调用结果通道
}

func NewGenerator(pset ParamSet) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNs:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}

//初始化工作
func (gen *myGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("Initializing the load generator...")
	var total64 = int64(gen.timeoutNs)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets
	buf.WriteString(fmt.Sprintf("Done. (concurrency = %d)", gen.concurrency))
	logger.Infoln(buf.String())
	return nil
}

//向载荷承受方发起一次调用
func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
	atomic.AddInt64(&gen.callCount, 1)
	if rawReq == nil {
		return &lib.RawResp{
			ID:  -1,
			Err: errors.New("Invalid raw request."),
		}
	}
	start := time.Now().UnixNano()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNs)
	end := time.Now().UnixNano()
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime,
		}
	} else {
		rawResp = lib.RawResp{
			ID:     rawReq.ID,
			Resp:   resp,
			Elapse: elapsedTime,
		}
	}
	return &rawResp
}

//发送调用结果
func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
	if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}
	select {
	case gen.resultCh <- result:
		return true
	default:
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

//打印被忽略的结果
func (gen *myGenerator) printIgnoredResult(result *lib.CallResult, cause string) {
	resultMsg := fmt.Sprintf("ID=%d, Code=%d, Msg=%s, Elapse=%v", result.ID, result.Code, result.Msg, result.Elapse)
	logger.Warnf("Ignored result: %s. (cause: %s)\n", resultMsg, cause)
}

//为停止载荷器做准备
func (gen *myGenerator) prepareToStop(ctxError error) {
	logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
	atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
	logger.Infof("Closing result channel...")
	close(gen.resultCh)
	atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)
}

//产生载荷并向承受方发送
func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
		default:
		}
		gen.asyncCall()
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
	}
}

//异步调用承受方的接口
func (gen *myGenerator) asyncCall() {
	gen.tickets.Take()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
				} else {
					errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
				}
				logger.Errorln(errMsg)
				result := &lib.CallResult{
					ID:   -1,
					Code: lib.RET_CODE_FATAL_CALL,
					Msg:  errMsg,
				}
				gen.sendResult(result)
			}
			gen.tickets.Return()
		}()

		rawReq := gen.caller.BuildReq()
		//调用状态： 0-未调用或者未调用中; 1-调用完成； 2-调用超时;
		var callStatus uint32
		timer := time.AfterFunc(gen.timeoutNs, func() {
			if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				return
			}
			result := &lib.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNs),
				Elapse: gen.timeoutNs,
			}
			gen.sendResult(result)
		})

		rawResp := gen.callOne(&rawReq)
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}
		timer.Stop()
		var result *lib.CallResult
		if rawResp.Err != nil {
			result = &lib.CallResult{
				ID:     rawResp.ID,
				Req:    rawReq,
				Code:   lib.RET_CODE_ERROR_CALL,
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse,
			}
		} else {
			result = gen.caller.CheckResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		gen.sendResult(result)
	}()
}

//启动载荷器
func (gen *myGenerator) Start() bool {
	logger.Infoln("Starting load generator...")
	if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
		if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
			return false
		}
	}

	//设定节流阀
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		logger.Infof("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval)
	}

	//初始化上下文和取消函数
	gen.ctx, gen.cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)

	//初始化计数
	gen.callCount = 0

	//设置状态为已启动
	atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)

	go func() {
		logger.Infoln("Genertating loads...")
		gen.genLoad(throttle)
		logger.Infof("Stopped. (call count: %d)", gen.callCount)
	}()
	return true
}

func (gen *myGenerator) Stop() bool {
	if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) {
		return false
	}
	gen.cancelFunc()
	for {
		if atomic.LoadUint32(&gen.status) == lib.STATUS_STOPPED {
			break
		}
		time.Sleep(time.Microsecond)
	}
	return true
}

func (gen *myGenerator) Status() uint32 {
	return atomic.LoadUint32(&gen.status)
}

func (gen *myGenerator) CallCount() int64 {
	return atomic.LoadInt64(&gen.callCount)
}
