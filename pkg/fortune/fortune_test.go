package fortune_test

import (
	"fmt"
	"testing"

	"github.com/kdevo/gofort/pkg/fortune"
)

func TestFortune(t *testing.T) {
	testCases := []struct {
		name   string
		teller fortune.FortuneTeller
	}{
		{"stream", fortune.NewStreamFortuneTeller()},
		{"simple", fortune.NewStreamFortuneTeller().WithBufferSize(64 * 1024)},
		{"concurrent", fortune.NewConcurrentFortuneTeller()},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fortune := tc.teller.Fortune()
			fmt.Println(fortune)
			if err := tc.teller.Err(); err != nil {
				t.Errorf("Got error: %s", err)
			}
			if fortune == "" {
				t.Fatal("Fortune() must always return a text!")
			}
		})
	}
}

// BenchmarkFortune sets the different FortuneTeller implementations into comparison.
// They are ordered from usually fastest to slowest.
//
// Interpretation:
// Even though your mileage could theoretically vary due to system-specific differences,
// the benchmark usually indicates that the reservoir sampling algorithm in NewStreamFortuneTeller is both,
// fastest and least memory consuming. The concurrent implementation NewConcurrentFortuneTeller is not faster,
// which is due to the common Go idiom "Concurrency is not parallelism". The files might be read concurrently,
// but due to already reaching the maximum I/O throughput by simple sequential scanning, performance does not improve.
// To the contrary, if we would do heavy calculations on the read fortunes for some reason,
// machines with multiple cores could even profit from parallelism and therefore more speed.
func BenchmarkFortune(b *testing.B) {
	benchmarks := []struct {
		name   string
		teller fortune.FortuneTeller
	}{
		{"++    reservoir sampling", fortune.NewStreamFortuneTeller()},
		{"+     reservoir sampling", fortune.NewStreamFortuneTeller().WithBufferSize(64 * 1024)},
		{"-             concurrent", fortune.NewConcurrentFortuneTeller()},
		{"--            concurrent", fortune.NewConcurrentFortuneTeller().WithBuffer(0)},
		{"---              readall", fortune.NewSimpleFortuneTeller()},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if bm.teller.Fortune() == "" {
					if err := bm.teller.Err(); err != nil {
						b.Errorf("Got error: %s", err)
					}
					b.Fatal("Fortune() must always return a text!")
				}
			}
		})
	}
}
