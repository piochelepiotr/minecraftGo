package worldcontent

import "math/rand"

const (
	// ChunkSize is the size of a chunk in blocks
	ChunkSize = 16
	WorldHeight = ChunkSize*3
)

type Player struct {
	PosX float32 `json:"pos_x"`
	PosY float32 `json:"pos_y"`
	PosZ float32 `json:"pos_z"`
}

type Config struct {
	Name   string `json:"name"`
	Seed   int64  `json:"seed"`
	Player Player `json:"player"`
	Creative bool `json:"creative"`
	ChunkSize int `json:"chunk_size"`
	WorldHeight int `json:"world_height"`
}

func GetRandomWorld(name string) Config {
	seed := rand.Int63()
	// generator := NewGenerator(Config{Seed: seed})
	player := Player{PosX: 0, PosY: -1, PosZ: 0}
	// player.PosY = float32(generator.GetHeight(0, 0)) + 10
	return Config{
		Creative: true,
		Seed: seed,
		Name: name,
		Player: player,
	}
}
