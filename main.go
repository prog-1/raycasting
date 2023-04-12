package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	winTitle     = "Cube"
	screenWidth  = 1000
	screenHeight = 1000
	dpi          = 100
)

type (
	point struct {
		x, y float64
	}
	game struct {
		m             *ebiten.Image
		p, dir, plane point
		pg            [][]int
		// lp            point
	}
)

func rotate(p *point, angle float64) {
	x, y := p.x, p.y

	p.x = (x*math.Cos(angle) - y*math.Sin(angle))
	p.y = (x*math.Sin(angle) + y*math.Cos(angle))

}

// func rotate2(p *point, angle float64) {
// 	x, y := p.x, p.y

// 	p.x = (x*math.Cos(angle) - y*math.Sin(angle))
// 	p.y = (x*math.Sin(angle) + y*math.Cos(angle))

// }

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		rotate(&g.dir, math.Pi/100)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		rotate(&g.dir, -math.Pi/100)

	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if g.pg[((int(g.p.y)-38)/50)+1][int(g.p.x)/50] == 0 {
			g.p.y += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if g.pg[(int(g.p.y))/50][int(g.p.x+38)/50-1] == 0 {
			g.p.x -= 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if g.pg[(int(g.p.y))/50][int(g.p.x-38)/50+1] == 0 {
			g.p.x += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if g.pg[(int(g.p.y+38))/50-1][int(g.p.x)/50] == 0 {
			g.p.y -= 5
		}
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.m, nil)
	ebitenutil.DrawCircle(screen, g.p.x, g.p.y, 10, color.RGBA{0xff, 0xff, 0x00, 0xff})
	for i := 0.0; i < 400; i++ {
		var side, step point
		cameraX := i/200.0 - 1
		ray := point{g.dir.x + g.plane.x*cameraX, g.dir.y + g.plane.y*cameraX}
		p := point{g.p.x, g.p.y}
		delta := point{math.Abs(1 / ray.x), math.Abs(1 / ray.y)}
		if ray.x < 0 {
			step.x = -1
			side.x = (g.p.x - p.x) * delta.x
		} else {
			step.x = 1
			side.x = (p.x + 1.0 - g.p.x) * delta.x
		}
		if ray.y < 0 {
			step.y = -1
			side.y = (g.p.y - p.y) * delta.y
		} else {
			step.y = 1
			side.y = (p.y + 1.0 - g.p.y) * delta.y
		}
		for g.pg[int(p.y)/50][int(p.x)/50] == 0 {
			if side.x < side.y {
				side.x += delta.x
				p.x += step.x
			} else {
				side.y += delta.y
				p.y += step.y
			}
		}
		ebitenutil.DrawLine(screen, g.p.x, g.p.y, p.x, p.y, color.RGBA{0xff, 0xff, 0x00, 0xff})
	}

}

func main() {
	ebiten.SetWindowTitle(winTitle)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
func DrawBackground(m *ebiten.Image, pg [][]int) *ebiten.Image {
	var cntW, cntH float64
	var prevI int
	for i := range pg {
		for j := range pg[i] {
			if i != prevI {
				cntW = 0
				cntH++
				prevI = i
			}
			if pg[i][j] > 0 {
				ebitenutil.DrawRect(m, 50*cntW, 50*cntH, 50, 50, color.RGBA{0xff, 0xff, 0xff, 0xff})
			}
			cntW++

		}
	}
	return m
}
func NewGame() *game {
	pg := [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	return &game{
		m:  DrawBackground(ebiten.NewImage(screenWidth, screenHeight), pg),
		p:  point{100, 100},
		pg: pg,
		// lp:    point{300, 500},
		dir:   point{1, 0},
		plane: point{0.5, 0.5},
	}
}
