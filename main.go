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
	mw = 24 // map width
	mh = 24 // map height
)

type Pair[T, U any] struct {
	X T
	Y U
}

type Player struct {
	Index Pair[int, int]
	Dir   v.Vec
}

type game struct {
	screenBuffer *ebiten.Image
	maze         *[mh][mw]int
	p            Player
}

func NewGame() *game {
	maze := [mh][mw]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
	return &game{
		ebiten.NewImage(sw, sh),
		&maze,
		InitPlayer(maze),
	}
}

func InitPlayer(maze [mh][mw]int) Player {
	var x, y int
	for x, y = 0, 0; maze[y][x] != 0; x, y = x+1, y+1 {
	}
	return Player{Pair[int, int]{x, y}, v.Vec{1, 0, 0}}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return sw, sh }
func (g *game) Update() error {
	// Movement:
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.maze[g.p.Index.Y-1][g.p.Index.X] == 0 {
			g.p.Index.Y--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.maze[g.p.Index.Y+1][g.p.Index.X] == 0 {
			g.p.Index.Y++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.maze[g.p.Index.Y][g.p.Index.X+1] == 0 {
			g.p.Index.X++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.maze[g.p.Index.Y][g.p.Index.X-1] == 0 {
			g.p.Index.X--
		}
	}

	// Rotation
	if ebiten.IsKeyPressed(ebiten.KeyQ) {

	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
	}
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	cubeSize := v.Vec{sw / float64(mw), sh / float64(mh), 0}
	for y := range g.maze {
		for x := range g.maze[y] {
			if g.maze[y][x] != 0 {
				ebitenutil.DrawRect(screen, cubeSize.X*float64(x), float64(y)*cubeSize.Y, cubeSize.X, cubeSize.Y, color.RGBA{255, 192, 203, 255})
			}
		}
	}
	playerCenterPos := v.Vec{cubeSize.X*float64(g.p.Index.X) + cubeSize.X/2, cubeSize.Y*float64(g.p.Index.Y) + cubeSize.Y/2, 0}
	ebitenutil.DrawRect(screen, cubeSize.X*float64(g.p.Index.X), cubeSize.Y*float64(g.p.Index.Y), cubeSize.X, cubeSize.Y, color.RGBA{255, 255, 0, 255})
	ebitenutil.DrawLine(screen, playerCenterPos.X, playerCenterPos.Y, playerCenterPos.X+g.p.Dir.X*10, playerCenterPos.Y-g.p.Dir.Y*10, color.RGBA{124, 252, 0, 255})
}

func main() {
	ebiten.SetWindowSize(sw, sh)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
