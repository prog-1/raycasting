package main

import (
	"image/color"
	"log"
	"math"
	"time"

	v "github.com/34thSchool/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sw, sh = 640, 480
	mw, mh = 24, 24                             // map width/height
	cw, ch = float64(sw) / mw, float64(sh) / mh // cell width/height
)

type Pair[T, U any] struct {
	X T
	Y U
}

type Player struct {
	Pos    v.Vec
	Dir    v.Vec
	ms, rs float64 // Movement and Rotation speed
	a, h   int     // base length and height of player's FOV triangle
}

type game struct {
	screenBuffer  *ebiten.Image
	maze          *[mh][mw]int
	p             Player
	prevFrameTime int64
	rc            int //Ray Count
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
		101,
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

	return Player{v.Vec{float64(x)*cw + cw/2, float64(y)*ch + ch/2, 0}, v.Vec{1, 0, 0}, 0.2, 0.01, 500, 250}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return sw, sh }
func (g *game) Update() error {
	dt := float64(time.Now().UnixMilli() - g.prevFrameTime)
	g.prevFrameTime = time.Now().UnixMilli()

	// dp - delta position(new pos - cur pos)
	upp := func(dp v.Vec) { // Update player's position
		np := v.Add(g.p.Pos, dp) // new pos
		in := Pair[int, int]{int(np.X / cw), int(np.Y / ch)}
		if g.maze[in.Y][in.X] == 0 {
			g.p.Pos = np
		}
	}

	// Movement:
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		upp(v.Mul(g.p.Dir, g.p.ms*dt))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		upp(v.Mul(g.p.Dir, -g.p.ms*dt))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		upp(v.Mul(*v.RotateZ(&g.p.Dir, -math.Pi/2), g.p.ms*dt))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		upp(v.Mul(*v.RotateZ(&g.p.Dir, math.Pi/2), g.p.ms*dt))
	}

	// Rotation
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.p.Dir = *v.RotateZ(&g.p.Dir, -g.p.rs*dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.p.Dir = *v.RotateZ(&g.p.Dir, g.p.rs*dt)
	}
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	drawLine := func(a, b v.Vec, clr color.Color) {
		ebitenutil.DrawLine(screen, a.X, a.Y, b.X, b.Y, clr)
	}

	// Maze
	for y := range g.maze {
		for x := range g.maze[y] {
			if g.maze[y][x] != 0 {
				ebitenutil.DrawRect(screen, cw*float64(x), float64(y)*ch, cw, ch, color.RGBA{255, 192, 203, 255})
			}
		}
	}

	// Player
	h := v.Mul(g.p.Dir, float64(g.p.h))

	// Rays
	r := v.Mul(*v.RotateZ(&g.p.Dir, math.Pi/2), float64(g.p.a)/2)
	pa, pb := v.Sub(h, r), v.Add(h, r)
	ab := v.Sub(pb, pa)
	for i := 0; i < g.rc; i++ {
		ak := v.Mul(v.Div(ab, float64(g.rc)), float64(i))
		pk := v.Add(pa, ak)
		k := v.Add(g.p.Pos, pk)
		drawLine(g.p.Pos, k, color.RGBA{255, 255, 0, 255})
	}
}

func main() {
	ebiten.SetWindowSize(sw, sh)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
