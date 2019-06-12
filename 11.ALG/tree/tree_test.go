package tree

import (
	"fmt"
	"testing"
)

var (
	root = NewTree(0)
	t1   = NewTree(1)
	t2   = NewTree(2)
	t3   = NewTree(3)
	t4   = NewTree(4)
	t5   = NewTree(5)
	t6   = NewTree(6)
	t7   = NewTree(7)
	t8   = NewTree(8)
	t9   = NewTree(9)
)

func TestBinaryTreeNode_MidTraverseRecursion(t *testing.T) {
	root.Left = t1
	root.Right = t2

	t1.Left = t3
	t1.Right = t4

	t3.Left = t7
	t7.Left = t9

	t4.Left = t8

	t2.Left = t5
	t5.Right = t6

	fmt.Println(root.PostTraverse())
}
