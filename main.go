package main

import (
	"image/color"
	"log"
	"math"
	"time"

	v "github.com/34thSchool/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sw, sh = 1920, 1080
	mw, mh = 4, 4                               // map width/height
	cw, ch = float64(sw) / mw, float64(sh) / mh // cell width/height
)

type Pair[T, U any] struct {
	X T
	Y U
}

type Player struct {
	Pos    v.Vec
	Dir    v.Vec
	ms, rs float64 // Movement and Rotation speed
	a, h   int     // base length and height of player's FOV triangle
}

type game struct {
	sb            *ebiten.Image
	maze          *[mh][mw]int
	p             Player
	prevFrameTime int64
	rc            int //Ray Count
}

func NewGame() *game {
	// Values on the map:
	// 0 - empty space
	// 1 - pink wall
	// 2 - white wall
	// maze := [mh][mw]int{
	// 	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	// 	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	// }
	maze := [mh][mw]int{
		{1, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 0, 0, 1},
		{2, 1, 1, 1},
	}
	return &game{
		ebiten.NewImage(sw, sh),
		&maze,
		InitPlayer(maze),
		0,
		101,
	}
}

func InitPlayer(maze [mh][mw]int) Player {
	var x, y int
	for y = range maze {
		for x = range maze[y] {
			if maze[y][x] == 0 {
				break
			}
		}
		if maze[y][x] == 0 {
			break
		}
	}

	return Player{v.Vec{float64(x)*cw + cw/2, float64(y)*ch + ch/2, 0}, v.Vec{1, 0, 0}, 0.2, 0.0025, 500, 250}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return sw, sh }
func (g *game) Update() error {
	dt := float64(time.Now().UnixMilli() - g.prevFrameTime)
	g.prevFrameTime = time.Now().UnixMilli()

	// dp - delta position(new pos - cur pos)
	upp := func(dp v.Vec) { // Update player's position
		np := v.Add(g.p.Pos, dp)                                       // new pos
		ni := Pair[int, int]{int(np.X / cw), int(np.Y / ch)}           // new position index
		ci := Pair[int, int]{int(g.p.Pos.X / cw), int(g.p.Pos.Y / ch)} // current position index
		if g.maze[ni.Y][ni.X] == 0 {
			if ni.X != ci.X && ni.Y != ci.Y {
				if g.maze[ci.Y+(ni.Y-ci.Y)][ci.X] == 0 || g.maze[ci.Y][ci.X+(ni.X-ci.X)] == 0 {
					g.p.Pos = np
				}
			} else {
				g.p.Pos = np
			}
		}
	}

	// Movement:
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		upp(v.Mul(g.p.Dir, g.p.ms*dt))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		upp(v.Mul(g.p.Dir, -g.p.ms*dt))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		upp(v.Mul(*v.RotateZ(&g.p.Dir, math.Pi/2), g.p.ms*dt))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		upp(v.Mul(*v.RotateZ(&g.p.Dir, -math.Pi/2), g.p.ms*dt))
	}

	// Rotation
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.p.Dir = *v.RotateZ(&g.p.Dir, g.p.rs*dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.p.Dir = *v.RotateZ(&g.p.Dir, -g.p.rs*dt)
	}
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	drawLine := func(a, b v.Vec, clr color.Color) {
		ebitenutil.DrawLine(screen, a.X, sh-a.Y, b.X, sh-b.Y, clr)
	}

	// Maze
	for y := range g.maze {
		for x := range g.maze[y] {
			if g.maze[y][x] == 1 {
				ebitenutil.DrawRect(screen, cw*float64(x), sh-float64(mh-y-1)*ch, cw, -ch, color.RGBA{255, 192, 203, 255})
			}
			if g.maze[y][x] == 2 {
				ebitenutil.DrawRect(screen, cw*float64(x), sh-float64(mh-y-1)*ch, cw, -ch, color.RGBA{255, 255, 203, 255})
			}
		}
	}

	// Player
	ebitenutil.DrawCircle(screen, g.p.Pos.X, sh-g.p.Pos.Y, 1, color.White)
	drawLine(g.p.Pos, v.Add(g.p.Pos, v.Mul(g.p.Dir, 10)), color.White)

	// Rays
	// h := v.Mul(g.p.Dir, float64(g.p.h))
	// r := v.Mul(*v.RotateZ(&g.p.Dir, math.Pi/2), float64(g.p.a)/2)
	// pa, pb := v.Sub(h, r), v.Add(h, r)
	// ab := v.Sub(pb, pa)
	// for i := 0; i < g.rc; i++ {
	// 	ak := v.Mul(v.Div(ab, float64(g.rc)), float64(i))
	// 	pk := v.Add(pa, ak)
	// 	n := v.Normalize(pk)
	// 	drawLine(g.p.Pos, v.Add(g.p.Pos, v.Mul(n, GetLengthToIntersection(g, n))), color.RGBA{255, 255, 0, 255})
	// }
	drawLine(g.p.Pos, v.Add(g.p.Pos, v.Mul(g.p.Dir, GetLengthToIntersection(g, g.p.Dir))), color.RGBA{255, 255, 0, 255})

	screen.DrawImage(g.sb, &ebiten.DrawImageOptions{})
}

// Returns length of the ray casted from the player position along the directional normal vector(v) after hitting the wall
func GetLengthToIntersection(g *game, d v.Vec) float64 {
	k := d.Y / d.X
	// Calculating starting distance to horizontal(lx) and vertical neighbors(ly)
	p := v.Vec{g.p.Pos.X / cw, g.p.Pos.Y / ch, 0} // player position in cells
	frac := func(a float64) float64 { return a - float64(int(a)) }
	fp := v.Vec{frac(p.X), frac(p.Y), 0} // fractional parts
	lp := v.Vec{1 - fp.X, 1 - fp.Y, 0}
	// X:
	fx := k * lp.X
	lx := math.Sqrt(lp.X*lp.X + fx*fx)
	// Y:
	fy := lp.Y / k
	ly := math.Sqrt(lp.Y*lp.Y + fy*fy)

	c := Pair[int, int]{int(p.X), int(p.Y)} // Player's cell
	l := float64(0)                         // Ray's line segment's length
	if d.X == 0 && d.Y == 0 {
		return 0
	}
	if d.X != 0 && (ly == 0 || lx < ly) {
		l += v.Mod(v.Vec{lp.X * cw, fx * ch, 0})
		if d.X > 0 { // Function grows along +X
			// Checkinig right neighbor
			if g.maze[c.X+1][c.Y] != 0 {
				return l
			}
		} else {
			// Checking left neighbor
			if g.maze[c.X-1][c.Y] != 0 {
				return l
			}
		}
	} else {
		l += v.Mod(v.Vec{fy * cw, lp.Y * ch, 0})
		if d.Y > 0 {
			// Checking upper neighbor
			if g.maze[c.X][c.Y+1] != 0 {
				return l
			}
		} else {
			// Checking lower neighbor
			if g.maze[c.X][c.Y-1] != 0 {
				return l
			}
		}
	}
	return 0
}

func main() {
	ebiten.SetWindowSize(sw, sh)
	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
