package perlin

import (
	"math"
	"math/rand"
)

type vec2 [2]float64

func (v1 vec2) dot(v2 vec2) float64 {
	return v1[0] * v2[0] + v1[1] * v2[1]
}

type Perlin struct {
	rand *rand.Rand
	P []int
}

func (p *Perlin) makePermutation() {
	p.P = make([]int, 512)
	for i := 0; i < 256; i++ {
		p.P[i] = i
	}
	p.rand.Shuffle(256, func(i, j int) {
		p.P[i], p.P[j] = p.P[j], p.P[i]
	})
	copy(p.P[256:], p.P[:256])
}

func fade(x float64) float64 {
	return ((6*x - 15)*x + 10)*x*x*x
}

func lerp(t, a1, a2 float64) float64 {
	return a1 + t*(a2-a1)
}

func NewPerlin(seed int64) *Perlin {
	rand := rand.New(rand.NewSource(seed))
	p := &Perlin{
		rand: rand,
	}
	p.makePermutation()
	return p
}

func getConstantVector(x int) vec2 {
	v := x & 3
	switch v {
	case 0:
		return vec2{1, 1}
	case 1:
		return vec2{-1, 1}
	case 2:
		return vec2{-1, -1}
	default:
		return vec2{1, -1}
	}
}

// Noise2D returns a noise value for x,y. The value returned is between 0;1
func (p *Perlin) Noise2D(x float64, y float64) float64 {
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255

	xf := x-math.Floor(x)
	yf := y-math.Floor(y)

	topRight := vec2{xf-1.0, yf-1.0}
	topLeft := vec2{xf, yf-1.0}
	bottomRight := vec2{xf-1.0, yf}
	bottomLeft := vec2{xf, yf}

	valueTopRight := p.P[p.P[X+1]+Y+1]
	valueTopLeft := p.P[p.P[X]+Y+1]
	valueBottomRight := p.P[p.P[X+1]+Y]
	valueBottomLeft := p.P[p.P[X]+Y]

	dotTopRight := topRight.dot(getConstantVector(valueTopRight));
	dotTopLeft := topLeft.dot(getConstantVector(valueTopLeft));
	dotBottomRight := bottomRight.dot(getConstantVector(valueBottomRight));
	dotBottomLeft := bottomLeft.dot(getConstantVector(valueBottomLeft));

	u := fade(xf)
	v := fade(yf)

	noise := (1 + lerp(u, lerp(v, dotBottomLeft, dotTopLeft), lerp(v, dotBottomRight, dotTopRight))) * 0.5
	if noise >= 1 {
		noise = 0.999
	}
	return noise
}