package util

import (
	"fmt"
	"testing"
)

func TestBitmap(t *testing.T) {
	bm := Bitmap64{
		base: 50051,
		bits: nil,
	}


	bm.Mask(50051)
	fmt.Println(bm.Test(50051))
}