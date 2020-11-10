package world

import (
	"io/ioutil"
	"log"
	"os"
)

// write chunk to disk
func (c *RawChunk) write() {
	path := "~/Documents/perso/minecraftGoSaves"
	encoded := c.encode()
	key := c.getKey()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
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