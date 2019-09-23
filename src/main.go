package main

import (
	"log"
	"v2sdl/display"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.Main(func() {})

	d, err := display.NewDisplay("Hello")
	if err != nil {
		log.Panic(err)
	}
	defer d.Close()
	d.EventLoop()
}
