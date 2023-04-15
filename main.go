package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/deeean/go-vector/vector2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sw, sh = 1920, 1080 //(in pixels)
	mw, mh = 4, 4       // map width/height(in cells)
)

type Pair[T, U any] struct {
	X T
	Y U
}

type Player struct {
	Pos    vector2.Vector2
	Dir    vector2.Vector2
	ms, rs float64 // Movement and Rotation speed
	l, h   int     // base length and height of player's FOV sector
}

type game struct {
	Maze *[mh][mw]int
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
	maze := [mh][mw]int{
		{1, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 9, 0, 1},
		{2, 1, 1, 1},
	}

	FlipMazeVertically := func(m *[mh][mw]int) {
		swap := func(a, b *int) {
			tmp := *a
			*a = *b
			*b = tmp
		}

		for c := 0; c < mw/2; c++ {
			for r := 0; r < mh/2; r++ {
				swap(&m[r][c], &m[mh-r-1][c])
			}
		}
	}
	fmt.Println(maze)
	FlipMazeVertically(&maze)

	NewPlayer := func() Player {
		var x, y int
		for y = range maze {
			for x = range maze[y] {
				if maze[y][x] == 9 {
					break
				}
			}
		}
		if x == mw && y == mh { // Maze's bottom right cell(which has to be a wall) -> No 9 in maze
			panic("Maze must have player spawn(9)")
		}
		return Player{
			Pos: vector2.Vector2{float64(x), float64(y)},
			Dir: vector2.Vector2{1, 0},
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
	// dt := float64(time.Now().UnixMilli() - g.pft)
	// g.pft = time.Now().UnixMilli()

	// // dp - delta position(new pos - cur pos)
	// upp := func(dp vector2.Vector2) { // Update player's position
	// 	np := vector2.Vector2{g.P.Pos.X + dp.X, g.P.Pos.Y + dp.Y} // New position
	// 	ci := Pair[int, int]{int(g.P.Pos.X), int(g.P.Pos.Y)}      // Index of the cell containing current position
	// 	ni := Pair[int, int]{int(np.X), int(np.Y)}                // Index of the cell containing new position
	// 	// Checking whether we're not moving into a wall:
	// 	if ni.X != ci.X || ni.Y != ci.Y { // We're not in the previous cell
	// 		if g.Maze[ni.Y][ni.X] == 0 {
	// 			if g.Maze[ci.Y][ci.X+(ni.X-ci.X)] == 0 || g.Maze[ci.Y+(ci.Y-ni.Y)][ci.X] == 0 {

	// 			}
	// 		}
	// 	} else {
	// 		g.P.Pos = np
	// 	}
	// }

	// // Movement:
	// if ebiten.IsKeyPressed(ebiten.KeyW) {
	// 	upp(v.Mul(g.P.Dir, g.P.ms*dt))
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyS) {
	// 	upp(v.Mul(g.P.Dir, -g.P.ms*dt))
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyA) {
	// 	upp(v.Mul(*v.RotateZ(&g.P.Dir, math.Pi/2), g.P.ms*dt))
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyD) {
	// 	upp(v.Mul(*v.RotateZ(&g.P.Dir, -math.Pi/2), g.P.ms*dt))
	// }

	// // Rotation
	// if ebiten.IsKeyPressed(ebiten.KeyQ) {
	// 	g.P.Dir = *v.RotateZ(&g.P.Dir, g.P.rs*dt)
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyE) {
	// 	g.P.Dir = *v.RotateZ(&g.P.Dir, -g.P.rs*dt)
	// }
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	const cw, ch = float64(sw) / mw, float64(sh) / mh // cell width/height(in pixels)

	// drawLine := func(a, b v.Vec, clr color.Color) {
	// 	ebitenutil.DrawLine(screen, a.X, sh-a.Y, b.X, sh-b.Y, clr)
	// }

	// Maze
	for y := range g.Maze {
		for x := range g.Maze[y] {
			if g.Maze[y][x] == 1 { // Pink
				ebitenutil.DrawRect(screen, cw*float64(x), sh-float64(y)*ch, cw, -ch, color.RGBA{255, 192, 203, 255})
			}
			if g.Maze[y][x] == 2 { // White
				ebitenutil.DrawRect(screen, cw*float64(x), sh-float64(y)*ch, cw, -ch, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	// Player
	ebitenutil.DrawCircle(screen, g.P.Pos.X, sh-g.P.Pos.Y, 1, color.White)
	// drawLine(g.P.Pos, v.Add(g.P.Pos, v.Mul(g.P.Dir, 10)), color.White)

	// // Rays
	// // h := v.Mul(g.P.Dir, float64(g.P.h))
	// // r := v.Mul(*v.RotateZ(&g.P.Dir, math.Pi/2), float64(g.P.a)/2)
	// // pa, pb := v.Sub(h, r), v.Add(h, r)
	// // ab := v.Sub(pb, pa)
	// // for i := 0; i < g.rc; i++ {
	// // 	ak := v.Mul(v.Div(ab, float64(g.rc)), float64(i))
	// // 	pk := v.Add(pa, ak)
	// // 	n := v.Normalize(pk)
	// // 	drawLine(g.P.Pos, v.Add(g.P.Pos, v.Mul(n, GetLengthToIntersection(g, n))), color.RGBA{255, 255, 0, 255})
	// // }
	// drawLine(g.P.Pos, v.Add(g.P.Pos, v.Mul(g.P.Dir, GetLengthToIntersection(g, g.P.Dir))), color.RGBA{255, 255, 0, 255})

	screen.DrawImage(g.sb, &ebiten.DrawImageOptions{})
}

// Returns length of the ray casted from the player position along the directional normal vector(v) after hitting the wall
// func GetLengthToIntersection(g *game, d v.Vec) float64 {
// 	k := d.Y / d.X
// 	// Calculating starting distance to horizontal(lx) and vertical neighbors(ly)
// 	p := v.Vec{g.P.Pos.X / cw, g.P.Pos.Y / ch, 0} // player position in cells
// 	// var x float64
// 	// if d.X < 0 {
// 	// 	x = p.X - math.Trunc(p.X)
// 	// } else {
// 	// 	x = math.Trunc(p.X+1) - p.X
// 	// }
// 	// var y float64
// 	// if d.Y < 0 {
// 	// 	y = p.Y - math.Trunc(p.Y)
// 	// } else {
// 	// 	y = math.Trunc(p.Y+1) - p.Y
// 	// }
// 	fp := v.Vec{p.X - math.Trunc(p.X), p.Y - math.Trunc(p.Y), 0} // fractional parts

// 	lp := v.Vec{1 - fp.X, 1 - fp.Y, 0}
// 	// X:
// 	fx := k * lp.X
// 	lx := math.Sqrt(lp.X*lp.X + fx*fx)
// 	// Y:
// 	fy := lp.Y / k
// 	ly := math.Sqrt(lp.Y*lp.Y + fy*fy)

// 	c := Pair[int, int]{int(p.X), int(p.Y)} // Player's cell
// 	l := float64(0)                         // Ray's line segment's length
// 	if d.X == 0 && d.Y == 0 {
// 		return 0
// 	}
// 	if d.X != 0 && (ly == 0 || lx < ly) {
// 		l += v.Mod(v.Vec{lp.X * cw, fx * ch, 0})
// 		if d.X > 0 { // Function grows along +X
// 			// Checkinig right neighbor
// 			if g.Maze[c.X+1][c.Y] != 0 {
// 				return l
// 			}
// 		} else {
// 			// Checking left neighbor
// 			if g.Maze[c.X-1][c.Y] != 0 {
// 				return l
// 			}
// 		}
// 	} else {
// 		l += v.Mod(v.Vec{fy * cw, lp.Y * ch, 0})
// 		if d.Y > 0 {
// 			// Checking upper neighbor
// 			if g.Maze[c.X][c.Y+1] != 0 {
// 				return l
// 			}
// 		} else {
// 			// Checking lower neighbor
// 			if g.Maze[c.X][c.Y-1] != 0 {
// 				return l
// 			}
// 		}
// 	}
// 	return 0
// }

func main() {
	ebiten.SetWindowSize(sw, sh)

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
