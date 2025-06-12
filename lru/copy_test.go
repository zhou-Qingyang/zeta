package main

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{8, 2, 1, 3}
	copy(slice2, slice1)
	fmt.Println(slice2) // 1 2 3 ,3
}
