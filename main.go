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

func GetTheSign(v float64) float64 { // returns -1 or 1
	if math.Ceil(v) == 0 {
		return math.Floor(v)
	}
	return math.Ceil(v)
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
		kx, ky := ray.y/ray.x, ray.x/ray.y
		prevx, prevy := int(pos.x)-int(pos.x)%15, int(pos.y)-int(pos.y)%15
		A, B := pos.x-float64(prevx), pos.y-float64(prevy)
		if ray.x > 0 {
			A = float64(prevx) + 15 - pos.x
		}
		if ray.y > 0 {
			B = float64(prevy) + 15 - pos.y
		}
		lx, ly := math.Sqrt(math.Pow(A*kx, 2)+math.Pow(A, 2)), math.Sqrt(math.Pow(B*ky, 2)+math.Pow(B, 2))
		startx, starty := Point{pos.x + A*GetTheSign(ray.x), pos.y + A*math.Abs(kx)*GetTheSign(ray.y)}, Point{pos.x + B*math.Abs(ky)*GetTheSign(ray.x), pos.y + B*GetTheSign(ray.y)}
		for startx.x >= 0 && startx.y >= 0 && starty.x >= 0 && starty.y >= 0 && Map[int(startx.y/15)][int(startx.x/15)] == 0 && Map[int(starty.y/15)][int(starty.x/15)] == 0 {
			if lx > ly {
				startx.x, startx.y = startx.x+15*GetTheSign(ray.x), startx.y+15*math.Abs(kx)*GetTheSign(ray.y)
				lx += math.Sqrt(math.Pow(15, 2) + math.Pow(15*kx, 2))
			} else {
				starty.x, starty.y = starty.x+15*math.Abs(ky)*GetTheSign(ray.x), starty.y+15*GetTheSign(ray.y)
				ly += math.Sqrt(math.Pow(15, 2) + math.Pow(15*ky, 2))
			}
			// fmt.Println("startx:", i, startx.x, startx.y)
			// fmt.Println("starty:", i, starty.x, starty.y)
		}
		if lx > ly {
			ebitenutil.DrawLine(img, pos.x, pos.y, startx.x, startx.y, color.RGBA{0xf8, 0xf0, 0x00, 0xff})
		} else {
			ebitenutil.DrawLine(img, pos.x, pos.y, starty.x, starty.y, color.RGBA{0xf8, 0xf0, 0x00, 0xff})
		}
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
