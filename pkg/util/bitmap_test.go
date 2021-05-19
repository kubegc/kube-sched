package util

import (
	"fmt"
	"testing"
)

func TestBitmap(t *testing.T) {
	bm := NewBitMap64(1000)
	fmt.Println(bm.Acquire())
	fmt.Println(bm.Acquire())
	bm.Release(0)
	fmt.Println(bm.Acquire())
}