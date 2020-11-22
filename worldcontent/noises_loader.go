package worldcontent

import (
	"github.com/aquilax/go-perlin"
)

type noises2D struct {
	biomeNoise float64
	// 1 different elevation for each biome
	elevationNoises [nBiomes]float64
}

type chunkNoises2D []noises2D

func (n chunkNoises2D) get(x, z int) noises2D {
	return n[z*ChunkSize + x]
}

func (n chunkNoises2D)set(x, z int, noise noises2D) {
	n[z*ChunkSize + x] = noise
}

func newChunkNoises2D() chunkNoises2D {
	return make(chunkNoises2D, ChunkSize * ChunkSize)
}

type noisesWithNeighbors struct {
	chunks [9]chunkNoises2D
	originX int
	originZ int
}

type point2d [2]int

// noisesLoader computes & stores noises for chunks, to avoid computing them multiple times
type noisesLoader struct {
	noises map[point2d]chunkNoises2D
	biomePerlin *perlin.Perlin
	elevationPerlins [nBiomes]*perlin.Perlin
	scales [nBiomes]int
}

func (n *noisesLoader) getNoisesWithNeighbors(chunkX, chunkZ int) noisesWithNeighbors {
	noises := noisesWithNeighbors{
		originX: chunkX,
		originZ: chunkZ,
	}
	i := 0
	for x := -1; x <= 1; x++ {
		for z := -1; z <= 1; z++ {
			noises.chunks[i] = n.get(chunkX + x*ChunkSize, chunkZ + z*ChunkSize)
			i++
		}
	}
	return noises
}

func (n noisesWithNeighbors) getNoise(x, z int) noises2D {
	chunkX := n.originX
	chunkZ := n.originZ
	// the middle is the origin chunk
	i := 4
	if x < n.originX {
		chunkX -= ChunkSize
		i -= 3
	} else if x >= n.originX + ChunkSize {
		chunkX += ChunkSize
		i += 3
	}
	if z < n.originZ {
		chunkZ -= ChunkSize
		i -= 1
	} else if z >= n.originZ + ChunkSize {
		chunkZ += ChunkSize
		i += 1
	}
	c := n.chunks[i]
	return c.get(x - chunkX, z - chunkZ)
}

func newNoisesLoader(seed int64) *noisesLoader {
	n := noisesLoader{
		noises: make(map[point2d]chunkNoises2D),
		biomePerlin: perlin.NewPerlin(2, 2, 3, seed),
	}
	for i := 0; i < nBiomes; i++ {
		n.elevationPerlins[i] = perlin.NewPerlin(2, 2, 3, seed*int64(i+2))
		n.scales[i] = 100
	}
	return &n
}

func (n *noisesLoader) get(chunkX, chunkZ int) chunkNoises2D {
	p := point2d{chunkX, chunkZ}
	if noise, ok := n.noises[p]; ok {
		return noise
	}
	noise := n.generateNoises(chunkX, chunkZ)
	n.noises[p] = noise
	return noise
}

func (n *noisesLoader) generateNoises(chunkX, chunkZ int) chunkNoises2D {
	noises := newChunkNoises2D()
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			noises.set(x, z, n.noiseAt(x + chunkX, z + chunkZ))
		}
	}
	return noises
}

func (n *noisesLoader) noiseAt(x, z int) noises2D {
	noises := noises2D{}
	noises.biomeNoise = noise2d(n.biomePerlin, x, z, biomeScale)
	for i := 0; i < nBiomes; i++ {
		noises.elevationNoises[i] = noise2d(n.elevationPerlins[i], x, z, float64(n.scales[i]))
	}
	return noises
}

func noise2d(p *perlin.Perlin, x int, y int, scale float64) float64 {
	c := p.Noise2D(float64(x)/scale, float64(y)/scale)
	c = (c + maxPerlin2D/2) / maxPerlin2D
	return c
}

