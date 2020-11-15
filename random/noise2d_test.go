package random

import (
	"fmt"
	"testing"
)

func TestNewNoise(t *testing.T) {
	total := 100
	n := 0
	p := 0.1
	r := NewNoise(234234)
	for x := 0; x < total; x++ {
		for y := 0; y < total; y++ {
			if r.Noise2D(x, y) <= p {
				n++
			}
		}
	}
	fmt.Println(float64(n)/float64(total*total))
}
