package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/*
--------------------STRUCTS--------------------
*/
const (
	winTitle     = "Raycasting"
	screenWidth  = 625
	screenHeight = 625
	cellSize     = 5
	dpi          = 100
)

// var (
// 	clr  = [5]color.RGBA{{255, 255, 255, 0xff}, {255, 0, 0, 0xff}, {0, 255, 0, 0xff}, {0, 0, 255, 0xff}, {0, 255, 255, 0xff}}
// 	clr2 = [5]color.RGBA{{50, 50, 50, 0xff}, {50, 0, 0, 0xff}, {0, 50, 0, 0xff}, {0, 0, 50, 0xff}, {50, 0, 20, 0xff}}
// )

type (
	point struct {
		x, y float64
	}
	game struct {
		m        *ebiten.Image
		p, dir   point
		pg       [][]int
		fisheye  bool
		showMap  bool
		textures [5][]color.Color
		buffer   [screenHeight][screenWidth]color.Color
	}
)

/*
--------------------STARTING GAME--------------------
*/
func loadTextures() (textures [5][]color.Color) {
	for i := 0; i < len(textures); i++ {
		textures[i] = make([]color.Color, 64*64)
		_, img, err := ebitenutil.NewImageFromFile(fmt.Sprintf("texture%d.png", i+1))
		if err != nil {
			panic(err)
		}
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				c := img.At(x, y)
				textures[i][64*y+x] = c
			}
		}
	}
	return
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
		m:        DrawBackground(ebiten.NewImage(screenWidth/cellSize, screenHeight/cellSize), pg),
		p:        point{screenWidth / 2, screenHeight / 2},
		pg:       pg,
		dir:      point{1, -1},
		fisheye:  true,
		textures: loadTextures(),
	}
}

/*
--------------------LAYOUT AND UPDATE--------------------
*/

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }

func (g *game) Update() error {
	/*
	   --------------------rotation and toggling fisheye & map--------------------
	*/

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
	/*
	   --------------------moving in 4 directions--------------------
	*/
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

/*
--------------------DRAW FUNCTION--------------------
*/

// func drawWalls(screen *ebiten.Image, dist, wall float64, col color.RGBA) {

// 	ebitenutil.DrawLine(screen, wall, screenHeight/2, wall, (screenHeight/2)+(5000)/dist, col)
// 	ebitenutil.DrawLine(screen, wall, screenHeight/2, wall, (screenHeight/2)-(5000)/dist, col)

// }

func (g *game) Draw(screen *ebiten.Image) {

	gap := 60.0 / (float64(screenWidth) - 1)
	for i, x := 30.0, float64(screenWidth)-1; i >= -30; i, x = i-gap, x-1 {
		/*
		   --------------------raycasting--------------------
		*/
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
		/*
		   --------------------fisheye--------------------
		*/
		if g.fisheye {
			if shadow {
				dist = side.x - delta.x
			} else {
				dist = side.y - delta.y
			}
		} else {
			if shadow {
				dist = (side.x - delta.x) * math.Cos(math.Pi/180*(i))
			} else {
				dist = (side.y - delta.y) * math.Cos(math.Pi/180*(i))
			}
		}
		dist /= 16
		// if shadow {
		// 	drawWalls(screen, dist, wall, clr[g.pg[int(maP.y/25)][int(maP.x/25)]])
		// } else {
		// 	drawWalls(screen, dist, wall, clr2[g.pg[int(maP.y/25)][int(maP.x/25)]])
		// }

		/*
		   --------------------textured walls--------------------
		*/

		lh := int(float64(screenHeight) / dist)
		Start := -lh/2 + screenHeight/2
		if Start < 0 {
			Start = 0
		}
		end := lh/2 + screenHeight/2
		if end >= screenHeight {
			end = screenHeight - 1
		}
		texNum := g.pg[int(maP.y/25)][int(maP.x)/25] - 1
		var wallX float64
		if shadow {
			wallX = g.p.y + dist*ray.y
		} else {
			wallX = g.p.x + dist*ray.x
		}
		wallX -= math.Floor(wallX)
		texX := int(wallX * 64.0)
		if (shadow && ray.x > 0) || (!shadow && ray.y < 0) {
			texX = 64.0 - texX - 1
		}
		texStep := float64(64) / float64(lh)
		texPos := float64(Start-screenHeight/2+lh/2) * texStep
		for y := Start; y < end; y++ {
			texY := int(texPos) & (64 - 1)
			texPos += texStep
			c := g.textures[texNum][64*texY+texX]
			if !shadow {
				r, g, b, a := c.RGBA()
				c = color.RGBA{uint8(float64(r>>8) * 0.5), uint8(float64(g>>8) * 0.5), uint8(float64(b>>8) * 0.5), uint8(float64(a >> 8))}
			}
			g.buffer[y][int(x)] = c
		}
		//drawing vectors on map
		if g.showMap {
			ebitenutil.DrawLine(screen, g.p.x/cellSize, g.p.y/cellSize, maP.x/cellSize, maP.y/cellSize, color.RGBA{255, 255, 0, 255})
		}
	}
	for y := screenHeight - 1; y >= 0; y-- {
		for x := screenWidth - 1; x >= 0; x-- {
			if g.buffer[y][x] != nil {
				screen.Set(x, y, g.buffer[y][x])
				g.buffer[y][x] = nil
			}
		}

	}
	//drawing map in left upper corner.
	if g.showMap {
		screen.DrawImage(g.m, nil)
		ebitenutil.DrawCircle(screen, g.p.x/cellSize, g.p.y/cellSize, 3, color.RGBA{0xff, 0xff, 0x00, 0xff})
	}

}

/*
--------------------MAIN and AUXILIARY--------------------
*/

func main() {
	ebiten.SetWindowTitle(winTitle)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func rotate(p *point, angle float64) {
	x, y := p.x, p.y
	p.x = (x*math.Cos(angle) - y*math.Sin(angle))
	p.y = (x*math.Sin(angle) + y*math.Cos(angle))

}
