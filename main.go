package main

import (
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
	cs     = 50         // cell size - width/height(width=height)(in pixels)
)

func RotateZ(v *vector2.Vector2, rad float64) *vector2.Vector2 {
	x := v.X*math.Cos(rad) - v.Y*math.Sin(rad)
	y := v.X*math.Sin(rad) + v.Y*math.Cos(rad)
	return &vector2.Vector2{X: x, Y: y}
}

type Pair[T, U any] struct {
	X T
	Y U
}

type Player struct {
	Pos    *vector2.Vector2
	Dir    *vector2.Vector2
	FOV    float64
	ms, rs float64 // Movement and Rotation speed
}

type game struct {
	Maze [][]int
	P    Player
	RC   int           //Ray Count
	pft  int64         // Previous frame time
	sb   *ebiten.Image // Screen buffer
}

func NewGame() *game {
	// Values in the maze:
	// 0 - empty space
	// 1 - pink wall
	// 9 - player spawn
	// NOTE: Maze must have player spawn walls along the perimeter
	g := &game{
		Maze: [][]int{
			{1, 1, 1, 1, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 0, 0, 0, 1},
			{1, 9, 0, 0, 1},
			{2, 1, 1, 1, 1},
		},
		P: Player{
			Dir: &vector2.Vector2{X: 1, Y: 0},
			ms:  0.2, rs: 0.00125,
			FOV: math.Pi / 3,
		},
		RC:  101,
		pft: 0,
		sb:  ebiten.NewImage(sw, sh),
	}
	FlipVertically := func(m *[][]int) {
		swap := func(a, b *int) {
			tmp := *a
			*a = *b
			*b = tmp
		}

		for c := 0; c < len(*m)/2; c++ {
			for r := 0; r < len((*m)[0])/2; r++ {
				swap(&(*m)[r][c], &(*m)[len(*m)-r-1][c])
			}
		}
	}
	FlipVertically(&g.Maze)

	// Player Pos:
	var x, y int
	for y = range g.Maze {
		for x = range g.Maze[y] {
			if g.Maze[y][x] == 9 {
				break
			}
		}
		if g.Maze[y][x] == 9 {
			break
		}
	}
	if x == len(g.Maze[0])-1 && y == len(g.Maze)-1 { // Maze's bottom right cell(which has to be a wall) -> No 9 in maze
		panic("Maze must have player spawn(9)")
	}
	g.P.Pos = &vector2.Vector2{X: float64(x) + 0.5, Y: float64(y) + 0.5}

	return g
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

	// Rays:
	for rad := g.P.FOV / 2; rad > -g.P.FOV/2; rad -= g.P.FOV / float64(g.RC) {
		d := RotateZ(g.P.Dir, rad)
		l := GetRayLengthToIntersection(g.P.Pos, d, &g.Maze)
		drawLine(g.P.Pos, g.P.Pos.Add(d.MulScalar(l)), color.RGBA{255, 255, 0, 255})
	}

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
	opts.GeoM.Scale(-1, -1)                                                                 // Flipping
	opts.GeoM.Translate(float64(screen.Bounds().Max.X/2), float64(screen.Bounds().Max.Y/2)) // Moving coord. center to screen center
	screen.DrawImage(g.sb, &opts)
	g.sb.Clear()
}

// Returns length of the ray casted from the point p in the direction of d to the wall in the map m
func GetRayLengthToIntersection(p, d *vector2.Vector2, m *[][]int) float64 {
	d = d.Normalize()
	frac := func(a float64) float64 {
		return a - float64(int(a))
	}

	var l vector2.Vector2                                                                                 // Total line length from stepping to the neighbors along OX and OY
	step := vector2.Vector2{X: math.Sqrt(1 + (d.Y/d.X)*(d.Y/d.X)), Y: math.Sqrt(1 + (d.X/d.Y)*(d.X/d.Y))} // Ray's length from the current intersection point to another along OX and OY
	var ms Pair[int, int]                                                                                 // maze step(+1 or -1)
	if d.X > 0 {
		ms.X = 1
		l.X = (1 - frac(p.X)) * step.X
	} else {
		ms.X = -1
		l.X = frac(p.X) * step.X
	}
	if d.Y > 0 {
		ms.Y = 1
		l.Y = (1 - frac(p.Y)) * step.Y
	} else {
		ms.Y = -1
		l.Y = frac(p.Y) * step.Y
	}
	var len float64                         // Ray's length to the intersection point with the wall
	c := Pair[int, int]{int(p.X), int(p.Y)} // current cell
	for (*m)[int(c.Y)][int(c.X)] == 0 /*empty*/ || (*m)[int(c.Y)][int(c.X)] == 9 /*player's spawning point*/ {
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
