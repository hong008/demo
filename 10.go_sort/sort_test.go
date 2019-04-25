package _0_go_sort

import (
	"fmt"
	"testing"
)

var num  = []int{100, 229, 90, 33, 77, 99, 21, 64, 72, 100, 90, 99, 300}


func TestBubbleSort(t *testing.T) {
	fmt.Println(BubbleSort(num))
}

func TestQuickSort(t *testing.T) {
	fmt.Println(QuickSort(num))
}

func TestSelectSort(t *testing.T) {
	fmt.Println(SelectSort(num))
}

func TestInsertSort(t *testing.T) {
	fmt.Println(InsertSort(num))
}
