package main

import (
	"image/color"
	"log"
	"time"

	v "github.com/34thSchool/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sw, sh   = 640, 480
	mw, mh   = 24, 24                             // map width/height
	sqw, sqh = float64(sw) / mw, float64(sh) / mh // square width/height
)

type Pair[T, U any] struct {
	X T
	Y U
}

type Player struct {
	Pos   v.Vec
	Dir   v.Vec
	speed float64
}

type game struct {
	screenBuffer  *ebiten.Image
	maze          *[mh][mw]int
	p             Player
	prevFrameTime int64
}

func NewGame() *game {
	// Values on the map:
	// 0 - empty space
	// 1 - player(yellow)
	// 2 - white wall
	// 3 - pink wall
	// 4 - green wall
	maze := [mh][mw]int{
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		{2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
	}
	return &game{
		ebiten.NewImage(sw, sh),
		&maze,
		InitPlayer(maze),
		0,
	}
}

func InitPlayer(maze [mh][mw]int) Player {
	var x, y int
	for y = range maze {
		for x = range maze[y] {
			if maze[y][x] == 0 {
				break
			}
		}
		if maze[y][x] == 0 {
			break
		}
	}

	return Player{v.Vec{float64(x)*sqw + sqw/2, float64(y)*sqh + sqh/2, 0}, v.Vec{1, 0, 0}, 0.2}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return sw, sh }
func (g *game) Update() error {
	dt := float64(time.Now().UnixMilli() - g.prevFrameTime)
	g.prevFrameTime = time.Now().UnixMilli()

	// Movement:
	in := Pair[int, int]{int(g.p.Pos.X / sqw), int(g.p.Pos.Y / sqh)} // pos = in * sqw => in = pos/sqw
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.maze[in.Y-1][in.X] == 0 {
			g.p.Pos.Y -= g.p.speed * dt
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.maze[in.Y+1][in.X] == 0 {
			g.p.Pos.Y += g.p.speed * dt
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.maze[in.Y][in.X-1] == 0 {
			g.p.Pos.X -= g.p.speed * dt
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.maze[in.Y][in.X+1] == 0 {
			g.p.Pos.X += g.p.speed * dt
		}
	}

	// Rotation
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.p.Dir.RotateZ(0.01 * dt)
	}
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	for y := range g.maze {
		for x := range g.maze[y] {
			if g.maze[y][x] != 0 {
				ebitenutil.DrawRect(screen, sqw*float64(x), float64(y)*sqh, sqw, sqh, color.RGBA{255, 192, 203, 255})
			}
		}
	}
	in := Pair[int, int]{int(g.p.Pos.X / sqw), int(g.p.Pos.Y / sqh)} // pos = in * sqw => in = pos/sqw
	ebitenutil.DrawRect(screen, sqw*float64(in.X), sqh*float64(in.Y), sqw, sqh, color.RGBA{255, 255, 0, 255} /*yellow*/)
	ebitenutil.DrawLine(screen, g.p.Pos.X, g.p.Pos.Y, g.p.Pos.X+g.p.Dir.X*10, g.p.Pos.Y-g.p.Dir.Y*10, color.RGBA{124, 252, 0, 255})
}

func main() {
	ebiten.SetWindowSize(sw, sh)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
