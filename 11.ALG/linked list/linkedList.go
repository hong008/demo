package linked_list

import "fmt"

//表节点
type LinkNode struct {
	Data interface{} //数据
	Next *LinkNode   //指向下一个节点
}

//表头信息
type HeadInfo struct {
	Length int       //链表长度
	Head   *LinkNode //表头指针
}

//表尾追加节点
func AppendNode(headInfo *HeadInfo, nodeData interface{}) {
	if headInfo == nil || headInfo.Head == nil {
		panic("unclear head info")
	}

	tempNode := headInfo.Head
	for tempNode.Next != nil {
		tempNode = tempNode.Next
	}

	//尾节点next指向nil
	node := &LinkNode{
		Data: nodeData,
		Next: nil,
	}

	//追加前的尾节点的next指向新的尾节点
	tempNode.Next = node
	//表头信息+1
	headInfo.Length += 1
}

//新节点插入在指定位置的后面
func AddNode(headInfo *HeadInfo, insertIndex int, nodeData interface{}) {
	if headInfo == nil || headInfo.Head == nil {
		panic("unclear head info")
	}

	//判断插入位置是否合法，索引从0开始，所以插入位置最大不能超过length-1
	if insertIndex < 0 || insertIndex >= headInfo.Length {
		panic("illegal insert position")
	}

	tempNode := headInfo.Head
	for i := 0; i < insertIndex; i++ {
		tempNode = tempNode.Next
	}

	node := &LinkNode{
		Data: nodeData,
		Next: tempNode.Next,
	}

	tempNode.Next = node
	headInfo.Length += 1
}

//获取指定位置的节点
func GetNode(headInfo *HeadInfo, index int) *LinkNode {
	if headInfo == nil || headInfo.Head == nil {
		panic("unclear head info")
	}

	if index < 0 || index >= headInfo.Length {
		panic("illegal insert position")
	}

	tempNode := headInfo.Head
	count := 0
	for tempNode.Next != nil {
		tempNode = tempNode.Next
		count++
		if count == index {
			break
		}
	}

	return tempNode
}

//删除指定位置的节点
func DeleteNode(headInfo *HeadInfo, index int) {
	if headInfo == nil || headInfo.Head == nil {
		panic("unclear head info")
	}

	if index < 0 || index >= headInfo.Length {
		panic("illegal insert position")
	}

	tempNode := headInfo.Head
	count := 0
	for tempNode.Next != nil {
		tempNode = tempNode.Next
		count++
		//找到需要被删除节点的前一个节点
		if count+1 == index {
			break
		}
	}
	tempNode.Next = tempNode.Next.Next
	headInfo.Length -= 1
}

//获取节点数
func GetNodeCount(headInfo *HeadInfo) int {
	if headInfo == nil || headInfo.Head == nil {
		panic("unclear head info")
	}
	tempNode := headInfo.Head
	count := 0
	for tempNode != nil {
		tempNode = tempNode.Next
		count++
	}
	return count
}

//遍历链表
func ListLink(headInfo *HeadInfo) (result string) {
	if headInfo == nil || headInfo.Head == nil {
		panic("unclear head info")
	}
	result += fmt.Sprintf("{Length=%+v; Head=%p}   ", headInfo.Length, headInfo.Head)
	tempNode := headInfo.Head
	for tempNode != nil {
		result += fmt.Sprintf("{Data=%+v; Next=%p}   ", tempNode.Data, tempNode.Next)
		tempNode = tempNode.Next
	}
	return
}
