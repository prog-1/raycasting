package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/deeean/go-vector/vector2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sw, sh = 1920, 1080 // screen width and height(in pixels)
	ms     = 500        // map width/height(width=height)(in pixels)
	cc     = 7          // cell count in row/column(row=column)
)

type Pair[T, U any] struct {
	X T
	Y U
}

func RotateZ(v *vector2.Vector2, rad float64) *vector2.Vector2 {
	x := v.X*math.Cos(rad) - v.Y*math.Sin(rad)
	y := v.X*math.Sin(rad) + v.Y*math.Cos(rad)
	return &vector2.Vector2{X: x, Y: y}
}

type Player struct {
	Pos    *vector2.Vector2
	Dir    *vector2.Vector2
	ms, rs float64 // Movement and Rotation speed
	l, w   int     // base length and width of player's FOV sector
}

type game struct {
	Maze *[cc][cc]int
	P    Player
	RC   int           //Ray Count
	pft  int64         // Previous frame time
	sb   *ebiten.Image // Screen buffer
}

func NewGame() *game {
	var maze *[cc][cc]int
	// Values in the maze:
	// 0 - empty space
	// 1 - pink wall
	// 9 - player spawn
	// NOTE: Maze must have player spawn walls along the perimeter
	maze = &[cc][cc]int{
		{1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1},
		{1, 9, 0, 0, 0, 0, 1},
		{2, 1, 1, 1, 1, 1, 1},
	}
	FlipMazeVertically := func(m *[cc][cc]int) *[cc][cc]int {
		m2 := *m
		swap := func(a, b *int) {
			tmp := *a
			*a = *b
			*b = tmp
		}

		for c := 0; c < cc/2; c++ {
			for r := 0; r < cc/2; r++ {
				swap(&m2[r][c], &m2[cc-r-1][c])
			}
		}
		return &m2
	}
	maze = FlipMazeVertically(maze)

	NewPlayer := func() Player {
		var x, y int
		for y = range maze {
			for x = range maze[y] {
				if maze[y][x] == 9 {
					break
				}
			}
			if maze[y][x] == 9 {
				break
			}
		}
		if x == cc && y == cc { // Maze's bottom right cell(which has to be a wall) -> No 9 in maze
			panic("Maze must have player spawn(9)")
		}
		return Player{
			Pos: &vector2.Vector2{X: float64(x) + 0.5, Y: float64(y) + 0.5},
			Dir: &vector2.Vector2{X: 1, Y: 0},
			ms:  0.2, rs: 0.0025,
			l: 1, w: 1}
	}
	return &game{
		Maze: maze,
		P:    NewPlayer(),
		RC:   101,
		pft:  0,
		sb:   ebiten.NewImage(sw, sh),
	}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return sw, sh }
func (g *game) Update() error {
	dt := float64(time.Now().UnixMilli() - g.pft)
	sp := dt / 100 // speed factor
	g.pft = time.Now().UnixMilli()

	// dp - delta position(new pos - cur pos)
	upp := func(dp vector2.Vector2) { // Update player's position
		np := vector2.Vector2{X: g.P.Pos.X + dp.X, Y: g.P.Pos.Y + dp.Y} // New position
		ci := Pair[int, int]{int(g.P.Pos.X), int(g.P.Pos.Y)}            // Index of the cell containing current position
		ni := Pair[int, int]{int(np.X), int(np.Y)}                      // Index of the cell containing new position
		// Checking whether we're not moving into a wall:
		if ni.X != ci.X || ni.Y != ci.Y { // We're not in the previous cell
			if g.Maze[ni.Y][ni.X] == 0 || g.Maze[ni.Y][ni.X] == 9 { // We're moving into the empty cell
				// Checking diagonal neighbors:
				if g.Maze[ci.Y][ci.X+(ni.X-ci.X)] == 0 || g.Maze[ci.Y][ci.X+(ni.X-ci.X)] == 9 {
					if g.Maze[ci.Y+(ni.Y-ci.Y)][ci.X] == 0 || g.Maze[ci.Y+(ni.Y-ci.Y)][ci.X] == 9 {
						g.P.Pos = &np
					}
				}
			}
		} else {
			g.P.Pos = &np
		}
	}

	// Movement:
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		upp(*g.P.Dir.MulScalar(g.P.ms * sp))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		upp(*g.P.Dir.MulScalar(-g.P.ms * sp))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		upp(*RotateZ(g.P.Dir, math.Pi/2).MulScalar(g.P.ms * sp))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		upp(*RotateZ(g.P.Dir, -math.Pi/2).MulScalar(g.P.ms * sp))
	}

	// Rotation:
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.P.Dir = RotateZ(g.P.Dir, g.P.rs*dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.P.Dir = RotateZ(g.P.Dir, -g.P.rs*dt)
	}
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	const cs = float64(ms) / cc // cell size(width/height(in pixels))

	drawLine := func(a, b *vector2.Vector2, clr color.Color) {
		u := a.MulScalar(cs)
		v := b.MulScalar(cs)
		ebitenutil.DrawLine(g.sb, u.X, u.Y, v.X, v.Y, clr)
	}

	// Maze:
	for y := range g.Maze {
		for x := range g.Maze[y] {
			if g.Maze[y][x] != 0 && g.Maze[y][x] != 9 {
				cp := vector2.Vector2{X: float64(x) * cs, Y: float64(y) * cs} // cell pos(pixels)
				var clr color.Color
				switch g.Maze[y][x] {
				case 1:
					clr = color.RGBA{255, 192, 203, 255} // Pink
				case 2:
					clr = color.White
				}

				ebitenutil.DrawRect(g.sb, cp.X, cp.Y, cs, cs, clr)
			}
		}
	}

	// Player:
	ebitenutil.DrawCircle(g.sb, g.P.Pos.X*cs, g.P.Pos.Y*cs, 10, color.White)
	drawLine(g.P.Pos, g.P.Pos.Add(g.P.Dir), color.White)

	// Rays:
	g.P.Pos = &vector2.Vector2{X: 1, Y: 1}
	g.P.Dir = &vector2.Vector2{X: 0.6, Y: 0.8}
	l := GetRayLengthToIntersection(*g.P.Pos, *g.P.Dir, *g.Maze)
	fmt.Println(l)
	drawLine(g.P.Pos, g.P.Pos.Add(g.P.Dir.MulScalar(l)), color.RGBA{0, 0, 255, 255})

	// Screen manipulations:
	opts := ebiten.DrawImageOptions{}
	p := ebiten.GeoM{}
	opts.GeoM.Translate(-g.P.Pos.X*cs, -g.P.Pos.Y*cs) // Translation: World -> Player local
	// x, y rotated by 90 deg. = -y,x
	p.SetElement(0, 0, -g.P.Dir.Y) // p1.X
	p.SetElement(0, 1, g.P.Dir.X)  // p1.Y
	p.SetElement(1, 0, g.P.Dir.X)  // p2.X
	p.SetElement(1, 1, g.P.Dir.Y)  // p2.Y
	p.Invert()
	opts.GeoM.Concat(p)                                                                     // Multiplication
	opts.GeoM.Scale(-1, -1)                                                                 // Flipping along OX
	opts.GeoM.Translate(float64(screen.Bounds().Max.X/2), float64(screen.Bounds().Max.Y/2)) // Moving coord. center to screen center
	screen.DrawImage(g.sb, &opts)
	g.sb.Clear()
}

// Returns length of the ray casted from the point p in the direction of d to the wall in the map m
func GetRayLengthToIntersection(p, d vector2.Vector2, m [cc][cc]int) float64 {
	d = *d.Normalize()
	mod := func(a float64) float64 {
		if a > 0 {
			return a
		} else {
			return -a
		}
	}
	frac := func(a float64) float64 {
		return a - float64(int(a))
	}

	var l vector2.Vector2     // Total line length from stepping to the neighbors along OX and OY
	var k1 vector2.Vector2    // Distance to the nearest neighbor along OX and OY, which we can express using starting point coordinates(p.X and P.Y)
	tg := mod(d.Y) / mod(d.X) // tangent of the angle between lx or ly and frac(k.X) or frac(k.Y)
	var ms Pair[int, int]     // maze step(+1 or -1)
	if d.X > 0 {
		ms.X = 1
		k1.X = 1 - frac(p.X)
	} else {
		ms.X = -1
		k1.X = frac(p.X)
	}
	if d.Y > 0 {
		ms.Y = 1
		k1.Y = 1 - frac(p.Y)
	} else {
		ms.Y = -1
		k1.Y = frac(p.Y)
	}
	k2 := vector2.Vector2{X: k1.X * tg, Y: k1.Y / tg}                                 // tg = y/x => y = xtg, x = y/tg
	l.X, l.Y = math.Sqrt(k1.X*k1.X+k2.X*k2.X), math.Sqrt(k1.Y*k1.Y+k2.Y*k2.Y)         // Ray's length from p to intersection point with the nearest neighbor along OX and OY
	step := vector2.Vector2{X: math.Sqrt(1 + k2.X*k2.X), Y: math.Sqrt(1 + k2.Y*k2.Y)} // Ray's length from the current intersection point to another along OX and OY
	var len float64                                                                   // Ray's length to the intersection point with the wall
	c := Pair[int, int]{int(p.X), int(p.Y)}                                           // current cell
	for m[int(c.X)][int(c.Y)] == 0 /*empty*/ || m[int(c.X)][int(c.Y)] == 9 /*player spawning point*/ {
		if l.X < l.Y {
			len = l.X
			l.X += step.X
			c.X += ms.X
		} else {
			len = l.Y
			l.Y += step.Y
			c.Y += ms.Y
		}
	}
	return len
}

func main() {
	ebiten.SetWindowSize(sw, sh)

	g := NewGame()
	g.P.Dir = RotateZ(g.P.Dir, math.Pi)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
