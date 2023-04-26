package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//-----------------------Copyright-Notice----------------------------
/*
	Copyright (c) 2023 Vladimir Stukalov
	All right reserved.
*/
//-------------------------Declaration------------------------------

const (
	sW = 1280 //screen width [in pixels]
	sH = 720  //screen height [in pixels]
)

type Game struct {
	screenBuffer  *ebiten.Image
	width, height int //screen width and height [in pixels]
	//variables
	gameMap   [][]int   //game map 2d matrix
	playerPos vector    //player position [in cells]
	viewDir   vector    //player view direction
	cellSize  float64   //[in pixels] (not in 'draw' due use in 'update')
	pt        time.Time //previous frame time (for movement)
}

//line with start & end points
type line struct {
	a, b vector
}

type vector struct {
	x, y float64
}

//---------------------------Update-------------------------------------

func (g *Game) Update() error {

	//Time
	ft := time.Now().Sub(g.pt).Seconds() //difference from cur frame and last saved
	g.pt = time.Now()                    //updating previous frame time
	mspd := ft * 80                      //movement speed
	rspd := ft * 100                     //rotation speed

	move := func(nextPos vector) {
		//collision handling
		if g.gameMap[int(g.playerPos.y)][int(nextPos.x)] == 0 { //if there is air when we move more along x
			g.playerPos.x = nextPos.x // move along x
		}
		if g.gameMap[int(nextPos.y)][int(g.playerPos.x)] == 0 { //if there is air when we move more along y
			g.playerPos.y = nextPos.y // move along y
		}
	}

	//Player WASD Movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		nextPos := add(g.playerPos, scale(scale(g.viewDir, 1/g.cellSize), mspd)) // nextPos = playerPos + viewdir*(1/25)
		move(nextPos)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		//same as W, bust substracting viewdir from player pos
		nextPos := subtract(g.playerPos, scale(scale(g.viewDir, 1/g.cellSize), mspd)) // nextPos = playerPos - viewdir*(1/25)
		move(nextPos)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		//rotating on 90 before adding
		nextPos := add(g.playerPos, rotate(scale(scale(g.viewDir, 1/g.cellSize), mspd), -math.Pi/2)) // nextPos = playerPos + viewdir*(1/25) on +90°
		move(nextPos)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		//rotating on -90 before adding
		nextPos := add(g.playerPos, rotate(scale(scale(g.viewDir, 1/g.cellSize), mspd), math.Pi/2)) // nextPos = playerPos + viewdir*(1/25) on -90°
		move(nextPos)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.viewDir = rotate(g.viewDir, -math.Pi/200*rspd)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.viewDir = rotate(g.viewDir, math.Pi/200*rspd)
	}

	return nil
}

//---------------------------Draw-------------------------------------

func (g *Game) Draw(screen *ebiten.Image) {

	//---------------------
	//-----Declaration-----
	//---------------------

	mW := 600.0 //map width [in pixels]
	mH := 600.0 //map height [in pixels]

	viewLen := 150.0 //lenght of player view and camera plane lines [in pixels]
	//if len of viewDir = len of camPlane, then FoV = 90 degrees
	mapPos := vector{(sW / 2) - (mW / 2), (sH / 2) - (mH / 2)} // map position on the screen (top left corner) [world coordinates]
	segNumber := 100.0                                         // number of segments of camera plane line for rays

	playerPixelPos := vector{(g.playerPos.x * g.cellSize) + mapPos.x, (g.playerPos.y * g.cellSize) + mapPos.y} // player position [in pixels]

	//---------------------
	//----Screen-Buffer----
	//---------------------

	if g.screenBuffer == nil {
		g.screenBuffer = ebiten.NewImage(g.width, g.height) // creating buffer if we don't have one
	}
	g.screenBuffer.Clear() //clearing previous screen to draw everything one more time

	//---------------------
	//-----Map-Drawing-----
	//---------------------

	ebitenutil.DrawRect(g.screenBuffer, mapPos.x, mapPos.y, mW, mH, color.RGBA{50, 50, 50, 255} /*Grey*/) //map background

	//draw cell with proper color
	drawCell := func(ci, cj float64, clr color.RGBA) { //ci & cj - cell index in gameMap [player coordinates]
		//map position [world coordinates]  +  cell position [world coordinates]
		cX := mapPos.x + (cj * g.cellSize)
		cY := mapPos.y + (ci * g.cellSize)
		ebitenutil.DrawRect(g.screenBuffer, cX, cY, g.cellSize, g.cellSize, clr)
	}

	//draw map cells
	for i := range g.gameMap { //for each column

		for j := range g.gameMap[i] { // for each row

			//If current cell is wall
			if g.gameMap[i][j] == 1 {
				drawCell(float64(i), float64(j), color.RGBA{200, 200, 200, 255} /*Light grey*/)

			} else if g.gameMap[i][j] == 2 {
				drawCell(float64(i), float64(j), color.RGBA{255, 0, 0, 255} /*Red*/)
			}

		}
	}

	//---------------------
	//----Field-of-View----
	//---------------------

	camA := add(scale(g.viewDir, viewLen), rotate(scale(g.viewDir, viewLen), math.Pi/2))  //first camera point
	camB := add(scale(g.viewDir, viewLen), rotate(scale(g.viewDir, viewLen), -math.Pi/2)) //last camera point
	camPlane := line{camA, camB}                                                          //camera plane line

	//---------------------
	//-----Ray-Drawing-----
	//---------------------

	//segment length for each projection
	segLenX := (camPlane.b.x - camPlane.a.x) / segNumber
	segLenY := (camPlane.b.y - camPlane.a.y) / segNumber

	for i := 0.0; i <= segNumber; i++ { //for each segment

		var r vector //ray point on camPlane
		r.x = (camPlane.a.x + (segLenX * i))
		r.y = (camPlane.a.y + (segLenY * i))
		r = norm(r) //ray direction unit vector

		res := scale(r, g.DDA(r)*g.cellSize) //calculate ray length

		ebitenutil.DrawLine(g.screenBuffer, playerPixelPos.x, playerPixelPos.y, playerPixelPos.x+res.x, playerPixelPos.y+res.y, color.RGBA{242, 207, 85, 200} /*Yellow*/) //draw ray
	}

	//---------------------
	//--------Debug--------
	//---------------------

	//player view line drawing
	ebitenutil.DrawLine(g.screenBuffer, playerPixelPos.x, playerPixelPos.y, playerPixelPos.x+scale(g.viewDir, viewLen).x, playerPixelPos.y+scale(g.viewDir, viewLen).y, color.RGBA{255, 146, 28, 200} /*Orange*/)

	//camera plane line drawing
	ebitenutil.DrawLine(g.screenBuffer, playerPixelPos.x+camPlane.a.x, playerPixelPos.y+camPlane.a.y, playerPixelPos.x+camPlane.b.x, playerPixelPos.y+camPlane.b.y, color.RGBA{132, 132, 255, 200} /*Blue*/)
	//adding player position to convert from player coordinates to world coordinates

	//Draw player
	ebitenutil.DrawCircle(g.screenBuffer, playerPixelPos.x, playerPixelPos.y, 10, color.RGBA{100, 180, 255, 255} /*Light blue*/)

	//---------------------
	//----Screen-Drawing---
	//---------------------

	// var opts ebiten.DrawImageOptions                    //declaring screen operations
	// opts.GeoM.Translate(-g.playerPos.x, -g.playerPos.y) // converting screen world coordinates to player's coordinates

	// //rotation
	// var m ebiten.GeoM //declaring matrix
	// /*
	// 	00 10
	// 	a  b
	// 	01 11
	// 	c  d
	// */
	// //setting matrix
	// //90° rotation = [x,y] -> [-y,x]
	// m.SetElement(0, 0, -g.viewDir.y) //a (-| viewdir.x)
	// m.SetElement(0, 1, g.viewDir.x)  //b (-| viewdir.y)
	// m.SetElement(1, 0, g.viewDir.x)  //c (viewdir.x)
	// m.SetElement(1, 1, g.viewDir.y)  //d (viewdir.y)
	// m.Invert()                       //taking inverse matrix
	// opts.GeoM.Concat(m)              //multiplying "opts matrix" with "our matrix"
	// opts.GeoM.Scale(1, -1)           // scaling matrix for proper player movement & rotation

	// opts.GeoM.Translate(float64(screen.Bounds().Max.X)/2, float64(screen.Bounds().Max.Y)/2) //centering the screen
	// screen.DrawImage(g.screenBuffer, &opts)                                                 //drawing screen buffer

	screen.DrawImage(g.screenBuffer, nil) //drawing static screen buffer

}

//-------------------------Functions----------------------------------

// DDA - calculates length of the ray [in cells]
// inputs unit (direction) vector
func (g Game) DDA(v vector) (rayLen float64) {

	var curCell vector                    //current cell [in cells]
	curCell.x = math.Trunc(g.playerPos.x) //pmp.x - frac.x
	curCell.y = math.Trunc(g.playerPos.y) //pmp.y - frac.y

	var step vector            //distance to row and column of ray in cell [in cells]
	step.x = math.Abs(1 / v.x) // √ 1^2 + k^2
	step.y = math.Abs(1 / v.y) // √ 1^2 + (1/k)^2

	var dist vector //initial distance to first row and column [in cells]

	var mapd vector //step where to go on each direction [in cells]
	if v.x < 0 {
		mapd.x = -1
		dist.x = (g.playerPos.x - curCell.x) * step.x
	} else {
		mapd.x = 1
		dist.x = (curCell.x + 1 - g.playerPos.x) * step.x //right neighbor cell - player x
	}
	if v.y < 0 {
		mapd.y = -1
		dist.y = (g.playerPos.y - curCell.y) * step.y
	} else {
		mapd.y = 1
		dist.y = (curCell.y + 1 - g.playerPos.y) * step.y //bottom neighbor cell - player y
	}

	for {
		if dist.x < dist.y {
			rayLen = dist.x
			dist.x += step.x
			curCell.x += mapd.x
		} else /* dist.x >= dist.y */ {
			rayLen = dist.y
			dist.y += step.y
			curCell.y += mapd.y
		}
		if g.gameMap[int(curCell.y)][int(curCell.x)] != 0 {
			return rayLen
		}
	}
}

// normalize - converts vector to unit vector
func norm(v vector) vector {
	mod := mod(v)
	return vector{x: v.x / mod, y: v.y / mod}
}

// returns module (length) of the vector
func mod(v vector) float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y))
}

func rotate(p vector, angle float64) vector {
	sin, cos := math.Sincos(angle)
	return vector{
		x: p.x*cos - p.y*sin,
		y: p.x*sin + p.y*cos}
}

func subtract(a, b vector) vector {
	return vector{x: a.x - b.x, y: a.y - b.y}
}

func add(a, b vector) vector {
	return vector{x: a.x + b.x, y: a.y + b.y}
}

func scale(a vector, v float64) vector { //v - value
	return vector{x: a.x * v, y: a.y * v}
}

func initGameMap() [][]int {
	return [][]int{ //24 x 24 cells
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 1, 1, 1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
}

//---------------------------Main-------------------------------------

func (g *Game) Layout(inWidth, inHeight int) (outWidth, outHeight int) {
	return g.width, g.height
}

func main() {

	//Window
	ebiten.SetWindowSize(sW, sH)
	ebiten.SetWindowTitle("Raycasting")
	ebiten.SetWindowResizable(true) //enablening window resize

	//Game instance
	g := NewGame(sW, sH)                      //creating game instance
	if err := ebiten.RunGame(g); err != nil { //running game
		log.Fatal(err)
	}
}

// New game instance function
func NewGame(width, height int) *Game {

	//playerPos := vector{(sW / 2) - 10, (sH / 2) - 10} //player initial position [in pixels]
	playerPos := vector{12, 12} //player initial position [in cells]
	viewDir := vector{0, -1}    //player view direction unit vector [in cells]
	cellSize := 25.0            // size of each cell [in pixels] (not in 'draw' due use in 'update')
	pt := time.Now()            //previous frame time (for movement)

	return &Game{width: width, height: height, gameMap: initGameMap(), playerPos: playerPos, viewDir: viewDir, cellSize: cellSize, pt: pt}
}
