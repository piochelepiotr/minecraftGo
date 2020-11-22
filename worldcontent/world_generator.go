package worldcontent

import (
	"encoding/json"
	"github.com/piochelepiotr/minecraftGo/random"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	biomeScale      float64 = 300
	structuresFolder = "./structures/"
)

var maxPerlin2D = math.Sqrt(2)
var maxPerlin3D = math.Sqrt(3)

type biome struct {
	name string
	blockType func(y, elevation int, oreNoise, caveNoise float64) block.Block
	perlin     *perlin.Perlin
	noise *random.Noise
	structures []*savedStructure
	scale int
	maxElevation int
	minElevation int
}

func makeBiome(seed int64, name string, blockType func(y, elevation int, oreNoise, caveNoise float64) block.Block, scale, minElevation, maxElevation int, structures []*structure) *biome {
	b := biome{
		name: name,
		blockType: blockType,
		scale: scale,
		minElevation: minElevation,
		maxElevation: maxElevation,
		perlin:     perlin.NewPerlin(2, 2, 3, seed),
		noise: random.NewNoise(seed),
	}
	var err error
	b.structures, err = loadStructures(name)
	if err != nil {
		log.Fatal("error loading structures", err)
	}
	// for _, s := range structures {
	// 	b.structures = append(b.structures, s.toSavedStructure())
	// }
	// for _, s := range b.structures {
	// 	if err := s.save(s.Name, name); err != nil {
	// 		log.Fatal("error saving struct", err)
	// 	}
	// }
	return &b
}

func (b *biome) getBlock(x, y, z, elevation int) block.Block {
	if y > elevation || y >= WorldHeight {
		return block.Air
	}
	if y == 0 {
		return block.BedRock
	}
	caveNoise := noise3d(b.perlin, x, z, y, 40.0)
	oreNoise := b.noise.Noise3D(x, y, z)
	return b.blockType(y, elevation, oreNoise, caveNoise)
}

type savedStructure struct {
	Name    string `json:"name"`
	Blocks  []byte  `json:"blocks"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Z       int     `json:"z"`
	P       float64 `json:"probability"`
	OriginX int     `json:"origin_x"`
	OriginZ int     `json:"origin_y"`
}

func loadStructure(biome, name string) (*savedStructure, error) {
	path := structuresFolder + biome + "/" + name
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s savedStructure
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return &s, err
}

func loadStructures(biome string) ([]*savedStructure, error) {
	path := structuresFolder + biome
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return nil, err
		}
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	structures := make([]*savedStructure, 0, len(names))
	for _, n := range names {
		s, err := loadStructure(biome, n)
		if err != nil {
			return nil, err
		}
		structures = append(structures, s)
	}
	return structures, nil
}

func (s *savedStructure) save(name string, biome string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := structuresFolder + biome
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(path + "/" + name + ".json", data, 0644); err != nil {
		return err
	}
	return nil
}

func structIndex(sy, sz, x, y, z int) int {
	return x * sy * sz + y * sz + z
}

func (s *savedStructure) toStructure() *structure{
	structure := structure{
		originX: s.OriginX,
		originZ: s.OriginZ,
		p: s.P,
		name: s.Name,
	}
	structure.blocks = make([][][]block.Block, s.X)
	for x := 0; x  < s.X; x++ {
		structure.blocks[x] = make([][]block.Block, s.Y)
		for y := 0; y < s.Y; y++ {
			structure.blocks[x][y] = make([]block.Block, s.Z)
			for z := 0; z < s.Z; z++ {
				structure.blocks[x][y][z] = block.Block(s.Blocks[structIndex(s.Y, s.Z, x, y, z)])
			}
		}
	}
	return &structure
}

func (s *structure) toSavedStructure() (saved *savedStructure) {
	saved = &savedStructure{}
	saved.X = s.x()
	saved.Y = s.y()
	saved.Z = s.z()
	saved.P = s.p
	saved.Name = s.name
	saved.OriginX = s.originX
	saved.OriginZ = s.originZ
	saved.Blocks = make([]byte, s.x() * s.y() * s.z())
	for x := 0; x  < saved.X; x++ {
		for y := 0; y < saved.Y; y++ {
			for z := 0; z < saved.Z; z++ {
				saved.Blocks[structIndex(saved.Y, saved.Z, x, y, z)] = byte(s.blocks[x][y][z])
			}
		}
	}
	return saved
}

type structure struct {
	name string
	blocks  [][][]block.Block
	p       float64
	originX int
	originZ int
}

func (s *structure) x() int {
	return len(s.blocks)
}

func (s *structure) y() int {
	return len(s.blocks[0])
}

func (s *structure) z() int {
	return len(s.blocks[0][0])
}

func makeStructure(x, y, z int) *structure {
	blocks := make([][][]block.Block, x)
	for ix := 0; ix < x; ix++ {
		blocks[ix] = make([][]block.Block, y)
		for iy := 0; iy < y; iy++ {
			blocks[ix][iy] = make([]block.Block, z)
			for iz := 0; iz < z; iz++ {
				blocks[ix][iy][iz] = block.Air
			}
		}
	}
	return &structure{
		blocks: blocks,
	}
}

func noise3d(perlin *perlin.Perlin, x int, y int, z int, scale float64) float64 {
	c := perlin.Noise3D(float64(x)/scale, float64(y)/scale, float64(z)/scale)
	return (c + maxPerlin3D/2) / maxPerlin3D
}

type Generator struct {
	// perlin *perlin.Perlin
	noises *noisesLoader
	biomes []*biome
	maxStructureSizeX int
	maxStructureSizeZ int
}

func makeBiomes(seed int64) []*biome {
	biomes := make([]*biome, 0)
	// biomes = append(biomes, makeBiome(seed*2, "flat", flatBlockType, plainScale, plainMinElevation, plainMinElevation, nil))
	biomes = append(biomes, makeBiome(seed*2, "plain", plainBlockType, plainScale, plainMinElevation, plainMaxElevation, []*structure{makePlainTree(), makePlainTallGrass()}))
	biomes = append(biomes, makeBiome(seed*3, "forest", forestBlockType, forestScale, forestMinElevation, forestMaxElevation, []*structure{makeTree(), makeTallGrass(), makeBirchSampling(), makeRose()}))
	biomes = append(biomes, makeBiome(seed*4, "desert", desertBlockType, desertScale, desertMinElevation, desertMaxElevation, []*structure{makeCactus()}))
	return biomes
}

func newGenerator(worldConfig Config) *Generator {
	biomes:= makeBiomes(worldConfig.Seed)
	g := &Generator{
		biomes: biomes,
		noises: newNoisesLoader(worldConfig.Seed, biomes),
	}
	g.computeMaxStructureSizes()
	return g
}

// blockType returns block, without placing structures
func (g *Generator) blockType(x, y, z int, noises *noisesWithNeighbors) block.Block {
	n := noises.getNoise(x, z)
	biome := g.biomes[n.biome]
	return biome.getBlock(x, y, z, n.elevation)
}

// generateChunkColumn allows you to create a chunk by passing the start point (the second chunk is at position ChunkSize-1)
func (g *Generator) generateChunkColumn(start geometry.Point2D) (chunks []*RawChunk) {
	n := WorldHeight/ChunkSize
	chunks = make([]*RawChunk, 0, n)
	noises := g.noises.getNoisesWithNeighbors(start[0], start[1])
	for y := 0; y < n; y ++ {
		chunk := RawChunk{}
		chunk.start = geometry.Point{start[0], y*ChunkSize, start[1]}
		chunk.blocks = make([]block.Block, ChunkSize*ChunkSize*ChunkSize)
		for x := 0; x < ChunkSize; x++ {
			for z := 0; z < ChunkSize; z++ {
				for y := 0; y < ChunkSize; y++ {
					chunk.blocks[index(x, y, z)] = g.blockType(chunk.start.X+x, chunk.start.Y+y, chunk.start.Z+z, noises)
				}
			}
		}
		chunks = append(chunks, &chunk)
	}
	g.addStructures(chunks, noises)
	return chunks
}

func (g *Generator) computeMaxStructureSizes() {
	var x, z int
	for _, b := range g.biomes {
		for _, s := range b.structures {
			if s.X > x {
				x = s.X
			}
			if s.Z > z {
				z = s.Z
			}
		}
	}
	g.maxStructureSizeX = x
	g.maxStructureSizeZ = z
}

func (g *Generator) addStructures(chunks []*RawChunk, noises *noisesWithNeighbors) {
	startX := chunks[0].start.X
	startZ := chunks[0].start.Z
	setBlock := func(x, y, z int, b block.Block) {
		if b == block.Air {
			return
		}
		if x < startX || x >= startX + ChunkSize || z < startZ || z >= startZ + ChunkSize {
			return
		}
		x = x - startX
		z = z - startZ
		// works only because y > 0
		chunkY := y / ChunkSize
		y = y % ChunkSize
		chunks[chunkY].blocks[index(x, y, z)] = b
	}
	for x := startX - g.maxStructureSizeX; x < startX + ChunkSize + g.maxStructureSizeX; x++ {
		for z := startZ - g.maxStructureSizeZ; z < startZ + ChunkSize + g.maxStructureSizeZ; z++ {
			n := noises.getNoise(x, z)
			if n.structure == -1 {
				continue
			}
			s := g.biomes[n.biome].structures[n.structure]
			xo := x + s.OriginX
			zo := z + s.OriginZ
			yo := noises.getNoise(xo, zo).elevation + 1
			for xs := 0; xs < s.X; xs++ {
				for ys := 0; ys < s.Y; ys++ {
					for zs := 0; zs < s.Z; zs++ {
						setBlock(xs+x, ys+yo, zs+z, block.Block(s.Blocks[structIndex(s.Y, s.Z, xs, ys, zs)]))
					}
				}
			}
		}
	}
}
