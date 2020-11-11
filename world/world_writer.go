package world

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

type player struct {
	PosX float32 `json:"pos_x"`
	PosY float32 `json:"pos_y"`
	PosZ float32 `json:"pos_z"`
}

type Config struct {
	Name   string `json:"name"`
	Seed   int64  `json:"seed"`
	Player player `json:"player"`
}

func randomWorld(name string) Config {
	seed := rand.Int63()
	generator := NewGenerator(Config{Seed: seed})
	player := player{PosX: 0, PosY: 0, PosZ: 0}
	player.PosY = float32(generator.GetHeight(0, 0)) + 10
	return Config{
		Seed: seed,
		Name: name,
		Player: player,
	}
}

// GenerateWorld generates a random seed, places the player and saves the world on disk
func GenerateWorld(name string) error {
	config := randomWorld(name)
	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	path := savesPath + "/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(path + "/" + "config.json", data, 0644); err != nil {
		return err
	}
	return nil
}

// LoadWorld loads the world config from local file
func LoadWorld(name string) (config Config, err error) {
	log.Printf("Loading world %s\n", name)
	path := savesPath + "/" + name + "/" + "config.json"
	if _, err := os.Stat(path); err != nil {
		return Config{}, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}

// LoadWorlds loads the list of available worlds
func LoadWorlds() ([]string, error) {
	files, err := ioutil.ReadDir(savesPath)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

