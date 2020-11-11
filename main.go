package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/piochelepiotr/minecraftGo/game"
	"github.com/piochelepiotr/minecraftGo/render"
)

const windowWidth = 800
const windowHeight = 600
const aspectRatio = windowWidth / windowHeight

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	d := render.NewDisplay(windowWidth, windowHeight)
	defer d.CloseDisplay()
	game.Start(aspectRatio, d)
}
