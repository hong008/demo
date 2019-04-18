package module

import "net/http"

//用于汇集组件内部计数的类型
type Counts struct {
	CalledCount    uint64 //调用计数
	AcceptedCount  uint64 //接受请求的计数
	CompletedCount uint64 //成功完成的计数
	HandlingNumber uint64 //实时处理数
}

//组件摘要结构
type SummaryStruct struct {
	ID        MID         `json:"id"`
	Called    uint64      `json:"called"`
	Accepted  uint64      `json:"accepted"`
	Completed uint64      `json:"completed"`
	Handling  uint64      `json:"handling"`
	Extra     interface{} `json:"extra, omitempty"`
}

//组件基本的API
type Module interface {
	ID() MID                         //获取当前组件的ID
	Addr() string                    //获取当前组件的网络地址
	Score() uint64                   //获取当前组件的评分
	SetScore(score uint64)           //设置当前组件的评分
	ScoreCalculator() CalculateScore //获取当前组件的评分计算器
	CalledCount() uint64             //获取当前组件被调用的计数
	AcceptedCount() uint64           //获取当前组件接受的调用次数：可以判断组件是否因为超负荷而拒绝调用
	CompletedCount() uint64          //获取当前组件已完成的调用次数
	HandlingNumber() uint64          //获取当前组件正在处理的调用数量
	Counts() Counts                  //一次性获取所有计数
	Summary() SummaryStruct          //获取当前组件的摘要
}

//下载器API
type Downloader interface {
	Module
	Download(req *Request) (*Response, error)
}

//用于解析http响应的函数
type ParseResponse func(httpResp *http.Response, respDepth uint32) ([]Data, []error)

//分析器API
type Analyzer interface {
	Module
	RespParsers() []ParseResponse             //获取分析器所用的响应解析器的列表
	Analyze(resp *Response) ([]Data, []error) //根据规则分析响应并返回分析后的结果：下一级的请求和这一级分析结果对应的数据条目
}

//处理结果条目的函数
type ProcessItem func(item Item) (result Item, err error)

//条目处理管道API
type Pipeline interface {
	Module
	ItemProcessors() []ProcessItem //处理条目的函数列表
	Send(item Item) []error        //向条目处理管道发送条目
	FailFast() bool                //当前条目处理管道是否是快速结束的， 只要发生错误就定义为快速出错
	SetFailFast(failFast bool)     //设置是否快速失败
}
