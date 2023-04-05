package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	winTitle     = "Raycasting"
	screenWidth  = 1000
	screenHeight = 1000
	dpi          = 100
)

type (
	point struct {
		x, y float64
	}
	game struct {
		m  *ebiten.Image
		p  point
		pg [][]int
		lp point
	}
)

func rotate(p *point, angle float64) {
	x, y := p.x-500, p.y-500

	p.x = (x*math.Cos(angle) - y*math.Sin(angle)) + 500
	p.y = (x*math.Sin(angle) + y*math.Cos(angle)) + 500

}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }

func (g *game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		rotate(&g.lp, -math.Pi/50)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		rotate(&g.lp, math.Pi/50)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.pg[((int(g.p.y)+500-38)/50)+1][int(g.p.x+500)/50] == 0 {
			g.p.y += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.pg[(int(g.p.y)+500)/50][int(g.p.x+500+38)/50-1] == 0 {
			g.p.x -= 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.pg[(int(g.p.y)+500)/50][int(g.p.x+500-38)/50+1] == 0 {
			g.p.x += 5
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
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
	for i := 0; i < 180; i++ {
		ebitenutil.DrawLine(screen, screenWidth/2, screenHeight/2, pnt.x, pnt.y, color.RGBA{0xff, 0xff, 0x00, 0xff})
		rotate(&pnt, math.Pi/360)
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
		{1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	return &game{
		m:  DrawBackground(ebiten.NewImage(screenWidth, screenHeight), pg),
		p:  point{0, 0},
		pg: pg,
		lp: point{300, 700},
		// rp: point{700, 300},
	}
}
