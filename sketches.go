package ada

import (
	"fmt"
	"math"
)

// Sketches ...
type Sketches struct {
	sketches    []*Sketch
	maxDuration uint64
}

// NewSketches ...
func NewSketches(maxDuration, w, d uint64, alpha float64) *Sketches {
	num := uint64(math.Log2(float64(maxDuration)))
	sketches := make([]*Sketch, num)
	for i := range sketches {
		sketches[i] = NewSketch(w, d, alpha)
	}
	return &Sketches{
		sketches:    sketches,
		maxDuration: maxDuration,
	}
}

// Update ...
func (sks *Sketches) Update(item []byte, timestamp, count uint64) error {
	for i := range sks.sketches {
		t := 1 + timestamp/uint64(math.Pow(2, float64(i)))
		sks.sketches[i].Update(item, t, count)
	}
	return nil
}

// Estimate ...
func (sks *Sketches) Estimate(item []byte, start, end uint64) (uint64, error) {
	if end-start > sks.maxDuration {
		return 0, fmt.Errorf("window to big [start, end] %d > %d", end-start, sks.maxDuration)
	}
	estimate := uint64(0)
	for start <= end {
		pow2 := float64(start & (^start + 1))
		logpow2 := math.Log2(pow2)
		if logpow2 < 0 {
			logpow2 = 0
		}
		for i := logpow2; i >= 0; i-- {
			if start+uint64(math.Pow(2, i))-1 <= end {
				t := 1 + start/uint64(math.Pow(2, i))
				estimate += sks.sketches[uint64(i)].Count(item, t)
				start += uint64(math.Pow(2, i))
			}
		}
	}
	return estimate, nil
}

// EstimateOverMaxDuration ...
func (sks *Sketches) EstimateOverMaxDuration(item []byte, start, end uint64) uint64 {
	estimate := uint64(0)
	for i := start; i <= end; i += sks.maxDuration {
		tmpEnd := i + sks.maxDuration - 1
		if tmpEnd > end {
			tmpEnd = end
		}
		res, _ := sks.Estimate(item, i, tmpEnd)
		estimate += res
	}
	return estimate
}
