package main

import (
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 700
	screenHeight = 700
)

// Point is a struct for representing 2D vectors.
type Point struct {
	x, y float64
}

type Player struct {
	pos Point
	dir Point
}

var Map = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

// NewPlayer initializes and returns a new Player instance.
func NewPlayer() *Player {
	return &Player{
		pos: Point{15 * 12, 15 * 23},
		dir: Point{0, -1},
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) && Map[int(p.pos.y+p.dir.y)/15][int(p.pos.x+p.dir.x)/15] == 0 {
		p.pos.x, p.pos.y = p.pos.x+p.dir.x, p.pos.y+p.dir.y
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && Map[int(p.pos.y-p.dir.y)/15][int(p.pos.x-p.dir.x)/15] == 0 {
		p.pos.x, p.pos.y = p.pos.x-p.dir.x, p.pos.y-p.dir.y
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && Map[int(p.pos.y+p.dir.x*math.Sin(math.Pi/2)+p.dir.y*math.Cos(math.Pi/2))/15][int(p.pos.x+p.dir.x*math.Cos(math.Pi/2)-p.dir.y*math.Sin(math.Pi/2))/15] == 0 {
		p.pos.x, p.pos.y = p.pos.x+p.dir.x*math.Cos(math.Pi/2)-p.dir.y*math.Sin(math.Pi/2), p.pos.y+p.dir.x*math.Sin(math.Pi/2)+p.dir.y*math.Cos(math.Pi/2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && Map[int(p.pos.y+p.dir.x*math.Sin(-math.Pi/2)+p.dir.y*math.Cos(-math.Pi/2))/15][int(p.pos.x+p.dir.x*math.Cos(-math.Pi/2)-p.dir.y*math.Sin(-math.Pi/2))/15] == 0 {
		p.pos.x, p.pos.y = p.pos.x+p.dir.x*math.Cos(-math.Pi/2)-p.dir.y*math.Sin(-math.Pi/2), p.pos.y+p.dir.x*math.Sin(-math.Pi/2)+p.dir.y*math.Cos(-math.Pi/2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.dir.x, p.dir.y = p.dir.x*math.Cos(math.Pi/180)-p.dir.y*math.Sin(math.Pi/180), p.dir.x*math.Sin(math.Pi/180)+p.dir.y*math.Cos(math.Pi/180)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.dir.x, p.dir.y = p.dir.x*math.Cos(-math.Pi/180)-p.dir.y*math.Sin(-math.Pi/180), p.dir.x*math.Sin(-math.Pi/180)+p.dir.y*math.Cos(-math.Pi/180)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
}

// Game is a game instance.
type Game struct {
	width, height int
	Player        *Player
	mapIsOpen     bool
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		Player: NewPlayer(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	g.Player.Update()
	return nil
}

func DrawMap(img *ebiten.Image, pos, dir Point) {
	for i := range Map {
		for j := range Map[i] {
			c := color.RGBA{}
			if Map[i][j] == 1 {
				c = color.RGBA{0xff, 0xff, 0xff, 0xff}
			} else if Map[i][j] == 2 {
				c = color.RGBA{0x0E, 0x30, 0x7F, 0xff}
			} else if Map[i][j] == 3 {
				c = color.RGBA{0x0A, 0x51, 0x0C, 0xff}
			} else if Map[i][j] == 4 {
				c = color.RGBA{0x78, 0x0A, 0x0A, 0xff}
			} else if Map[i][j] == 5 {
				c = color.RGBA{0x96, 0x1E, 0x62, 0xff}
			}
			ebitenutil.DrawRect(img, float64(15*j), float64(15*i), 15, 15, c)
		}
	}
	ebitenutil.DrawCircle(img, pos.x, pos.y, 3, color.White)
	// 320 rays for 60 degrees
	startdir := Point{dir.x*math.Cos(-math.Pi/6) - dir.y*math.Sin(-math.Pi/6), dir.x*math.Sin(-math.Pi/6) + dir.y*math.Cos(-math.Pi/6)}
	for i := 0; i <= 319; i++ {
		ray := Point{startdir.x*math.Cos(float64(i)*60/320*math.Pi/180) - startdir.y*math.Sin(float64(i)*60/320*math.Pi/180), startdir.x*math.Sin(float64(i)*60/320*math.Pi/180) + startdir.y*math.Cos(float64(i)*60/320*math.Pi/180)}
		d := Point{pos.x + ray.x, pos.y + ray.y}
		for ; Map[int(d.y/15)][int(d.x/15)] == 0; d.x, d.y = d.x+ray.x, d.y+ray.y {
		}
		ebitenutil.DrawLine(img, pos.x, pos.y, d.x, d.y, color.RGBA{0xf8, 0xf0, 0x00, 0xff})
	}
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) && !g.mapIsOpen {
		g.mapIsOpen = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.mapIsOpen = false
	}
	if g.mapIsOpen {
		DrawMap(screen, g.Player.pos, g.Player.dir)
	}
	g.Player.Draw(screen)
}

func main() {
	//rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
