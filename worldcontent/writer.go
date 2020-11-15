package worldcontent

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	savesPath = "saves"
)

func writeChunk(config Config, c *RawChunk) {
	path := savesPath + "/" + config.Name
	encoded := c.encode()
	key := "chunk_" + c.start.GetKey()
	if err := ioutil.WriteFile(path + "/" + key, encoded, 0644); err != nil {
		log.Printf("Error when writing chunk to file. %v\n", err)
	}
}

type ChunkWriter struct {
	in <-chan *RawChunk
	done chan struct{}
	worldConfig Config
}

func NewChunkWriter(worldConfig Config, in <-chan *RawChunk) (done <-chan struct{}) {
	w := &ChunkWriter{
		in: in,
		done: make(chan struct{}, 0),
		worldConfig: worldConfig,
	}
	go w.start()
	return w.done
}

func (w *ChunkWriter) start() {
	for chunk := range w.in {
		writeChunk(w.worldConfig, chunk)
	}
	log.Println("closed chunk writer")
	close(w.done)
}

func WriteWorld(config Config) error {
	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	path := savesPath + "/" + config.Name
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

