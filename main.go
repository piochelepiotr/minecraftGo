package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/piochelepiotr/minecraftGo/game"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

const windowWidth = 800
const windowHeight = 600

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	err := profiler.Start(
		profiler.WithService("minecraft-go"),
		profiler.WithEnv("dev"),
		profiler.WithVersion("v0.5"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()
	rand.Seed(time.Now().UnixNano())
	d := render.NewDisplay(windowWidth, windowHeight)
	defer d.CloseDisplay()
	game.Start(d)
}
