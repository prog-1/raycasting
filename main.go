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
	mazeScale =     15
)

type Point struct {
	x, y float64
}

type Game struct {
	width, height int
	maze          [][]int
	player        Point
	playerEyesDir Point
	fov           Point
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
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.maze[int(g.player.y+g.playerEyesDir.y/mazeScale)][int(g.player.x+g.playerEyesDir.x/mazeScale)] == 0 {
		g.player.x = g.player.x + g.playerEyesDir.x/mazeScale
		g.player.y = g.player.y + g.playerEyesDir.y/mazeScale
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.maze[int(g.player.y-g.playerEyesDir.y/mazeScale)][int(g.player.x+g.playerEyesDir.x/mazeScale)] == 0 {
		g.player.x = g.player.x + g.playerEyesDir.x/mazeScale
		g.player.y = g.player.y - g.playerEyesDir.y/mazeScale
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.maze[int(g.player.y-g.playerEyesDir.y/mazeScale)][int(g.player.x-g.playerEyesDir.x/mazeScale)] == 0 {
		g.player.x = g.player.x - g.playerEyesDir.x/mazeScale
		g.player.y = g.player.y - g.playerEyesDir.y/mazeScale
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.maze[int(g.player.y+g.playerEyesDir.y/mazeScale)][int(g.player.x-g.playerEyesDir.x/mazeScale)] == 0 {
		g.player.x = g.player.x - g.playerEyesDir.x/mazeScale
		g.player.y = g.player.y + g.playerEyesDir.y/mazeScale
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
	g.DrawWalls(screen)
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
				vector.DrawFilledRect(screen, float32(x), float32(y), float32(mazeScale), float32(mazeScale), color.White, false)
			}
			x += mazeScale
		}
		x = 0
		y += mazeScale
	}
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	screen.Set(int(g.player.x)*mazeScale, int(g.player.y)*mazeScale, color.RGBA{255, 0, 0, 255})
}

func (g *Game) DrawFov(screen *ebiten.Image) {
	for i := 0.0; i < screenWidth; i++ {
		var dist float64
		cameraX := 2*i/float64(screenWidth) - 1
		ray := Point{g.playerEyesDir.x + g.fov.x*cameraX, g.playerEyesDir.y + g.fov.y*cameraX}
		dot := Point{math.Floor(g.player.x*mazeScale), math.Floor(g.player.y*mazeScale)}
		deltaDistX, deltaDistY := float64(math.Abs(1/ray.x)), float64(math.Abs(1/ray.y))
		stepX, stepY, sideDistX, sideDistY := g.startDists(ray, Point{dot.x/mazeScale,dot.y/mazeScale}, deltaDistX, deltaDistY)

		for g.maze[int(dot.y)/mazeScale][int(dot.x)/mazeScale] == 0 {
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
		vector.StrokeLine(screen, float32(g.player.x*float64(mazeScale)), float32(g.player.y*float64(mazeScale)), float32(dot.x), float32(dot.y), 1, color.RGBA{255, 255, 0, 255}, false)
	}
}

func (g *Game) DrawWalls(screen *ebiten.Image) {
	perpWallDist := g.perpendicul()
	for i := 0.0; i < screenWidth; i++ {
		cameraX := 2*i/float64(screenWidth) - 1
		ray := Point{g.playerEyesDir.x + g.fov.x*cameraX, g.playerEyesDir.y + g.fov.y*cameraX}
		dot := Point{math.Floor(g.player.x), math.Floor(g.player.y)}
		deltaDistX, deltaDistY := float64(math.Abs(1/ray.x)), float64(math.Abs(1/ray.y))
		stepX, stepY, sideDistX, sideDistY := g.startDists(ray, dot, deltaDistX, deltaDistY)
		var wallVert bool
		var wallHeight float64

		for g.maze[int(dot.y)][int(dot.x)] == 0 {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				dot.x += float64(stepX)
				wallVert = false
			} else {
				sideDistY += deltaDistY
				dot.y += float64(stepY)
				wallVert = true
			}
		}
		if !wallVert {
		perpWallDist = sideDistX - deltaDistX
		} else {
			perpWallDist = sideDistY - deltaDistY
		}
		if ebiten.IsKeyPressed(ebiten.KeyF) {
			perpWallDist *= math.Sqrt(1 + math.Pow(math.Sqrt(math.Pow(ray.x, 2)+math.Pow(ray.y, 2)), 2))
		}
			wallHeight = screenHeight /perpWallDist
			drawStart := -wallHeight/ 2 + screenHeight / 2
		if drawStart < 0{
			drawStart = 0
		}
		drawEnd := wallHeight/ 2 + screenHeight / 2
		if drawEnd >= screenHeight{
			drawEnd = screenHeight - 1
		}
		

	vector.StrokeLine(screen, float32(i),float32(drawStart),float32(i),float32(drawEnd),1,color.RGBA{0, 180, 0, 255},false)
	}
}

func (g *Game) perpendicul() float64 {
	var dist float64
	cameraX := 2*(screenWidth/2)/float64(screenWidth) - 1
	ray := Point{g.playerEyesDir.x + g.fov.x*cameraX, g.playerEyesDir.y + g.fov.y*cameraX}
	dot := Point{math.Floor(g.player.x), math.Floor(g.player.y)}
	deltaDistX, deltaDistY := float64(math.Abs(1/ray.x)), float64(math.Abs(1/ray.y))
	stepX, stepY, sideDistX, sideDistY := g.startDists(ray, dot, deltaDistX, deltaDistY)
	for g.maze[int(dot.y)/mazeScale][int(dot.x)/mazeScale] == 0 {
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			dot.x += float64(stepX)
			dist = sideDistX
		} else {
			sideDistY += deltaDistY
			dot.y += float64(stepY)
			dist = sideDistY
		}
	}
	return dist
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
		player:        Point{1.5, 1.5},
		playerEyesDir: Point{-1, 0},
		fov:           Point{0, 0.3},
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
