package world

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	savesPath = "saves"
)

// write chunk to disk
func (c *RawChunk) write() {
	path := savesPath
	encoded := c.encode()
	key := "chunk_" + c.Start.GetKey()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0755); err != nil {
			log.Printf("error creating directory: %v", err)
		}
	}
	if err := ioutil.WriteFile(path + "/" + key, encoded, 0644); err != nil {
		log.Printf("Error when writing chunk to file. %v\n", err)
	}
}

type ChunkWriter struct {
	in <-chan *RawChunk
	done chan struct{}
}

func NewChunkWriter(in <-chan *RawChunk) (done <-chan struct{}) {
	w := &ChunkWriter{
		in: in,
		done: make(chan struct{}, 0),
	}
	go w.start()
	return w.done
}

func (w *ChunkWriter) start() {
	for chunk := range w.in {
		chunk.write()
	}
	log.Println("closed chunk writer")
	close(w.done)
}