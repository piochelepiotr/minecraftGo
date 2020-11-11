package world

import (
	"io/ioutil"
	"log"
)

const (
	savesPath = "saves"
)

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
		w.write(chunk)
	}
	log.Println("closed chunk writer")
	close(w.done)
}

// write chunk to disk
func (w *ChunkWriter) write(c *RawChunk) {
	path := savesPath + "/" + w.worldConfig.Name
	encoded := c.encode()
	key := "chunk_" + c.Start.GetKey()
	if err := ioutil.WriteFile(path + "/" + key, encoded, 0644); err != nil {
		log.Printf("Error when writing chunk to file. %v\n", err)
	}
}

