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
	width  = 640.0
	height = 640.0
)

var (
	maze = [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1},
		{1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
	oneBlockWidthLength  = width / len(maze[0])
	oneBlockHeightLength = height / len(maze)
	FOV                  = 30.0
)

type Point struct {
	x, y float64
}

type Game struct {
	pos     Point
	dir     Point
	showMap bool
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.showMap = !g.showMap
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		a := int(g.pos.x+g.dir.x) / int(oneBlockWidthLength)
		b := int(g.pos.y+g.dir.y) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.x += g.dir.x
			g.pos.y += g.dir.y
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		a := int(g.pos.x-g.dir.x) / int(oneBlockWidthLength)
		b := int(g.pos.y-g.dir.y) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.x -= g.dir.x
			g.pos.y -= g.dir.y
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		tmp := Rotate(g.dir, -math.Pi/2)
		a := int(g.pos.x+tmp.x) / int(oneBlockWidthLength)
		b := int(g.pos.y+tmp.y) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.x += tmp.x
			g.pos.y += tmp.y
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		tmp := Rotate(g.dir, math.Pi/2)
		a := int(g.pos.x+tmp.x) / int(oneBlockWidthLength)
		b := int(g.pos.y+tmp.y) / int(oneBlockHeightLength)
		if maze[b][a] == 0 {
			g.pos.x += tmp.x
			g.pos.y += tmp.y
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dir = Rotate(g.dir, -math.Pi/140)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dir = Rotate(g.dir, math.Pi/140)
	}

	return nil
}
func Rotate(a Point, angle float64) Point {
	a.x, a.y = a.x*math.Cos(angle)-a.y*math.Sin(angle), a.x*math.Sin(angle)+a.y*math.Cos(angle)
	return a
}

func DrawMap(screen *ebiten.Image) {
	for i := range maze {
		for j, WallType := range maze[i] {
			if WallType != 0 {
				ebitenutil.DrawRect(screen, float64(j*oneBlockWidthLength), float64(i*oneBlockHeightLength), float64(oneBlockWidthLength), float64(oneBlockHeightLength), color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func abs(a float64) float64 {
	if a < 0 {
		return a * -1
	}
	return a
}

func raycast(screen *ebiten.Image, pos Point, dirx, diry float64, draw bool) (int, float64) {
	mapX, mapY := int(pos.x), int(pos.y)

	deltaDistX := abs(1 / dirx)
	deltaDistY := abs(1 / diry)
	var sideDistX, sideDistY float64
	var stepX, stepY int
	if dirx < 0 {
		stepX = -1
		sideDistX = (pos.x - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX) + 1.0 - pos.x) * deltaDistX
	}
	if diry < 0 {
		stepY = -1
		sideDistY = (pos.y - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY) + 1.0 - pos.y) * deltaDistY
	}
	var side int
	for {
		//jump to next map square, either in x-direction, or in y-direction
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepX
			side = 0
		} else {
			sideDistY += deltaDistY
			mapY += stepY
			side = 1
		}
		//Check if ray has hit a wall
		// fmt.Println(mapX, mapY)
		if maze[mapY/int(oneBlockHeightLength)][mapX/int(oneBlockWidthLength)] > 0 {
			if draw {
				ebitenutil.DrawLine(screen, pos.x, pos.y, float64(mapX), float64(mapY), color.RGBA{255, 255, 0, 255})
			}
			return side, math.Sqrt(math.Pow(float64(mapX)-pos.x, 2) + math.Pow(float64(mapY)-pos.y, 2))
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, "Hello, World!")
	if g.showMap {
		DrawMap(screen)
		ebitenutil.DrawCircle(screen, g.pos.x, g.pos.y, 3, color.RGBA{255, 255, 0, 255})
	}
	// fmt.Println(screen.Size())
	for i, line := -FOV, 0.0; i <= FOV; i, line = i+2*FOV/(width-1), line+1 {
		c := color.RGBA{91, 97, 89, 250}
		tmp := Rotate(g.dir, i*math.Pi/180.0)
		// ebitenutil.DrawLine(screen, g.pos.x, g.pos.y, g.pos.x+tmp.x*1000, g.pos.y+tmp.y*1000, color.RGBA{255, 255, 0, 255})
		side, distance := raycast(screen, g.pos, tmp.x, tmp.y, g.showMap)
		if side == 1 {
			c = color.RGBA{155, 161, 153, 255}
		}
		ebitenutil.DrawLine(screen, line, height/2, line, (height/2)+(height/2)/distance, c)
		ebitenutil.DrawLine(screen, line, height/2, line, (height/2)-(height/2)/distance, c)

	}
	// ebitenutil.DrawRect(screen, float64(a*oneBlockWidthLength), float64(b*oneBlockHeightLength), float64(oneBlockWidthLength), float64(oneBlockHeightLength), color.RGBA{0, 0, 255, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{Point{320, 350}, Point{0, -1}, false}); err != nil {
		log.Fatal(err)
	}
}
