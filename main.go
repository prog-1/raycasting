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
	screenWidth  = 960
	screenHeight = 720
	cellSize     = 8
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
		fisheye       bool
	}
)

func rotate(p *point, angle float64) {
	x, y := p.x, p.y

	p.x = (x*math.Cos(angle) - y*math.Sin(angle))
	p.y = (x*math.Sin(angle) + y*math.Cos(angle))

}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		rotate(&g.dir, -math.Pi/180)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		rotate(&g.dir, math.Pi/180)

	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.pg[int(g.p.y/8-g.dir.y/8)][int(g.p.x/8-g.dir.x/8)] == 0 {
		g.p.x -= g.dir.x / 8
		g.p.y -= g.dir.y / 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.pg[int(g.p.y/8-g.dir.x/8)][int(g.p.x/8+g.dir.y/8)] == 0 {
		g.p.x += g.dir.y / 8
		g.p.y -= g.dir.x / 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.pg[int(g.p.y/8+g.dir.x/8)][int(g.p.x/8-g.dir.y/8)] == 0 {
		g.p.x -= g.dir.y / 8
		g.p.y += g.dir.x / 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.pg[int(g.p.y/8+g.dir.y/8)][int(g.p.x/8+g.dir.x/8)] == 0 {
		g.p.x += g.dir.x / 8
		g.p.y += g.dir.y / 8
	}

	return nil
}
func drawWalls(screen *ebiten.Image, dist, wall, m float64, col color.RGBA) {

	ebitenutil.DrawLine(screen, wall, screenWidth/2, wall, (screenWidth/2)+(screenWidth/2)*4/dist, col)
	ebitenutil.DrawLine(screen, wall, screenWidth/2, wall, (screenWidth/2)-(screenWidth/2)*4/dist, col)

	for tmp := wall; tmp > wall-m; tmp -= m / 20 {
		ebitenutil.DrawLine(screen, tmp, screenWidth/2, tmp, (screenWidth/2)+(screenWidth/2)*4/dist, col)
		ebitenutil.DrawLine(screen, tmp, screenWidth/2, tmp, (screenWidth/2)-(screenWidth/2)*4/dist, col)
	}

}
func (g *game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.m, nil)
	ebitenutil.DrawCircle(screen, g.p.x, g.p.y, 3, color.RGBA{0xff, 0xff, 0x00, 0xff})
	var wall float64
	m := float64(screenWidth) / 60
	for i := 0.0; i <= 60; i++ {
		var side, step point
		var ray = g.dir
		rotate(&ray, math.Pi/180*(i-30))
		p := point{math.Floor(g.p.x), math.Floor(g.p.y)}
		delta := point{math.Abs(1 / ray.x), math.Abs(1 / ray.y)}
		if ray.x < 0 {
			step.x = -1
			side.x = (g.p.x - float64(int(p.x))) * delta.x
		} else {
			step.x = 1
			side.x = (float64(int(p.x)) + 1.0 - g.p.x) * delta.x
		}
		if ray.y < 0 {
			step.y = -1
			side.y = (g.p.y - float64(int(p.y))) * delta.y
		} else {
			step.y = 1
			side.y = (float64(int(p.y)) + 1.0 - g.p.y) * delta.y
		}
		shadow := true
		for g.pg[int(p.y)/cellSize][int(p.x)/cellSize] == 0 {
			if side.x < side.y {
				side.x += delta.x
				p.x += step.x
				shadow = true
			} else {
				side.y += delta.y
				p.y += step.y
				shadow = false
			}
			screen.Set(int(p.x), int(p.y), color.RGBA{255, 255, 0, 255})

		}
		dist := math.Sqrt(math.Pow(p.x/8*25-g.p.x/8*25, 2) + math.Pow(p.y/8*25-g.p.y/8*25, 2))
		if side.x < side.y {
			side.x += delta.x
			p.x += step.x
		} else {
			side.y += delta.y
			p.y += step.y
		}
		if shadow {
			switch g.pg[int(p.y/8)][int(p.x/8)] {
			case 1:
				drawWalls(screen, dist, wall, m, color.RGBA{0xff, 0xff, 0xff, 0xff})

			case 2:
				drawWalls(screen, dist, wall, m, color.RGBA{0xff, 0, 0, 0xff})
			case 3:
				drawWalls(screen, dist, wall, m, color.RGBA{0, 0xff, 0, 0xff})

			case 4:
				drawWalls(screen, dist, wall, m, color.RGBA{0, 0, 0xff, 0xff})

			case 5:
				drawWalls(screen, dist, wall, m, color.RGBA{227, 61, 148, 0xff})
			}

		} else {
			switch g.pg[int(p.y/8)][int(p.x/8)] {
			case 1:
				drawWalls(screen, dist, wall, m, color.RGBA{100, 100, 100, 0xff})

			case 2:
				drawWalls(screen, dist, wall, m, color.RGBA{100, 0, 0, 0xff})
			case 3:
				drawWalls(screen, dist, wall, m, color.RGBA{0, 100, 0, 0xff})

			case 4:
				drawWalls(screen, dist, wall, m, color.RGBA{0, 0, 100, 0xff})

			case 5:
				drawWalls(screen, dist, wall, m, color.RGBA{100, 30, 60, 0xff})
			}
		}

		wall += m

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
			switch pg[i][j] {
			case 1:
				ebitenutil.DrawRect(m, cellSize*cntW, cellSize*cntH, cellSize, cellSize, color.RGBA{0xff, 0xff, 0xff, 0xff})
			case 2:
				ebitenutil.DrawRect(m, cellSize*cntW, cellSize*cntH, cellSize, cellSize, color.RGBA{0xff, 0, 0, 0xff})
			case 3:
				ebitenutil.DrawRect(m, cellSize*cntW, cellSize*cntH, cellSize, cellSize, color.RGBA{0, 0xff, 0, 0xff})
			case 4:
				ebitenutil.DrawRect(m, cellSize*cntW, cellSize*cntH, cellSize, cellSize, color.RGBA{0, 0, 0xff, 0xff})
			case 5:
				ebitenutil.DrawRect(m, cellSize*cntW, cellSize*cntH, cellSize, cellSize, color.RGBA{227, 61, 148, 0xff})
			}
			cntW++

		}
	}
	return m
}
func NewGame() *game {
	pg := [][]int{{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
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
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 0, 1},
		{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}

	return &game{
		m:       DrawBackground(ebiten.NewImage(200, 200), pg),
		p:       point{100, 100},
		pg:      pg,
		dir:     point{1, -1},
		plane:   point{0.5, 0.5},
		fisheye: true,
	}
}
