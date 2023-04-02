package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//---------------------------Plan-----------------------------------
/*

- map (rectangle) with size 600x600px ✔
- 2d matrix of map ✔
- 24x24 grid | each cell 25x25px ✔
- wall filling depnding on matrix ✔


- player point (just circle) ✔
- player straight view ray ✔
- camera plane ✔
- FoV rays ✔
- i = i-1 ✔

- player movement
- player rotation

- collisions

*/
//-------------------------Declaration------------------------------

const (
	sW        = 1280 //screen width
	sH        = 720  //screen height
	mW        = 600  //map width
	mH        = 600  //map height
	cellSize  = 25   // size of each cell in pixels
	segNumber = 100  // number of segments of camera plane vector for rays
)

type Game struct {
	width, height int //screen width and height
	//global variables
	gameMap     [][]int //game map 2d matrix
	mapPos      point   //map position on the screen (top left corner)
	playerStart point   //player start position
	viewDir     line    //player view direction vector (ray)
	camPlane    line    // camera 1/2 plane vector
}

type point struct {
	x, y float64
}

//line with start & end point
type line struct {
	a, b vector
}

type vector struct {
	x, y float64
}

//---------------------------Update-------------------------------------

func (g *Game) Update() error {
	//all logic on update
	//can be divided on seperate functions (i.e. "UpdateCircle")
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

	//draw player
	ebitenutil.DrawCircle(screen, g.playerStart.x, g.playerStart.y, 10, color.RGBA{100, 180, 255, 255} /*Light blue*/)

	//FOV:

	//RAYS
	//segment length for each projection
	segLenX := (g.camPlane.b.x - g.camPlane.a.x) / segNumber
	segLenY := (g.camPlane.b.y - g.camPlane.a.y) / segNumber

	for i := 0.0; i <= segNumber; i++ { //for each segment
		ebitenutil.DrawLine(screen, g.viewDir.a.x, g.viewDir.a.y, g.camPlane.a.x+(segLenX*i), g.camPlane.a.y+(segLenY*i), color.RGBA{242, 207, 85, 200} /*Yellow*/) //draw ray
	}

	//draw player view vector
	ebitenutil.DrawLine(screen, g.viewDir.a.x, g.viewDir.a.y, g.viewDir.b.x, g.viewDir.b.y, color.RGBA{255, 146, 28, 200} /*Orange*/)

	//draw camera plane vector
	ebitenutil.DrawLine(screen, g.camPlane.a.x, g.camPlane.a.y, g.camPlane.b.x, g.camPlane.b.y, color.RGBA{132, 132, 255, 200} /*Blue*/)

}

//-------------------------Functions----------------------------------

func len(v vector) float64 {
	return math.Sqrt((v.x * v.x) + (v.y * v.y))
}

func subtract(a, b vector) (res vector) {
	res.x = a.x - b.x
	res.y = a.y - b.y
	return res
}

//draw cell with proper color
func drawCell(screen *ebiten.Image, mapPos point, ci, cj int, clr color.RGBA) { //ci & cj - cell index in gameMap

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

//New game instance function
func NewGame(width, height int) *Game {

	mapPos := point{(sW / 2) - (mW / 2), (sH / 2) - (mH / 2)} // map position
	playerStart := point{sW / 2, sH / 2}                      //player start position

	//Field of View
	len := 150.0 //lenght of player view and camera plane vectors
	//if len viewDir = len camPlane, then FoV = 90 degrees
	viewDir := line{vector{playerStart.x, playerStart.y}, vector{playerStart.x, playerStart.y - len}} //player view direciton vector
	camPlane := line{vector{viewDir.b.x - len, viewDir.b.y}, vector{viewDir.b.x + len, viewDir.b.y}}  // camera plane vector

	return &Game{width: width, height: height, gameMap: initGameMap(), mapPos: mapPos, playerStart: playerStart, viewDir: viewDir, camPlane: camPlane}
}
