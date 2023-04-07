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
	screenWidth  = 640
	screenHeight = 480
)

type Point struct {
	x, y float64
}
type Game struct {
	width, height       int
	cellSize            float64
	gameMap             [][]int
	playerPos, dir, fov Point
	colors              []color.Color
}

func rotate(a Point, ang float64) Point {
	return Point{a.x*math.Cos(ang) - a.y*math.Sin(ang), a.x*math.Sin(ang) + a.y*math.Cos(ang)}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.dir = rotate(g.dir, math.Pi/360)
		g.fov = rotate(g.fov, math.Pi/360)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.dir = rotate(g.dir, math.Pi/-360)
		g.fov = rotate(g.fov, math.Pi/-360)
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) && g.gameMap[int(g.playerPos.y+g.dir.y/10)][int(g.playerPos.x+g.dir.x/10)] == 0 {
		g.playerPos.x += g.dir.x / 10
		g.playerPos.y += g.dir.y / 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.gameMap[int(g.playerPos.y-g.dir.x/10)][int(g.playerPos.x+g.dir.y/10)] == 0 {
		g.playerPos.x += g.dir.y / 10
		g.playerPos.y -= g.dir.x / 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.gameMap[int(g.playerPos.y-g.dir.y/10)][int(g.playerPos.x-g.dir.x/10)] == 0 {
		g.playerPos.x -= g.dir.x / 10
		g.playerPos.y -= g.dir.y / 10
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.gameMap[int(g.playerPos.y+g.dir.x/10)][int(g.playerPos.x-g.dir.y/10)] == 0 {
		g.playerPos.x -= g.dir.y / 10
		g.playerPos.y += g.dir.x / 10
	}
	return nil
}

func (g *Game) DrawFieldOfView(screen *ebiten.Image) {
	for i := 0.0; i < 100; i++ {
		var a, step Point
		d := Point{g.dir.x + g.fov.x*i/50, g.dir.y + g.fov.y*i/50}
		p := Point{g.playerPos.x * g.cellSize, g.playerPos.y * g.cellSize}
		deltad := Point{math.Abs(1 / d.x), math.Abs(1 / d.y)}
		if d.x < 0 {
			step.x = -1
		} else {
			step.x = 1
			a.x = deltad.x
		}
		if d.y < 0 {
			step.y = -1
		} else {
			step.y = 1
			a.y = deltad.y
		}
		for g.gameMap[int(p.y/g.cellSize)][int(p.x/g.cellSize)] == 0 {
			if a.x < a.y {
				a.x += deltad.x
				p.x += step.x
			} else {
				a.y += deltad.y
				p.y += step.y
			}
		}
		vector.StrokeLine(screen, float32(g.playerPos.x*g.cellSize), float32(g.playerPos.y*g.cellSize), float32(p.x), float32(p.y), 2, color.RGBA{255, 255, 133, 255}, false)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.gameMap {
		for j := range g.gameMap[i] {
			if g.gameMap[j][i] != 0 {
				vector.DrawFilledRect(screen, float32(i*int(g.cellSize)), float32(j*int(g.cellSize)), float32(g.cellSize), float32(g.cellSize), g.colors[g.gameMap[j][i]-1], false)
			}
		}
	}
	vector.DrawFilledCircle(screen, float32(g.playerPos.x*g.cellSize), float32(g.playerPos.y*g.cellSize), 2, color.RGBA{255, 255, 133, 255}, false)
	g.DrawFieldOfView(screen)
}

func NewGame(width, height int) *Game {
	return &Game{
		width:     width,
		height:    height,
		cellSize:  15,
		playerPos: Point{2, 2},
		dir:       Point{1, 0},
		fov:       Point{0.3, 0.3},
		gameMap: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 3, 3, 0, 2, 0, 0, 0, 2, 0, 3, 3, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 2, 2, 0, 2, 2, 2, 2, 2, 2, 2, 0, 2, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 2, 0, 0, 2, 0, 2, 0, 2, 0, 2, 0, 0, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 2, 0, 2, 2, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 2, 2, 2, 0, 0, 2, 0, 2, 0, 0, 2, 2, 2, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 2, 2, 2, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		colors: []color.Color{color.RGBA{255, 255, 255, 150},
			color.RGBA{0, 255, 0, 150},
			color.RGBA{255, 0, 0, 150}},
	}
}
func (g *Game) Layout(outWidth, outHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
