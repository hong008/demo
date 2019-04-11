package __go_lock

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type Data []byte

type DataFile interface {
	Read() (rsn int64, d Data, err error) //读取数据块
	Write(d Data) (wsn int64, err error)  //写数据块
	RSN() int64                           //获取最后读取的数据块的序列号
	WSN() int64                           //获取最后写入的数据块的序列号
	DataLen() uint32                      //获取数据块长度
	Close() error                         //关闭文件
}

type myDataFile struct {
	f       *os.File
	fmutex  sync.RWMutex //读写锁
	woffset int64        //写偏移量
	roffset int64        //读偏移量
	/*wmutex  sync.Mutex   //读文件的互斥锁
	rmutex  sync.Mutex   //写文件的互斥锁*/
	dataLen uint32
}

func NewDataFile(path string, len uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if len <= 0 {
		return nil, errors.New("Invalid data length!")
	}
	return &myDataFile{
		f:       f,
		dataLen: len,
	}, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	//先获取当前的读偏移量, 然后更新
	/*var offSet int64
	df.rmutex.Lock()
	offSet = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()*/

	var offSet int64
	for {
		offSet = atomic.LoadInt64(&df.roffset)
		if atomic.CompareAndSwapInt64(&df.roffset, offSet, offSet+int64(df.dataLen)) {
			break
		}
	}

	//根据当前获取的读偏移量读取数据快
	rsn = offSet / int64(df.dataLen) //数据序列号
	bytes := make([]byte, df.dataLen)

	for {
		df.fmutex.RLock()

		_, err = df.f.ReadAt(bytes, offSet)
		if err == io.EOF {
			df.fmutex.RUnlock()
			continue
		}

		if err != nil {
			df.fmutex.RUnlock()
			return
		}

		d = bytes
		df.fmutex.RUnlock()
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	//获取写的偏移量
	var offSet int64

	/*df.wmutex.Lock()
	offSet = df.woffset
	df.woffset += int64(df.dataLen)
	df.wmutex.Unlock()*/

	for {
		offSet = atomic.LoadInt64(&df.woffset)
		if atomic.CompareAndSwapInt64(&df.woffset, offSet, offSet+int64(df.dataLen)) {
			break
		}
	}

	//写入数据
	wsn = offSet / int64(df.dataLen)

	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:int(df.dataLen)]
	} else {
		bytes = d
	}

	df.fmutex.Lock()
	defer df.fmutex.Unlock()

	_, err = df.f.Write(bytes)
	return
}

func (df *myDataFile) RSN() int64 {
	/*df.rmutex.Lock()
	rsn := df.roffset / int64(df.dataLen)
	df.rmutex.Unlock()*/
	return atomic.LoadInt64(&df.roffset) / int64(df.dataLen)
}

func (df *myDataFile) WSN() int64 {
	/*df.wmutex.Lock()
	wsn := df.woffset / int64(df.dataLen)
	df.wmutex.Unlock()*/
	return atomic.LoadInt64(&df.woffset) / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	return df.f.Close()
}
