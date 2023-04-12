package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 960
	screenHeight = 720
	cellSize     = 8
	rayNum       = 100
	texWidth     = 64
	texHeight    = 64
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
	textures        [5]*ebiten.Image
}

func NewGame(width, height int) *Game {
	return &Game{
		width:    width,
		height:   height,
		pos:      &Point{13 * cellSize, 13 * cellSize},
		dir:      &Point{-1, 0},
		plane:    &Point{0, 0.5},
		textures: loadTextures(),
	}
}

func loadTextures() (textures [5]*ebiten.Image) {
	for i := 1; i <= len(textures); i++ {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("texture%d.png", i))
		if err != nil {
			panic(err)
		}
		textures[i-1] = img
	}
	return
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func rotate(p *Point, angle float64) *Point {
	res := new(Point)
	res.x = p.x*math.Cos(angle) - p.y*math.Sin(angle)
	res.y = p.x*math.Sin(angle) + p.y*math.Cos(angle)
	return res
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) && worldMap[int(g.pos.y/cellSize+g.dir.y/2)][int(g.pos.x/cellSize+g.dir.x/2)] == 0 {
		g.pos.x += g.dir.x / 2
		g.pos.y += g.dir.y / 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && worldMap[int(g.pos.y/cellSize-g.dir.y/2)][int(g.pos.x/cellSize-g.dir.x/2)] == 0 {
		g.pos.x -= g.dir.x / 2
		g.pos.y -= g.dir.y / 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && worldMap[int(g.pos.y/cellSize-g.dir.x/2)][int(g.pos.x/cellSize+g.dir.y/2)] == 0 {
		g.pos.x += g.dir.y / 2
		g.pos.y -= g.dir.x / 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && worldMap[int(g.pos.y/cellSize+g.dir.x/2)][int(g.pos.x/cellSize-g.dir.y/2)] == 0 {
		g.pos.x -= g.dir.y / 2
		g.pos.y += g.dir.x / 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dir = rotate(g.dir, -math.Pi/90)
		g.plane = rotate(g.plane, -math.Pi/90)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dir = rotate(g.dir, math.Pi/90)
		g.plane = rotate(g.plane, math.Pi/90)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for x := 0; x < g.width; x++ {
		cameraX := 2*float64(x)/float64(g.width) - 1
		rayDir := Point{g.dir.x + g.plane.x*cameraX, g.dir.y + g.plane.y*cameraX}
		p := Point{g.pos.x, g.pos.y}
		deltaDist := Point{math.Abs(1/rayDir.x) / cellSize, math.Abs(1/rayDir.y) / cellSize}
		var sideDist, step Point
		var perpWallDist float64
		var side int
		if rayDir.x < 0 {
			step.x = -1
			sideDist.x = (g.pos.x - p.x) / cellSize * deltaDist.x
		} else {
			step.x = 1
			sideDist.x = (p.x + 1.0 - g.pos.x) / cellSize * deltaDist.x
		}
		if rayDir.y < 0 {
			step.y = -1
			sideDist.y = (g.pos.y - p.y) / cellSize * deltaDist.y
		} else {
			step.y = 1
			sideDist.y = (p.y + 1.0 - g.pos.y) / cellSize * deltaDist.y
		}
		for worldMap[int(p.y)/cellSize][int(p.x)/cellSize] == 0 {
			if sideDist.x < sideDist.y {
				sideDist.x += deltaDist.x
				p.x += step.x
				side = 0
			} else {
				sideDist.y += deltaDist.y
				p.y += step.y
				side = 1
			}
		}
		if side == 0 {
			perpWallDist = sideDist.x - deltaDist.x
		} else {
			perpWallDist = sideDist.y - deltaDist.y
		}
		lineHeight := int(float64(g.height) / perpWallDist)
		drawStart := -lineHeight/2 + g.height/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + g.height/2
		if drawEnd >= g.height {
			drawEnd = g.height - 1
		}
		texNum := worldMap[int(p.y)/cellSize][int(p.x)/cellSize] - 1
		var wallX float64
		if side == 0 {
			wallX = g.pos.y/cellSize + perpWallDist*rayDir.y
		} else {
			wallX = g.pos.x/cellSize + perpWallDist*rayDir.x
		}
		wallX -= math.Floor(wallX)
		texX := int(wallX * float64(texWidth))
		if (side == 0 && rayDir.x > 0) || (side == 1 && rayDir.y < 0) {
			texX = texWidth - texX - 1
		}
		texStep := float64(1.0 * texHeight / lineHeight)
		texPos := float64(drawStart-g.height/2+lineHeight/2) * texStep
		for y := drawStart; y < drawEnd; y++ {
			texY := int(texPos) & (texHeight - 1)
			texPos += texStep
			c := g.textures[texNum].At(texX, texY)
			if side == 1 {
				r, g, b, a := c.RGBA()
				c = color.RGBA{uint8(float64(r>>8) * 0.5), uint8(float64(g>>8) * 0.5), uint8(float64(b>>8) * 0.5), uint8(float64(a >> 8))}
			}
			screen.Set(g.width-x, y, c)
		}
	}
	for i := range worldMap {
		for j := range worldMap[i] {
			if worldMap[j][i] != 0 {
				vector.DrawFilledRect(screen, float32(i*cellSize), float32(j*cellSize), cellSize, cellSize, clr[worldMap[j][i]-1], false)
			}
		}
	}
	vector.DrawFilledCircle(screen, float32(g.pos.x), float32(g.pos.y), 3, color.RGBA{255, 255, 0, 150}, false)
	for i := 0.0; i < rayNum; i++ {
		cameraX := 2*i/float64(rayNum) - 1
		rayDir := Point{g.dir.x + g.plane.x*cameraX, g.dir.y + g.plane.y*cameraX}
		p := Point{g.pos.x, g.pos.y}
		deltaDist := Point{math.Abs(1 / rayDir.x), math.Abs(1 / rayDir.y)}
		var sideDist, step Point
		if rayDir.x < 0 {
			step.x = -1
			sideDist.x = (g.pos.x - p.x) * deltaDist.x
		} else {
			step.x = 1
			sideDist.x = (p.x + 1.0 - g.pos.x) * deltaDist.x
		}
		if rayDir.y < 0 {
			step.y = -1
			sideDist.y = (g.pos.y - p.y) * deltaDist.y
		} else {
			step.y = 1
			sideDist.y = (p.y + 1.0 - g.pos.y) * deltaDist.y
		}
		for worldMap[int(p.y)/cellSize][int(p.x)/cellSize] == 0 {
			if sideDist.x < sideDist.y {
				sideDist.x += deltaDist.x
				p.x += step.x
			} else {
				sideDist.y += deltaDist.y
				p.y += step.y
			}
		}
		vector.StrokeLine(screen, float32(g.pos.x), float32(g.pos.y), float32(p.x), float32(p.y), 3, color.RGBA{255, 255, 0, 50}, false)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
