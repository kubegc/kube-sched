package util

import "fmt"

type Bitmap interface {
	Clear()
	NextAndMask() int
	Mask(x int)
	Unmask(x int)
	Test(x int) bool
}


type Bitmap64 struct {
	base int
	bits []uint64
}

func NewBitMap64(base int) *Bitmap64 {
	return &Bitmap64{
		base: base,
	}
}

func (bm *Bitmap64) Clear() {
	bm.bits = nil
}

func (bm *Bitmap64) Test(x int) bool {
	diff := x - bm.base
	idx, offset := diff / 64, diff % 64
	return (bm.bits[idx] >> offset) & 1 > 0
}
func (bm *Bitmap64) NextAndMask() int {
	for i := 0; i < len(bm.bits); i++ {
		fmt.Println(bm.bits[i])
		if bm.bits[i] != 0xffffffffffffffff {
			for j := 0; j < 64; j++ {
				if (bm.bits[i] >> j) & 1 == 0 {
					pos := 64 * i + j + bm.base
					bm.Mask(pos)
					return pos
				}
			}
		}
	}
	bm.Mask(len(bm.bits) * 64 + bm.base)
	return len(bm.bits) * 64 + bm.base
}

func (bm *Bitmap64) Mask(x int) {
	diff := x - bm.base
	idx, offset := diff / 64, diff % 64
	for i := len(bm.bits); i <= idx; i++ {
		bm.bits = append(bm.bits, 0)
	}
	bm.bits[idx] = bm.bits[idx] | (1 << offset)
}

func (bm *Bitmap64) Unmask(x int) {
	diff := x - bm.base
	idx, offset := diff / 64, diff % 64
	bm.bits[idx] = bm.bits[idx] & (^(1 << offset))
}


