package worldcontent

import (
	"fmt"
	"github.com/aquilax/go-perlin"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	max := float64(0)
	min := float64(0)
	n := 4
	// we verified that value for n = 1 is between -sqrt(2)/2 and sqrt(2)/2
	for i := 0; i < 1000; i++ {
		p := perlin.NewPerlin(2, 2, 1, rand.Int63())
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				x := p.Noise2D(float64(x)+0.5, float64(y)+0.5)
				if x > max {
					max = x
				}
				if x < min {
					min = x
				}
				// fmt.Println(x)
			}
		}
	}
	fmt.Println("max", max, "min", min)
}
