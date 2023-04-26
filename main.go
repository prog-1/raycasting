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
	fisheye             bool
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
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.fisheye = !g.fisheye
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

func (g *Game) DrawWalls(screen *ebiten.Image) {
	for i := 0; i < screenWidth; i++ {
		cameraX := 2*float64(i)/float64(screenWidth) - 1
		rayDirX, rayDirY := g.dir.x+g.fov.x*cameraX, g.dir.y+g.fov.y*cameraX
		mapX, mapY := g.playerPos.x, g.playerPos.y
		deltaDistX, deltaDistY := math.Abs(1/rayDirX), math.Abs(1/rayDirY)
		var sideDistX, sideDistY, stepX, stepY, perpWallDist float64
		var side int
		hit := 0
		for hit == 0 {
			if rayDirX < 0 {
				stepX = -1
				sideDistX = (g.playerPos.x - mapX) * deltaDistX
			} else {
				stepX = 1
				sideDistX = (mapX + 1.0 - g.playerPos.x) * deltaDistX
			}
			if rayDirY < 0 {
				stepY = -1
				sideDistY = (g.playerPos.y - mapY) * deltaDistY
			} else {
				stepY = 1
				sideDistY = (mapY + 1.0 - g.playerPos.y) * deltaDistY
			}

			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				side = 0
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				side = 1
			}
			if g.gameMap[int(mapY)][int(mapX)] > 0 {
				hit = 1
			}
		}
		if side == 0 {
			perpWallDist = sideDistX - deltaDistX
		} else {
			perpWallDist = sideDistY - deltaDistY
		}
		if g.fisheye {
			perpWallDist *= math.Sqrt(1 + math.Pow(math.Sqrt(math.Pow(rayDirX, 2)+math.Pow(rayDirY, 2)), 2))
		}
		lineHeight := int(float64(screenHeight) / perpWallDist)
		drawStart := -lineHeight/2 + screenHeight/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + screenHeight/2
		if drawEnd >= screenHeight {
			drawEnd = screenHeight - 1
		}
		c := g.colors[g.gameMap[int(mapY)][int(mapX)]-1]
		if side == 1 {
			r, g, b, a := c.RGBA()
			c = color.RGBA{uint8(float64(r>>8) * 0.3), uint8(float64(g>>8) * 0.3), uint8(float64(b>>8) * 0.3), uint8(float64(a >> 8))}
		}
		vector.StrokeLine(screen, float32(screenWidth-i), float32(drawStart), float32(screenWidth-i), float32(drawEnd), 3, c, false)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawWalls(screen)
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
		fov:       Point{0.5, 0.5},
		gameMap: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 5, 0, 0, 0, 2, 0, 2, 2, 0, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2, 0, 3, 3, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 4, 0, 4, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 4, 0, 0, 4, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 5, 5, 5, 5, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 5, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 2, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		colors: []color.Color{color.RGBA{255, 255, 255, 150},
			color.RGBA{0, 255, 0, 150},
			color.RGBA{255, 0, 0, 150},
			color.RGBA{0, 66, 255, 255},
			color.RGBA{255, 0, 132, 255}},
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
