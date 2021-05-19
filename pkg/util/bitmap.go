package util

type Bitmap interface {
	clear()
	nextAndMask() int
	mask(x int)
	unmask(x int)
	test(x int) bool

	Acquire() int
	Release(x int)
}


type Bitmap64 struct {
	capacity int
	current int
	bits []uint64
}

func NewBitMap64(capacity int) *Bitmap64 {
	return &Bitmap64{
		capacity: capacity,
		current: 0,
		bits: make([]uint64, (capacity + 63) / 64),
	}
}

func (bm *Bitmap64) clear() {
	bm.capacity = 0
	bm.bits = nil
	bm.current = 0
}

func (bm *Bitmap64) test(x int) bool {
	idx, offset := x / 64, x % 64
	return (bm.bits[idx] >> offset) & 1 > 0
}
func (bm *Bitmap64) nextAndMask() int {
	for i := bm.current; i < bm.current + bm.capacity; i ++ {
		i %= bm.capacity
		idx, offset := i / 64, i % 64

		if bm.bits[idx] >> offset & 1 == 0 {
			bm.mask(i)
			bm.current = i + 1
			return i
		}
	}
	return -1
}

func (bm *Bitmap64) mask(x int) {
	idx, offset := x / 64, x % 64
	bm.bits[idx] = bm.bits[idx] | (1 << offset)
}

func (bm *Bitmap64) unmask(x int) {
	idx, offset := x / 64, x % 64
	bm.bits[idx] = bm.bits[idx] & (^(1 << offset))
}

func(bm *Bitmap64) Acquire() int {
	return bm.nextAndMask()
}

func(bm *Bitmap64) Release(x int) {
	bm.unmask(x)
}


