package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	winTitle     = "Cube"
	screenWidth  = 625
	screenHeight = 625
	cellSize     = 5
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
		showMap       bool
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
		rotate(&g.dir, -math.Pi/100)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		rotate(&g.dir, math.Pi/100)

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.showMap = !g.showMap

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.fisheye = !g.fisheye

	}
	for i := 0; i < 6; i++ {
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.pg[int(g.p.y/25-g.dir.y/8)][int(g.p.x/25-g.dir.x/8)] == 0 {
			g.p.x -= g.dir.x / 8
			g.p.y -= g.dir.y / 8
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.pg[int(g.p.y/25-g.dir.x/8)][int(g.p.x/25+g.dir.y/8)] == 0 {
			g.p.x += g.dir.y / 8
			g.p.y -= g.dir.x / 8
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.pg[int(g.p.y/25+g.dir.x/8)][int(g.p.x/25-g.dir.y/8)] == 0 {
			g.p.x -= g.dir.y / 8
			g.p.y += g.dir.x / 8
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.pg[int(g.p.y/25+g.dir.y/8)][int(g.p.x/25+g.dir.x/8)] == 0 {
			g.p.x += g.dir.x / 8
			g.p.y += g.dir.y / 8
		}
	}

	return nil
}
func drawWalls(screen *ebiten.Image, dist, wall float64, col color.RGBA) {

	ebitenutil.DrawLine(screen, wall, screenHeight/2, wall, (screenHeight/2)+(5000)/dist, col)
	ebitenutil.DrawLine(screen, wall, screenHeight/2, wall, (screenHeight/2)-(5000)/dist, col)

}
func (g *game) Draw(screen *ebiten.Image) {

	gap := 60.0 / (float64(screenWidth) - 1)
	for i, wall := -30.0, 0.0; i <= 30; i, wall = i+gap, wall+1 {
		var side, step point
		var dist float64
		var ray = g.dir
		rotate(&ray, math.Pi/180*(i))
		maP := point{math.Floor(g.p.x), math.Floor(g.p.y)}
		delta := point{math.Abs(1 / ray.x), math.Abs(1 / ray.y)}
		if ray.x < 0 {
			step.x = -1
			side.x = (g.p.x - float64(int(maP.x))) * delta.x
		} else {
			step.x = 1
			side.x = (float64(int(maP.x)) + 1.0 - g.p.x) * delta.x
		}
		if ray.y < 0 {
			step.y = -1
			side.y = (g.p.y - float64(int(maP.y))) * delta.y
		} else {
			step.y = 1
			side.y = (float64(int(maP.y)) + 1.0 - g.p.y) * delta.y
		}
		shadow := true
		for g.pg[int(maP.y)/25][int(maP.x)/25] == 0 {
			if side.x < side.y {
				side.x += delta.x
				maP.x += step.x
				shadow = true
			} else {
				side.y += delta.y
				maP.y += step.y
				shadow = false
			}

		}
		if g.fisheye {
			if shadow {
				dist = side.x - delta.x
			} else {
				dist = side.y - delta.y
			}
		} else {
			//I invented something
			if shadow {
				dist = (side.x - delta.x) * math.Cos(math.Pi/180*(i))
			} else {
				dist = (side.y - delta.y) * math.Cos(math.Pi/180*(i))
			}
		}
		// dist /= 20
		if shadow {
			switch g.pg[int(maP.y/25)][int(maP.x/25)] {
			case 1:
				drawWalls(screen, dist, wall, color.RGBA{0xff, 0xff, 0xff, 0xff})

			case 2:
				drawWalls(screen, dist, wall, color.RGBA{0xff, 0, 0, 0xff})
			case 3:
				drawWalls(screen, dist, wall, color.RGBA{0, 0xff, 0, 0xff})

			case 4:
				drawWalls(screen, dist, wall, color.RGBA{0, 0, 0xff, 0xff})

			case 5:
				drawWalls(screen, dist, wall, color.RGBA{227, 61, 148, 0xff})
			}

		} else {
			switch g.pg[int(maP.y/25)][int(maP.x/25)] {
			case 1:
				drawWalls(screen, dist, wall, color.RGBA{100, 100, 100, 0xff})

			case 2:
				drawWalls(screen, dist, wall, color.RGBA{100, 0, 0, 0xff})
			case 3:
				drawWalls(screen, dist, wall, color.RGBA{0, 100, 0, 0xff})

			case 4:
				drawWalls(screen, dist, wall, color.RGBA{0, 0, 100, 0xff})

			case 5:
				drawWalls(screen, dist, wall, color.RGBA{100, 30, 60, 0xff})
			}
		}
		if g.showMap {
			ebitenutil.DrawLine(screen, g.p.x/cellSize, g.p.y/cellSize, maP.x/cellSize, maP.y/cellSize, color.RGBA{255, 255, 0, 255})
		}

	}
	if g.showMap {
		screen.DrawImage(g.m, nil)
		ebitenutil.DrawCircle(screen, g.p.x/cellSize, g.p.y/cellSize, 3, color.RGBA{0xff, 0xff, 0x00, 0xff})
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
	pg := [][]int{
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
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 0, 1},
		{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}

	return &game{
		m:       DrawBackground(ebiten.NewImage(screenWidth/cellSize, screenHeight/cellSize), pg),
		p:       point{screenWidth / 2, screenHeight / 2},
		pg:      pg,
		dir:     point{1, -1},
		fisheye: true,
	}
}
