package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 960
	screenHeight = 720
	cellSize     = 8
	rayNum       = 100
)

var (
	worldMap = [25][25]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}
	clr = [5]color.Color{color.RGBA{255, 255, 255, 150}, color.RGBA{255, 0, 0, 150}, color.RGBA{0, 255, 0, 150}, color.RGBA{0, 0, 255, 150}, color.RGBA{0, 255, 255, 150}}
)

type Point struct {
	x, y float64
}

type Game struct {
	width, height   int
	pos, dir, plane *Point
}

func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		pos:    &Point{13, 13},
		dir:    &Point{-1, 0},
		plane:  &Point{0.3, 0.3},
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func rotate(p *Point, angle float64) {
	p.x = p.x*math.Cos(angle) - p.y*math.Sin(angle)
	p.y = p.x*math.Sin(angle) + p.y*math.Cos(angle)
}

func (g *Game) Update() error {
	var mult float64
	if worldMap[int(g.pos.y)][int(g.pos.x)] == 0 {
		mult = 1
	} else {
		mult = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.pos.x += g.dir.x / 10 * mult
		g.pos.y += g.dir.y / 10 * mult
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.pos.x -= g.dir.x / 10 * mult
		g.pos.y -= g.dir.y / 10 * mult
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.pos.x += g.dir.y / 10 * mult
		g.pos.y -= g.dir.x / 10 * mult
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.pos.x -= g.dir.y / 10 * mult
		g.pos.y += g.dir.x / 10 * mult
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		rotate(g.dir, -math.Pi/180)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		rotate(g.dir, math.Pi/180)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range worldMap {
		for j := range worldMap[i] {
			if worldMap[j][i] != 0 {
				vector.DrawFilledRect(screen, float32(i*cellSize), float32(j*cellSize), cellSize, cellSize, clr[worldMap[j][i]-1], false)
			}
		}
	}
	vector.DrawFilledCircle(screen, float32(g.pos.x*cellSize), float32(g.pos.y*cellSize), 3, color.RGBA{255, 255, 0, 150}, false)
	left, right := Point{g.dir.x - g.plane.x, g.dir.y - g.plane.y}, Point{g.dir.x + g.plane.x, g.dir.y + g.plane.y}
	for i := 0.0; i < rayNum; i++ {
		p := Point{g.pos.x, g.pos.y}
		for p.x > 0 && p.x < 24 && p.y > 0 && p.y < 24 {
			p.x, p.y = p.x+left.x, p.y+left.y
		}
		vector.StrokeLine(screen, float32(g.pos.x*cellSize), float32(g.pos.y*cellSize), float32(p.x*cellSize), float32(p.y*cellSize), 3, color.RGBA{255, 255, 0, 50}, false)
		left.x, left.y = left.x+right.x/rayNum, left.y+right.y/rayNum
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
