package main

import (
	"image/color"
	"log"
	"math"

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
	sW       = 1280 //screen width
	sH       = 720  //screen height
	mW       = 600  //map width
	mH       = 600  //map height
	cellSize = 25   // size of each cell in pixels
	viewLen  = 150  //lenght of player view and camera plane lines
	//if len of viewDir = len of camPlane, then FoV = 90 degrees
	segNumber = 100 // number of segments of camera plane line for rays
	playerCol = 10  //player collision size
)

type Game struct {
	screenBuffer  *ebiten.Image
	width, height int //screen width and height
	//global variables
	gameMap   [][]int //game map 2d matrix
	mapPos    vector  //map position on the screen (top left corner)
	playerPos vector  //player position
	viewDir   vector  //player view direction
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

	//Player WASD Movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		//collision handling
		if nextPos := add(g.playerPos, multiply(g.viewDir, playerCol)); g.gameMap[int((nextPos.y-g.mapPos.y)/cellSize)][int((nextPos.x-g.mapPos.x)/cellSize)] == 0 { //if next position of the player is not the wall
			//g.playerPos = add(g.playerPos, divide(g.viewDir, viewLen))
			g.playerPos = add(g.playerPos, g.viewDir)
			//adding viewDir to playerPos to move forward on 1 pixel
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if nextPos := subtract(g.playerPos, multiply(g.viewDir, playerCol)); g.gameMap[int((nextPos.y-g.mapPos.y)/cellSize)][int((nextPos.x-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos = subtract(g.playerPos, g.viewDir) //same as W, but subtracting viewDir
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if nextPos := add(g.playerPos, rotate(multiply(g.viewDir, playerCol), -math.Pi/2)); g.gameMap[int((nextPos.y-g.mapPos.y)/cellSize)][int((nextPos.x-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos = add(g.playerPos, rotate(g.viewDir, -math.Pi/2)) //rotating on 90 before adding
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if nextPos := add(g.playerPos, rotate(multiply(g.viewDir, playerCol), math.Pi/2)); g.gameMap[int((nextPos.y-g.mapPos.y)/cellSize)][int((nextPos.x-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos = add(g.playerPos, rotate(g.viewDir, math.Pi/2)) //rotating on -90 before adding
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.viewDir = rotate(g.viewDir, -math.Pi/200)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.viewDir = rotate(g.viewDir, math.Pi/200)
	}

	return nil
}

//---------------------------Draw-------------------------------------

func (g *Game) Draw(screen *ebiten.Image) {

	// # Screen buffer #
	if g.screenBuffer == nil {
		g.screenBuffer = ebiten.NewImage(g.width, g.height) // creating buffer if we don't have one
	}
	g.screenBuffer.Clear() //clearing previous screen to draw everything one more time

	//###
	//draw map background
	ebitenutil.DrawRect(g.screenBuffer, g.mapPos.x, g.mapPos.y, mW, mH, color.RGBA{50, 50, 50, 255} /*Grey*/)

	//draw map cells
	for i := range g.gameMap { //for each column

		for j := range g.gameMap[i] { // for each row

			//If current cell is wall
			if g.gameMap[i][j] == 1 {
				drawCell(g.screenBuffer, g.mapPos, i, j, color.RGBA{200, 200, 200, 255} /*Light grey*/)

			} else if g.gameMap[i][j] == 2 {
				drawCell(g.screenBuffer, g.mapPos, i, j, color.RGBA{255, 0, 0, 255} /*Red*/)
			}

		}
	}

	//Draw player
	ebitenutil.DrawCircle(g.screenBuffer, g.playerPos.x, g.playerPos.y, playerCol, color.RGBA{100, 180, 255, 255} /*Light blue*/)

	//# FIELD OF VIEW #

	camA := add(multiply(g.viewDir, viewLen), rotate(multiply(g.viewDir, viewLen), math.Pi/2))  //first camera point
	camB := add(multiply(g.viewDir, viewLen), rotate(multiply(g.viewDir, viewLen), -math.Pi/2)) //last camera point
	camPlane := line{camA, camB}                                                                //camera plane line

	//# Ray Drawing #
	//segment length for each projection
	segLenX := (camPlane.b.x - camPlane.a.x) / segNumber
	segLenY := (camPlane.b.y - camPlane.a.y) / segNumber

	for i := 0.0; i <= segNumber; i++ { //for each segment

		var r vector //ray point on camPlane
		r.x = (camPlane.a.x + (segLenX * i))
		r.y = (camPlane.a.y + (segLenY * i))
		r = norm(r) //ray unit vector

		res := g.DDA(r)                                                                                                                                       //calculate ray length
		ebitenutil.DrawLine(g.screenBuffer, g.playerPos.x, g.playerPos.y, g.playerPos.x+res.x, g.playerPos.y+res.y, color.RGBA{242, 207, 85, 200} /*Yellow*/) //draw ray
	}

	//player's location in map
	//fmt.Println(divide(subtract(g.playerPos, g.mapPos), cellSize))

	//player view line drawing
	ebitenutil.DrawLine(g.screenBuffer, g.playerPos.x, g.playerPos.y, g.playerPos.x+multiply(g.viewDir, viewLen).x, g.playerPos.y+multiply(g.viewDir, viewLen).y, color.RGBA{255, 146, 28, 200} /*Orange*/)

	//camera plane line drawing
	ebitenutil.DrawLine(g.screenBuffer, g.playerPos.x+camPlane.a.x, g.playerPos.y+camPlane.a.y, g.playerPos.x+camPlane.b.x, g.playerPos.y+camPlane.b.y, color.RGBA{132, 132, 255, 200} /*Blue*/)
	//adding player position to convert from player coordinates to world coordinates

	//###
	// # DRAWING SCREEN #
	var opts ebiten.DrawImageOptions                    //declaring screen operations
	opts.GeoM.Translate(-g.playerPos.x, -g.playerPos.y) // converting screen world coordinates to player's coordinates

	//rotation
	var m ebiten.GeoM //declaring matrix
	/*
		00 10
		a  b
		01 11
		c  d
	*/
	//setting matrix
	//90° rotation = [x,y] -> [-y,x]
	m.SetElement(0, 0, -g.viewDir.y) //a (-| viewdir.x)
	m.SetElement(0, 1, g.viewDir.x)  //b (-| viewdir.y)
	m.SetElement(1, 0, g.viewDir.x)  //c (viewdir.x)
	m.SetElement(1, 1, g.viewDir.y)  //d (viewdir.y)
	m.Invert()                       //taking inverse matrix
	opts.GeoM.Concat(m)              //multiplying "opts matrix" with "our matrix"
	opts.GeoM.Scale(1, -1)           // scaling matrix for proper player movement & rotation

	opts.GeoM.Translate(float64(screen.Bounds().Max.X)/2, float64(screen.Bounds().Max.Y)/2) //centering the screen
	screen.DrawImage(g.screenBuffer, &opts)                                                 //drawing screen buffer

}

//-------------------------Functions----------------------------------

//DDA - calculates length of the ray
//inputs unit (direction) vector.
//outputs vector in player's coord
func (g Game) DDA(v vector) (res vector) {

	initV := v

	k := v.y / v.x

	var signX, signY float64
	if v.x < 0 { //if vector's x in player coord is negative
		signX = -1
	} else {
		signX = 1
	}
	if v.y < 0 { //if vector's x in player coord is negative
		signY = -1
	} else {
		signY = 1
	}

	pi := divide(subtract(g.playerPos, g.mapPos), cellSize) // position index

	fi := subtract(pi, vector{float64(int(pi.x)), float64(int(pi.y))}) // fraction index

	/*
		PLAN:
		length from pp to y edge
		length from pp to x edge
		choose which is shorter
		make pythagor to find length
		repeat

		but it's not working normally ;(
	*/

	var rayLen float64 //length of the ray

	// Start iteration
	A := fi.x * cellSize
	B := fi.y * cellSize

	fx := A * k
	fy := B / k

	lx := math.Sqrt((A * A) + (fx * fx))
	ly := math.Sqrt((B * B) + (fy * fy))

	if lx < ly {
		v = add(v, vector{lx, lx})
		rayLen += lx
		ly -= lx
		lx = math.Sqrt(1 + k*k) //calculating new lx
		pi.x += 1 * signX
	} else {
		v = add(v, vector{ly, ly})
		rayLen += ly
		lx -= ly
		ly = math.Sqrt(1 + 1/(k*k)) //calculating new ly
		pi.y += 1 * signY
	}

	wall := g.gameMap[int(pi.y)][int(pi.x)]

	for wall == 0 {
		// Following iterations
		if lx < ly { //for X edge
			v = multiply(add(v, vector{lx, lx}), signX) //adding length to vector
			rayLen += lx                                //increasing our length
			ly -= lx                                    //substracting our ly to match current position
			lx = math.Sqrt(1 + k*k)                     //calculating new lx
			pi.x += 1 * signX
		} else { // for Y edge
			v = multiply(add(v, vector{ly, ly}), signY)
			rayLen += ly                //increasing our length
			lx -= ly                    //substracting lx to match current position
			ly = math.Sqrt(1 + 1/(k*k)) //calculating new ly
			pi.y += 1 * signY
		}
		wall = g.gameMap[int(pi.y)][int(pi.x)]
	}

	return multiply(initV, rayLen) //make final vector
}

//normalize - converts vector to unit vector
func norm(v vector) (res vector) {
	mod := mod(v)
	res.x = v.x / mod
	res.y = v.y / mod
	return res
}

//returns module (length) of the vector
func mod(v vector) float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y))
}

func rotate(p vector, angle float64) (res vector) {

	//Rotation
	res.x = p.x*math.Cos(angle) - p.y*math.Sin(angle)
	res.y = p.x*math.Sin(angle) + p.y*math.Cos(angle)

	return res
}

func subtract(a, b vector) (res vector) {
	res.x = a.x - b.x
	res.y = a.y - b.y
	return res
}

func add(a, b vector) (res vector) {
	res.x = a.x + b.x
	res.y = a.y + b.y
	return res
}

func multiply(a vector, v float64) (res vector) { //v - value
	res.x = a.x * v
	res.y = a.y * v
	return res
}

func divide(a vector, v float64) (res vector) { //v - value
	res.x = a.x / v
	res.y = a.y / v
	return res
}

// draw cell with proper color
func drawCell(screen *ebiten.Image, mapPos vector, ci, cj int, clr color.RGBA) { //ci & cj - cell index in gameMap

	//  map position  +  cell position
	cX := mapPos.x + float64((cj * cellSize))
	cY := mapPos.y + float64((ci * cellSize))
	ebitenutil.DrawRect(screen, cX, cY, cellSize, cellSize, clr)
}

func initGameMap() [][]int {
	return [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1},
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

	mapPos := vector{(sW / 2) - (mW / 2), (sH / 2) - (mH / 2)} // map position
	playerStartPos := vector{(sW / 2) - 10, (sH / 2) - 10}     //player initial position
	viewDir := vector{0, -1 /*prev: -viewLen*/}                //player view direction unit vector

	return &Game{width: width, height: height, gameMap: initGameMap(), mapPos: mapPos, playerPos: playerStartPos, viewDir: viewDir /*camPlane: camPlane*/}
}
