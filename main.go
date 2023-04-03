package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width  = 640.0
	height = 640.0
)

var (
	maze = [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
	oneBlockWidthLength  = width / len(maze[0])
	oneBlockHeightLength = height / len(maze)
)

type Point struct {
	x, y float64
}

type Game struct {
	pos Point
	dir Point
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		a := int(g.pos.x) / int(oneBlockWidthLength)
		b := int(g.pos.y-1) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.y--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		a := int(g.pos.x) / int(oneBlockWidthLength)
		b := int(g.pos.y+1) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.y++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		a := int(g.pos.x-1) / int(oneBlockWidthLength)
		b := int(g.pos.y) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.x--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		a := int(g.pos.x+1) / int(oneBlockWidthLength)
		b := int(g.pos.y) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.x++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dir = Rotate(g.dir, -math.Pi/360)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dir = Rotate(g.dir, math.Pi/360)
	}

	return nil
}
func Rotate(a Point, angle float64) Point {
	a.x, a.y = a.x*math.Cos(angle)-a.y*math.Sin(angle), a.x*math.Sin(angle)+a.y*math.Cos(angle)
	return a
}

func DrawMap(screen *ebiten.Image) {
	for i := range maze {
		for j, WallType := range maze[i] {
			if WallType != 0 {
				ebitenutil.DrawRect(screen, float64(j*oneBlockWidthLength), float64(i*oneBlockHeightLength), float64(oneBlockWidthLength), float64(oneBlockHeightLength), color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	DrawMap(screen)

	ebitenutil.DrawCircle(screen, g.pos.x, g.pos.y, 3, color.RGBA{255, 255, 0, 255})

	for i := -30.0; i < 31; i += 0.5 {
		tmp := Rotate(g.dir, i*math.Pi/180)
		ebitenutil.DrawLine(screen, g.pos.x, g.pos.y, g.pos.x+tmp.x, g.pos.y+tmp.y, color.RGBA{255, 255, 0, 255})
	}

	// ebitenutil.DrawRect(screen, float64(a*oneBlockWidthLength), float64(b*oneBlockHeightLength), float64(oneBlockWidthLength), float64(oneBlockHeightLength), color.RGBA{0, 0, 255, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{Point{320, 350}, Point{0, -1000}}); err != nil {
		log.Fatal(err)
	}
}
