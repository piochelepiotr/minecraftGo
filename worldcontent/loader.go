package worldcontent

import (
	"encoding/json"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"io/ioutil"
	"os"
)

func loadChunkFromSaves(config Config, start geometry.Point) (chunk RawChunk, err error) {
	path := savesPath + "/" + config.Name + "/" + "chunk_" + start.GetKey()
	if _, err := os.Stat(path); err != nil {
		return chunk, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return chunk, err
	}
	chunk = decode(config, data, start)
	return chunk, nil
}

// getChunk tries to load the chunk from a saved file. If there is nothing, generates one using the generator
func getChunk(worldConfig Config, start geometry.Point, generator *Generator) RawChunk {
	if chunk, err := loadChunkFromSaves(worldConfig, start); err == nil {
		return chunk
	}
	return generator.generateChunk(start, worldConfig.WorldHeight)
}

// LoadWorld loads the world config from local file
func LoadWorld(name string) (config Config, err error) {
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
	path := savesPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			return nil, err
		}
	}
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