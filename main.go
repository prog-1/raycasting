package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.maze[int(g.player.y-1)/g.mazeScale][int(g.player.x)/g.mazeScale] == 0 {
		g.player.y--
		g.playerEyesDir.y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.maze[int(g.player.y)/g.mazeScale][int(g.player.x-1)/g.mazeScale] == 0 {
		g.player.x--
		g.playerEyesDir.x--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.maze[int(g.player.y+1)/g.mazeScale][int(g.player.x)/g.mazeScale] == 0 {
		g.player.y++
		g.playerEyesDir.y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.maze[int(g.player.y)/g.mazeScale][int(g.player.x+1)/g.mazeScale] == 0 {
		g.player.x++
		g.playerEyesDir.x++
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
	for n := 10; n > 0; n-- {
	}
	ebitenutil.DrawLine(screen, g.player.x, g.player.y, a.x, a.y, color.RGBA{255, 255, 0, 255})
	ebitenutil.DrawLine(screen, g.player.x, g.player.y, b.x, b.y, color.RGBA{255, 255, 0, 255})

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
		playerEyesDir: Point{25, 0, 0},
		fov:           Point{5, 0, 0},
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
