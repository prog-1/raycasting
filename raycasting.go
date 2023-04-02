package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//---------------------------Plan-----------------------------------
/*

- map (rectangle) with size 600x600px ✔
- 2d matrix of map ✔
- 24x24 grid | each cell 25x25px ✔
- wall filling depnding on matrix ✔


- player point (just circle)
- player straight view ray
- player field of view rays from formula
- player rays visualization

*/
//-------------------------Declaration------------------------------

const (
	sW       = 1280 //screen width
	sH       = 720  //screen height
	mW       = 600  //map width
	mH       = 600  //map height
	cellSize = 25
)

type Game struct {
	width, height int //screen width and height
	//global variables
	gameMap     [][]int //game map 2d matrix
	playerStart point   //player start position
}

type point struct {
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
	ebitenutil.DrawRect(screen, (sW/2)-(mW/2), (sH/2)-(mH/2), mW, mH, color.RGBA{50, 50, 50, 255})

	//draw map cells
	for i := range g.gameMap { //for each column

		for j := range g.gameMap[i] { // for each row

			//If current cell is wall
			if g.gameMap[i][j] == 1 {
				drawCell(screen, i, j, color.RGBA{255, 255, 255, 255}) //White

			} else if g.gameMap[i][j] == 2 {
				drawCell(screen, i, j, color.RGBA{255, 0, 0, 255}) //Red
			}

		}
	}

	//draw player
	ebitenutil.DrawCircle(screen, g.playerStart.x, g.playerStart.y, 10, color.RGBA{100, 180, 255, 255})
}

//-------------------------Functions----------------------------------

//draw cell with proper color
func drawCell(screen *ebiten.Image, ci, cj int, clr color.RGBA) { //ci & cj - cell index in gameMap
	//       map location     +   cell position
	cX := (sW / 2) - (mW / 2) + (cj * cellSize)
	cY := (sH / 2) - (mH / 2) + (ci * cellSize)
	ebitenutil.DrawRect(screen, float64(cX), float64(cY), cellSize, cellSize, clr)
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

	playerStart := point{sW / 2, sH / 2}

	return &Game{width: width, height: height, gameMap: initGameMap(), playerStart: playerStart}
}
