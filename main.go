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
	x, y, z float64
}

type Game struct {
	width, height int
	maze          [][]int
	mazeScale     int
	player        Point
	playerEyesDir Point
	fov           Point
}

func sum(a, b Point) Point {
	return Point{a.x + b.x, a.y + b.y, a.z + b.z}
}

func sub(a, b Point) Point {
	return Point{a.x - b.x, a.y - b.y, a.z - b.z}
}

func rotate(p *Point, angle float64) {
	p.x = p.x*math.Cos(angle) - p.y*math.Sin(angle)
	p.y = p.x*math.Sin(angle) + p.y*math.Cos(angle)
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
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
		rotate(&g.playerEyesDir, math.Pi/60)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		rotate(&g.playerEyesDir, math.Pi/-60)
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
				vector.DrawFilledRect(screen, float32(x), float32(y), float32(g.mazeScale), float32(g.mazeScale), color.RGBA{128, 128, 128, 50}, false)
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
	a, b := sub(g.playerEyesDir, g.fov), sum(g.playerEyesDir, g.fov)
	for i := 0; i < 100; i++ {
		ray := Point{g.player.x, g.player.y, 0}
		for int(ray.x) > 1*g.mazeScale && int(ray.x) < 9*g.mazeScale && int(ray.y) > 1*g.mazeScale && int(ray.y) < 9*g.mazeScale {
			ray.x, ray.y = ray.x+a.x, ray.y+a.y
		}
		vector.StrokeLine(screen, float32(g.player.x), float32(g.player.y), float32(ray.x), float32(ray.y), 1, color.RGBA{255, 255, 0, 255}, false)
		a.x, a.y = a.x+b.x/100, a.y+b.y/100
	}
}

func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		maze: [][]int{
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 0, 0, 1, 1, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 1, 0, 0, 1, 1, 0, 0, 1},
			{1, 0, 1, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 1, 1, 1, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		mazeScale:     25,
		player:        Point{25, 25, 0},
		playerEyesDir: Point{-1, 0, 0},
		fov:           Point{0.2, 0.2, 0},
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
