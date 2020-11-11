package world

type player struct {
	posX float32 `json:"pos_x"`
	posY float32 `json:"pos_y"`
	posZ float32 `json:"pos_z"`
}

type Config struct {
	name string `json:"name"`
	player player `json:"player"`
}

// GenerateWorld generates a random seed, places the player and saves the world on disk
func GenerateWorld(name string) {

}

// LoadWorld loads the world config from local file
func LoadWorld(name string) Config {
	return Config{}
}