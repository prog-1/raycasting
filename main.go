package main

import (
	"image/color"
	"log"

	v "github.com/34thSchool/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sw = 640
	sh = 480
	mw = 10 // map width
	mh = 10 // map height
)

type game struct {
	screenBuffer *ebiten.Image
	maze         [mh][mw]int
}

func NewGame() *game {
	return &game{
		ebiten.NewImage(sw, sh),
		[mh][mw]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 1, 0, 0, 1},
			{1, 0, 1, 1, 0, 0, 1, 1, 1, 1},
			{1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
			{1, 1, 1, 1, 0, 0, 1, 0, 1, 1},
			{1, 0, 0, 0, 0, 0, 1, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return sw, sh }
func (g *game) Update() error {
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	cubeSize := v.Vec{sw / mw, sh / mh, 0}
	for y := range g.maze {
		for x := range g.maze[y] {
			if g.maze[y][x] != 0 {
				ebitenutil.DrawRect(screen, cubeSize.X*float64(x), float64(y)*cubeSize.Y, cubeSize.X, cubeSize.Y, color.White)
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(sw, sh)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
