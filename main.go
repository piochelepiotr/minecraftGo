package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/piochelepiotr/minecraftGo/game"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
)

const windowWidth = 800
const windowHeight = 600

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	rand.Seed(time.Now().UnixNano())
	d := render.NewDisplay(windowWidth, windowHeight)
	defer d.CloseDisplay()
	game.Start(d)
}
