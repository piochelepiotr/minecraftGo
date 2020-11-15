package random

import "math/rand"

const sampleSize = 10000

type Noise struct {
	samples []float64
	seed int64
}

func NewNoise(seed int64) *Noise{
	r := rand.New(rand.NewSource(seed))
	noise := &Noise{
		seed: seed,
		samples: make([]float64, sampleSize),
	}
	step := 1.0/float64(sampleSize)
	c := float64(0)
	for i := 0; i < sampleSize; i++ {
		noise.samples[i] = c
		c += step
	}
	r.Shuffle(sampleSize, func(i, j int) {
		noise.samples[i], noise.samples[j] = noise.samples[j], noise.samples[i]
	})
	return noise
}

func (n *Noise) Noise2D(x, y int) float64 {
	i := (x*int(n.seed) + y) % sampleSize
	if i < 0 {
		i = -i
	}
	return n.samples[i]
}

func (n *Noise) Noise3D(x, y, z int) float64 {
	i := (x*int(n.seed) + z*int(n.seed)*int(n.seed) + y) % sampleSize
	if i < 0 {
		i = -i
	}
	return n.samples[i]
}
