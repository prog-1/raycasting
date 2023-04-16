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
	ms     = 500        // map width/height(width=height)(in pixels)
	cc     = 4          // cell count in row/column(row=column)
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
	l, h   int     // base length and height of player's FOV sector
}

type game struct {
	Maze *[cc][cc]int
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
	maze := [cc][cc]int{
		{1, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 9, 0, 1},
		{2, 1, 1, 1},
	}

	FlipMazeVertically := func(m *[cc][cc]int) {
		swap := func(a, b *int) {
			tmp := *a
			*a = *b
			*b = tmp
		}

		for c := 0; c < cc/2; c++ {
			for r := 0; r < cc/2; r++ {
				swap(&m[r][c], &m[cc-r-1][c])
			}
		}
	}
	FlipMazeVertically(&maze)

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
			l: 1, h: 1}
	}
	return &game{
		Maze: &maze,
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

	// Rotation
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
		x, y := a.X*cs, (sh - a.Y*cs)
		x2, y2 := b.X*cs, (sh - b.Y*cs)
		ebitenutil.DrawLine(g.sb, x, y, x2, y2, clr)
	}

	// Maze
	for y := range g.Maze {
		for x := range g.Maze[y] {
			if g.Maze[y][x] == 1 { // Pink
				ebitenutil.DrawRect(g.sb, float64(x)*cs, sh-float64(y)*cs, cs, -cs, color.RGBA{255, 192, 203, 255})
			}
			if g.Maze[y][x] == 2 { // White
				ebitenutil.DrawRect(g.sb, float64(x)*cs, sh-float64(y)*cs, cs, -cs, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	// Player
	ebitenutil.DrawCircle(g.sb, g.P.Pos.X*cs, sh-g.P.Pos.Y*cs, 2, color.White)
	drawLine(g.P.Pos, g.P.Pos.Add(g.P.Dir), color.White)

	// Screen manipulations
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sw/2-cc/2*cs, -sh/2+cc/2*cs)
	screen.DrawImage(g.sb, &opts)
	g.sb.Clear()
}

func main() {
	ebiten.SetWindowSize(sw, sh)

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
