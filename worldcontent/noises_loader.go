package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/random"
	"math"
	"sort"
)

type noises2D struct {
	biome int
	// for the correct biome
	elevation int
	structure int
}

func elevation(noise float64, min int, max int) int {
	return min + int(float64(max-min)*noise)
}

// if distance != 1, we are at the border
const borderSize = 0.4
func distanceFromBiomeBorder(p float64) float64 {
	lower := float64(int(p))
	upper := lower + 1
	toLower := (p - lower)*(1/borderSize)
	toUpper := (upper - p)*(1/borderSize)
	minDistance := math.Min(toLower, toUpper)
	return math.Min(minDistance, 1)
}


type chunkNoises2D []noises2D

func (n chunkNoises2D) get(x, z int) noises2D {
	return n[z*ChunkSize + x]
}

func (n chunkNoises2D)set(x, z int, noise noises2D) {
	n[z*ChunkSize + x] = noise
}

type noisesWithNeighbors struct {
	chunks []chunkNoises2D
	originX int
	originZ int
}

type point2d [2]int

// noisesLoader computes & stores noises for chunks, to avoid computing them multiple times
type noisesLoader struct {
	noises           map[point2d]chunkNoises2D
	biomePerlin      *perlin.Perlin
	elevationPerlins []*perlin.Perlin
	biomes           []*biome
	structureNoise   *random.Noise
	structuresCumProbaPerBiome [][]float64
}

func (n *noisesLoader) computeStructuresCumProbaPerBiome() {
	for _, b := range n.biomes {
		probas := make([]float64, len(b.structures))
		sum := float64(0)
		for i, s := range b.structures {
			sum += s.P
			probas[i] = sum
		}
		n.structuresCumProbaPerBiome = append(n.structuresCumProbaPerBiome, probas)
	}
}

func (n *noisesLoader) getNoisesWithNeighbors(chunkX, chunkZ int) *noisesWithNeighbors {
	noises := noisesWithNeighbors{
		originX: chunkX,
		originZ: chunkZ,
		chunks: make([]chunkNoises2D, 9),
	}
	i := 0
	for x := -1; x <= 1; x++ {
		for z := -1; z <= 1; z++ {
			noises.chunks[i] = n.get(chunkX + x*ChunkSize, chunkZ + z*ChunkSize)
			i++
		}
	}
	return &noises
}

func (n *noisesWithNeighbors) getNoise(x, z int) noises2D {
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
	return n.chunks[i].get(x - chunkX, z - chunkZ)
}

func newNoisesLoader(seed int64, biomes []*biome) *noisesLoader {
	n := noisesLoader{
		noises:         make(map[point2d]chunkNoises2D),
		biomePerlin:    perlin.NewPerlin(1.1, 2, 3, seed),
		biomes:         biomes,
		structureNoise: random.NewNoise(seed),
	}
	n.elevationPerlins = make([]*perlin.Perlin, 0, len(biomes))
	for i := 0; i < len(biomes); i++ {
		n.elevationPerlins = append(n.elevationPerlins, perlin.NewPerlin(2, 2, 3, seed*int64(i+2)))
	}
	n.computeStructuresCumProbaPerBiome()
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
	noises := make(chunkNoises2D, ChunkSize*ChunkSize)
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			noises.set(x, z, n.noiseAt(x + chunkX, z + chunkZ))
		}
	}
	return noises
}

func (n *noisesLoader) noiseAt(x, z int) noises2D {
	noises := noises2D{}
	biomeNoise := noise2d(n.biomePerlin, x, z, biomeScale)*float64(len(n.biomes))
	noises.biome = int(biomeNoise)
	distanceToNextBiome := distanceFromBiomeBorder(biomeNoise)
	noises.elevation = elevation(noise2d(n.elevationPerlins[noises.biome], x, z, float64(n.biomes[noises.biome].scale))*distanceToNextBiome, n.biomes[noises.biome].minElevation, n.biomes[noises.biome].maxElevation)
	structureNoise := n.structureNoise.Noise2D(x, z)
	noises.structure = sort.SearchFloat64s(n.structuresCumProbaPerBiome[noises.biome], structureNoise)
	if noises.structure == len(n.biomes[noises.biome].structures) {
		// nothing
		noises.structure = -1
	}
	return noises
}

func noise2d(p *perlin.Perlin, x int, y int, scale float64) float64 {
	c := p.Noise2D(float64(x)/scale, float64(y)/scale)
	c = (c + maxPerlin2D/2) / maxPerlin2D
	if c < 0 {
		c = 0
	} else if c >= 1 {
		c = 0.9999999
	}
	return c
}