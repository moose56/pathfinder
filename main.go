package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	grid := NewGrid(14, 14)
	game := NewGame(grid, 50, 20)

	ebiten.SetWindowSize(game.Width(), game.Height())
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
