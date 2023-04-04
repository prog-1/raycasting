package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//-------------------------Copyright---------------------------------
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
)

type Game struct {
	width, height int //screen width and height
	//global variables
	gameMap   [][]int //game map 2d matrix
	mapPos    vector  //map position on the screen (top left corner)
	playerPos vector  //player position
	playerDir vector  //player view direction
	//camPlane  line    // camera plane(screen)
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
		if g.gameMap[int((g.playerPos.y-1-g.mapPos.y)/cellSize)][int((g.playerPos.x-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos.y--

		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		//collision handling
		if g.gameMap[int((g.playerPos.y+1-g.mapPos.y)/cellSize)][int((g.playerPos.x-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos.y++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		//collision handling
		if g.gameMap[int((g.playerPos.y-g.mapPos.y)/cellSize)][int((g.playerPos.x-1-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos.x--
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		//collision handling
		if g.gameMap[int((g.playerPos.y-g.mapPos.y)/cellSize)][int((g.playerPos.x+1-g.mapPos.x)/cellSize)] == 0 {
			g.playerPos.x++
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.playerDir = rotate(g.playerDir, -math.Pi/200)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.playerDir = rotate(g.playerDir, math.Pi/200)
	}

	return nil
}

//---------------------------Draw-------------------------------------

func (g *Game) Draw(screen *ebiten.Image) {

	//draw map background
	ebitenutil.DrawRect(screen, g.mapPos.x, g.mapPos.y, mW, mH, color.RGBA{50, 50, 50, 255} /*Grey*/)

	//draw map cells
	for i := range g.gameMap { //for each column

		for j := range g.gameMap[i] { // for each row

			//If current cell is wall
			if g.gameMap[i][j] == 1 {
				drawCell(screen, g.mapPos, i, j, color.RGBA{200, 200, 200, 255} /*Light grey*/)

			} else if g.gameMap[i][j] == 2 {
				drawCell(screen, g.mapPos, i, j, color.RGBA{255, 0, 0, 255} /*Red*/)
			}

		}
	}

	//Draw player
	ebitenutil.DrawCircle(screen, g.playerPos.x, g.playerPos.y, 10, color.RGBA{100, 180, 255, 255} /*Light blue*/)

	//# FIELD OF VIEW #

	camA := add(g.playerDir, rotate(g.playerDir, -math.Pi/2)) //first camera point
	camB := add(g.playerDir, rotate(g.playerDir, math.Pi/2))  //last camera point
	camPlane := line{camA, camB}                              //camera plane line

	//# Ray Drawing #
	//segment length for each projection
	segLenX := (camPlane.b.x - camPlane.a.x) / segNumber
	segLenY := (camPlane.b.y - camPlane.a.y) / segNumber

	for i := 0.0; i <= segNumber; i++ { //for each segment
		ebitenutil.DrawLine(screen, g.playerPos.x, g.playerPos.y, g.playerPos.x+(camPlane.a.x+(segLenX*i)), g.playerPos.y+(camPlane.a.y+(segLenY*i)), color.RGBA{242, 207, 85, 200} /*Yellow*/) //draw ray
	}

	//player view line drawing
	ebitenutil.DrawLine(screen, g.playerPos.x, g.playerPos.y, g.playerPos.x+g.playerDir.x, g.playerPos.y+g.playerDir.y, color.RGBA{255, 146, 28, 200} /*Orange*/)

	//camera plane line drawing
	ebitenutil.DrawLine(screen, g.playerPos.x+camPlane.a.x, g.playerPos.y+camPlane.a.y, g.playerPos.x+camPlane.b.x, g.playerPos.y+camPlane.b.y, color.RGBA{132, 132, 255, 200} /*Blue*/)

}

//-------------------------Functions----------------------------------

func rotate(p vector, angle float64) (res vector) {

	//Rotation
	res.x = p.x*math.Cos(angle) - p.y*math.Sin(angle)
	res.y = p.x*math.Sin(angle) + p.y*math.Cos(angle)

	return res
}

func len(v vector) float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y))
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
	playerStartPos := vector{sW / 2, sH / 2}                   //player initial position

	playerDir := vector{0, -viewLen} //player view direction
	//camPlane := line{vector{playerDir.x - viewLen, playerDir.y}, vector{playerDir.x + viewLen, playerDir.y}} //camera plane line

	return &Game{width: width, height: height, gameMap: initGameMap(), mapPos: mapPos, playerPos: playerStartPos, playerDir: playerDir /*camPlane: camPlane*/}
}