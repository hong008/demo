package tree

import (
	"container/list"
	"fmt"
)

//二叉树
type BinaryTreeNode struct {
	Data  interface{}     //数据
	Left  *BinaryTreeNode //左子树
	Right *BinaryTreeNode //右子树
}

func NewTree(data interface{}) *BinaryTreeNode {
	return &BinaryTreeNode{
		Data: data,
	}
}

func (tree *BinaryTreeNode) String() string {
	return fmt.Sprintf(" %v", tree.Data)
}

/*二叉树遍历（递归）*/
//前序遍历:以当前节点为根节点，根——>左——>右
func (tree *BinaryTreeNode) PreTraverseRecursion() (treeString string) {

	if tree == nil {
		return
	}

	treeString += tree.String()

	if tree.Left != nil {
		treeString += tree.Left.PreTraverseRecursion()
	}

	if tree.Right != nil {
		treeString += tree.Right.PreTraverseRecursion()
	}
	//fmt.Println(fmt.Sprintf("treeString = [%v]", treeString))
	return
}

//中序遍历:以当前节点为根节点，左——>根——>右
func (tree *BinaryTreeNode) MidTraverseRecursion() (treeString string) {
	if tree == nil {
		return
	}

	if tree.Left != nil {
		treeString += tree.Left.MidTraverseRecursion()
	}

	treeString += tree.String()

	if tree.Right != nil {
		treeString += tree.Right.MidTraverseRecursion()
	}
	//fmt.Println(fmt.Sprintf("treeString = [%v]", treeString))
	return
}

//后续遍历：以当前节点为根节点，左——>右——>根
func (tree *BinaryTreeNode) PostTraverseRecursion() (treeString string) {
	if tree == nil {
		return
	}

	if tree.Left != nil {
		treeString += tree.Left.PostTraverseRecursion()
	}

	if tree.Right != nil {
		treeString += tree.Right.PostTraverseRecursion()
	}

	treeString += tree.String()

	//fmt.Println(fmt.Sprintf("treeString = [%v]", treeString))
	return
}

/*非递归遍历：利用栈结构*/
//栈结构
type Stack struct {
	*list.List
}

//出栈
func (s *Stack) Pop() interface{} {
	if s == nil || s.Len() <= 0 {
		return nil
	}
	value := s.Back()
	s.Remove(value)
	return value.Value
}

//进栈
func (s *Stack) Push(d interface{}) {
	if s == nil {
		return
	}
	s.PushBack(d)
}

//获取栈顶元素
func (s *Stack) Top() interface{} {
	if s == nil {
		return nil
	}
	return s.Back().Value
}

/*	非递归前序遍历
	根——>左——>右
	1、从根节点开始访问，每访问一个元素，执行入栈操作并输出当前节点
	2、访问到最左边的子节点时，开始出栈
	3、每出栈一个元素需要该节点是否存在右节点，如果存在则重复操作1
*/
func (tree *BinaryTreeNode) PreTraverse() (result string) {
	if tree == nil {
		return
	}
	stack := &Stack{
		List: list.New(),
	}

	node := tree
	for node != nil || stack.Len() > 0 {
		if node != nil {
			stack.Push(node)
			result += node.String()
			node = node.Left
		} else {
			node = stack.Pop().(*BinaryTreeNode)
			node = node.Right
		}
	}
	return
}

/*	非递归中序遍历
	左——>根——>右
	1、从根节点开始遍历到最左边的子节点，每访问一个节点就入栈（此处用node访问每个节点）
	2、访问到最左边的子节点时开始出栈，出栈时做输出操作
	3、每次出栈一个元素，需要判断该元素是否存在右节点，如果存在，则重复步骤1
*/
func (tree *BinaryTreeNode) MidTraverse() (result string) {
	if tree == nil {
		return
	}

	stack := &Stack{
		List: list.New(),
	}

	node := tree
	for node != nil || stack.Len() > 0 {
		if node != nil {
			stack.Push(node)
			node = node.Left
		} else {
			node = stack.Pop().(*BinaryTreeNode)
			result += node.String()
			node = node.Right
		}
	}
	return
}

/*	非递归后续遍历
	左——>右——>根
	1、从根节点开始遍历到最左边的子节点，每访问一个节点就入栈（此处用node访问每个节点）
	2、最后一个左子节点入栈后开始出栈操作，出栈时做输出操作
	3、出栈条件：栈顶元素的右子节点为空或者右子节点已经出栈（此处用top纪录当前栈顶元素，last纪录最后出栈的元素）
	4、如果栈顶元素的右子节点不为空且未出栈，则继续步骤1
	为什么要纪录最后出站的元素？
	如果一个节点同时存在左右子节点，按照后序遍历的规则，最后一个出栈元素为一定为该节点的右子节点，此时该节点的子节点已经遍历完，需要将该节点出栈并输出
*/
func (tree *BinaryTreeNode) PostTraverse() (result string) {
	if tree == nil {
		return
	}

	stack := &Stack{
		List: list.New(),
	}

	node := tree
	var topNode, lastNode *BinaryTreeNode //top为栈顶元素、last为最后出栈的元素

	for node != nil || stack.Len() > 0 {
		if node != nil {
			stack.Push(node)
			node = node.Left
		} else {
			topNode = stack.Top().(*BinaryTreeNode)
			if topNode.Right == nil || topNode.Right == lastNode {
				stack.Pop()
				result += topNode.String()
				lastNode = topNode
			} else {
				node = topNode.Right
			}
		}
	}
	return
}

//广度优先遍历
func (tree *BinaryTreeNode) LevelTraverse() (result string) {
	if tree == nil {
		return
	}

	treeList := list.New()
	treeList.PushBack(tree)

	for treeList.Len() > 0 {
		element := treeList.Front()
		node := element.Value.(*BinaryTreeNode)

		result += node.String()
		treeList.Remove(element)

		if node.Left != nil {
			treeList.PushBack(node.Left)
		}
		if node.Right != nil {
			treeList.PushBack(node.Right)
		}
	}
	return
}
