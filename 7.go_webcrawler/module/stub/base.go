package stub

import "demo/7.go_webcrawler/module"

type ModuleInternal interface {
	module.Module
	IncrCalledCount()    //增加调用计数
	IncrAcceptedCount()   //增加接收计数
	IncrCompletedCount() //增加完成得计数
	IncrHandlingNumber() //增加实时处理计数
	DecrHandlingNumber() //减少实时处理计数
	Clear()              //清空所有计数
}


