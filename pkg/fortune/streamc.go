package fortune

import (
	"math/rand"
	"time"
)

// ConcurrentFortuneTeller runs multiple fortune tellers concurrently and performs reservoir sampling on the result.
type ConcurrentFortuneTeller struct {
	tellers []FortuneTeller
	buffer  int
	rand    *rand.Rand
	err     error
}

func NewConcurrentFortuneTeller() *ConcurrentFortuneTeller {
	streams := make([]FortuneTeller, 0)
	for c := range categories() {
		streams = append(streams, NewStreamFortuneTeller().WithCategories(c))
	}
	return &ConcurrentFortuneTeller{
		tellers: streams,
		buffer:  64,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// WithBuffer allows for a max of n fortunes held in memory when processing.
func (f *ConcurrentFortuneTeller) WithBuffer(n int) *ConcurrentFortuneTeller {
	f.buffer = n
	return f
}

func (f *ConcurrentFortuneTeller) Fortune() string {
	c := make(chan string, f.buffer)
	for _, ft := range f.tellers {
		go func(ft FortuneTeller) {
			c <- ft.Fortune()
		}(ft)
	}

	// Reservoir sampling based on candidates we get from c
	reservoir := ""
	for i := 1; i <= len(f.tellers); i++ {
		fortune := <-c
		if err := f.tellers[i-1].Err(); err != nil {
			f.err = err
			continue
		}
		if f.rand.Float32() <= 1/float32(i) {
			reservoir = fortune
		}
	}
	return reservoir
}

func (f *ConcurrentFortuneTeller) Err() error {
	return f.err
}
