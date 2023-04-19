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
	screenWidth  = 960
	screenHeight = 640
)

type Point struct {
	x, y float64
}

type Game struct {
	width, height int
	maze          [][]int
	wallColors    map[int]color.RGBA
	mazeScale     int
	player        Point
	playerEyesDir Point
	fov           Point
}

func sum(a, b Point) Point {
	return Point{a.x + b.x, a.y + b.y}
}

func sub(a, b Point) Point {
	return Point{a.x - b.x, a.y - b.y}
}

func rotate(p Point, angle float64) Point {
	return Point{p.x*math.Cos(angle) - p.y*math.Sin(angle), p.x*math.Sin(angle) + p.y*math.Cos(angle)}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) startDists(ray, dot Point, deltaX, deltaY float64) (stepX int, stepY int, sideX float64, sideY float64) {
	if ray.x < 0 {
		stepX = -1
		sideX = (g.player.x - dot.x) * deltaX
	} else {
		stepX = 1
		sideX = (dot.x + 1.0 - g.player.x) * deltaX
	}
	if ray.y < 0 {
		stepY = -1
		sideY = (g.player.y - dot.y) * deltaY
	} else {
		stepY = 1
		sideY = (dot.y + 1.0 - g.player.y) * deltaY
	}
	return stepX, stepY, sideX, sideY
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.maze[int(g.player.y+g.playerEyesDir.y)/g.mazeScale][int(g.player.x+g.playerEyesDir.x)/g.mazeScale] == 0 {
		g.player.x = g.player.x + g.playerEyesDir.x
		g.player.y = g.player.y + g.playerEyesDir.y
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.maze[int(g.player.y-g.playerEyesDir.y)/g.mazeScale][int(g.player.x+g.playerEyesDir.x)/g.mazeScale] == 0 {
		g.player.x = g.player.x + g.playerEyesDir.x
		g.player.y = g.player.y - g.playerEyesDir.y
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.maze[int(g.player.y-g.playerEyesDir.y)/g.mazeScale][int(g.player.x-g.playerEyesDir.x)/g.mazeScale] == 0 {
		g.player.x = g.player.x - g.playerEyesDir.x
		g.player.y = g.player.y - g.playerEyesDir.y
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.maze[int(g.player.y+g.playerEyesDir.y)/g.mazeScale][int(g.player.x-g.playerEyesDir.x)/g.mazeScale] == 0 {
		g.player.x = g.player.x - g.playerEyesDir.x
		g.player.y = g.player.y + g.playerEyesDir.y
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.playerEyesDir = rotate(g.playerEyesDir, math.Pi/60)
		g.fov = rotate(g.fov, math.Pi/60)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.playerEyesDir = rotate(g.playerEyesDir, math.Pi/-60)
		g.fov = rotate(g.fov, math.Pi/-60)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// if ebiten.IsKeyPressed(ebiten.KeyM) {
	g.DrawMinimap(screen)
	g.DrawPlayer(screen)
	g.DrawFov(screen)
	// }
}

func (g *Game) DrawMinimap(screen *ebiten.Image) {
	var x, y int
	for i := range g.maze {
		for _, j := range g.maze[i] {
			if j > 0 {
				vector.DrawFilledRect(screen, float32(x), float32(y), float32(g.mazeScale), float32(g.mazeScale), g.wallColors[j], false)
			}
			x += g.mazeScale
		}
		x = 0
		y += g.mazeScale
	}
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	screen.Set(int(g.player.x), int(g.player.y), color.RGBA{255, 0, 0, 255})
}

func (g *Game) DrawFov(screen *ebiten.Image) {
	for i := 0.0; i < screenWidth; i++ {
		var dist float64
		cameraX := 2*i/float64(screenWidth) - 1
		ray := Point{g.playerEyesDir.x + g.fov.x*cameraX, g.playerEyesDir.y + g.fov.y*cameraX}
		dot := Point{math.Floor(g.player.x), math.Floor(g.player.y)}
		deltaDistX, deltaDistY := float64(math.Abs(1/ray.x)), float64(math.Abs(1/ray.y))
		stepX, stepY, sideDistX, sideDistY := g.startDists(ray, dot, deltaDistX, deltaDistY)

		for g.maze[int(dot.y)/g.mazeScale][int(dot.x)/g.mazeScale] == 0 {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				dot.x += float64(stepX)
				dist += deltaDistX
			} else {
				sideDistY += deltaDistY
				dot.y += float64(stepY)
				dist += deltaDistY
			}
		}
		vector.StrokeLine(screen, float32(g.player.x), float32(g.player.y), float32(dot.x), float32(dot.y), 1, color.RGBA{255, 255, 0, 255}, false)
		wh := screenHeight / dist
		vector.StrokeLine(screen, float32(i), float32(screenHeight/2), float32(i), float32(screenHeight/2+wh/2), 1, color.RGBA{255, 255, 0, 200}, false)
		vector.StrokeLine(screen, float32(i), float32(screenHeight/2), float32(i), float32(screenHeight/2-wh/2), 1, color.RGBA{255, 255, 0, 200}, false)
	}
}

func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		maze: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 2, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 1},
			{1, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 4, 4, 4, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 4, 1},
			{1, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4, 4, 4, 4, 4, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		wallColors: map[int]color.RGBA{
			1: color.RGBA{255, 255, 255, 255},
			2: color.RGBA{200, 10, 10, 255},
			3: color.RGBA{10, 200, 10, 255},
			4: color.RGBA{10, 10, 200, 255},
		},
		mazeScale:     15,
		player:        Point{20, 20},
		playerEyesDir: Point{-1, 0},
		fov:           Point{0.3, 0.3},
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
