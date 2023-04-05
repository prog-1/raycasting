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

var c = color.RGBA{R: 255, G: 255, B: 255, A: 255}

type (
	point struct {
		x, y float64
	}
	game struct {
		m      *ebiten.Image
		p      point
		pg     [][]int
		lp, rp point
	}
)

func rotate(p *point, b bool) {
	x, y := p.x, p.y
	if b {
		p.x = x*math.Cos(math.Pi/200) - y*math.Sin(math.Pi/200)
		p.y = x*math.Sin(math.Pi/200) + y*math.Cos(math.Pi/200)
	} else {
		p.x = x*math.Cos(-math.Pi/200) - y*math.Sin(-math.Pi/200)
		p.y = x*math.Sin(-math.Pi/200) + y*math.Cos(-math.Pi/200)
	}

}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		rotate(&g.lp, true)
		rotate(&g.rp, true)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		rotate(&g.lp, false)
		rotate(&g.rp, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if g.pg[((int(g.p.y)+500-38)/50)+1][int(g.p.x+500)/50] == 0 {
			g.p.y += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if g.pg[(int(g.p.y)+500)/50][int(g.p.x+500+38)/50-1] == 0 {
			g.p.x -= 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if g.pg[(int(g.p.y)+500)/50][int(g.p.x+500-38)/50+1] == 0 {
			g.p.x += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if g.pg[(int(g.p.y+38)+500)/50-1][int(g.p.x+500)/50] == 0 {
			g.p.y -= 5
		}
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-g.p.x, -g.p.y)
	screen.DrawImage(g.m, op)
	ebitenutil.DrawCircle(screen, screenWidth/2, screenHeight/2, 10, color.RGBA{0xff, 0xff, 0x00, 0xff})
	pnt := g.lp
	diff := g.rp.x - g.lp.x
	for pnt.x <= g.rp.x {
		ebitenutil.DrawLine(screen, screenWidth/2, screenHeight/2, pnt.x+screenWidth/2, pnt.y+screenHeight/2, color.RGBA{0xff, 0xff, 0x00, 0xff})
		pnt.x += diff * 0.01
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
		p:  point{0, 0},
		pg: pg,
		lp: point{-200, 200},
		rp: point{200, 200},
	}
}
