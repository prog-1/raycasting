package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 960
	screenHeight = 720
	cellSize     = 8
	rayNum       = 100
	texWidth     = 64
	texHeight    = 64
)

var (
	worldMap = [25][25]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}
	clr = [5]color.Color{color.RGBA{255, 255, 255, 150}, color.RGBA{255, 0, 0, 150}, color.RGBA{0, 255, 0, 150}, color.RGBA{0, 0, 255, 150}, color.RGBA{0, 255, 255, 150}}
)

type Point struct {
	x, y float64
}

type Game struct {
	width, height   int
	pos, dir, plane *Point
	textures        [5][]color.Color
	buffer          [screenHeight][screenWidth]color.Color
	minimap         *ebiten.Image
	fisheye         bool
	lastTime        time.Time
}

func NewGame(width, height int) *Game {
	return &Game{
		width:    width,
		height:   height,
		pos:      &Point{13, 13},
		dir:      &Point{0, -1},
		plane:    &Point{0.5, 0},
		textures: loadTextures(),
		minimap:  ebiten.NewImage(len(worldMap[0])*cellSize, len(worldMap)*cellSize),
		lastTime: time.Now(),
	}
}

func loadTextures() (textures [5][]color.Color) {
	for i := 0; i < len(textures); i++ {
		textures[i] = make([]color.Color, texWidth*texHeight)
		loadImage(textures[i], fmt.Sprintf("texture%d.png", i+1))
	}
	return
}

func loadImage(texture []color.Color, path string) {
	_, img, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	for y := 0; y < texHeight; y++ {
		for x := 0; x < texWidth; x++ {
			c := img.At(x, y)
			texture[texWidth*y+x] = c
		}
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func rotate(p *Point, angle float64) *Point {
	res := new(Point)
	res.x = p.x*math.Cos(angle) - p.y*math.Sin(angle)
	res.y = p.x*math.Sin(angle) + p.y*math.Cos(angle)
	return res
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) && worldMap[int(g.pos.y+g.dir.y/20)][int(g.pos.x+g.dir.x/20)] == 0 {
		g.pos.x += g.dir.x / 20
		g.pos.y += g.dir.y / 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && worldMap[int(g.pos.y-g.dir.y/20)][int(g.pos.x-g.dir.x/20)] == 0 {
		g.pos.x -= g.dir.x / 20
		g.pos.y -= g.dir.y / 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && worldMap[int(g.pos.y-g.dir.x/20)][int(g.pos.x+g.dir.y/20)] == 0 {
		g.pos.x += g.dir.y / 20
		g.pos.y -= g.dir.x / 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && worldMap[int(g.pos.y+g.dir.x/20)][int(g.pos.x-g.dir.y/20)] == 0 {
		g.pos.x -= g.dir.y / 20
		g.pos.y += g.dir.x / 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.dir = rotate(g.dir, -math.Pi/90)
		g.plane = rotate(g.plane, -math.Pi/90)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.dir = rotate(g.dir, math.Pi/90)
		g.plane = rotate(g.plane, math.Pi/90)
	}
	if t := time.Now(); ebiten.IsKeyPressed(ebiten.KeyF) && t.Sub(g.lastTime).Milliseconds() > 300 {
		g.fisheye = !g.fisheye
		g.lastTime = t
	}
	return nil
}

func (g *Game) DrawMinimap(screen *ebiten.Image) {
	g.minimap.Clear()
	whiteRect, redRect, greenRect, blueRect, lightBlueRect := ebiten.NewImage(cellSize, cellSize), ebiten.NewImage(cellSize, cellSize), ebiten.NewImage(cellSize, cellSize), ebiten.NewImage(cellSize, cellSize), ebiten.NewImage(cellSize, cellSize)
	whiteRect.Fill(color.RGBA{255, 255, 255, 150})
	redRect.Fill(color.RGBA{255, 0, 0, 150})
	greenRect.Fill(color.RGBA{0, 255, 0, 150})
	blueRect.Fill(color.RGBA{0, 0, 255, 150})
	lightBlueRect.Fill(color.RGBA{0, 255, 255, 150})
	for i := range worldMap {
		for j := range worldMap[i] {
			if worldMap[j][i] != 0 {
				var opts ebiten.DrawImageOptions
				opts.GeoM.Translate(float64(i*cellSize), float64(j*cellSize))
				switch worldMap[j][i] {
				case 1:
					g.minimap.DrawImage(whiteRect, &opts)
				case 2:
					g.minimap.DrawImage(redRect, &opts)
				case 3:
					g.minimap.DrawImage(greenRect, &opts)
				case 4:
					g.minimap.DrawImage(blueRect, &opts)
				case 5:
					g.minimap.DrawImage(lightBlueRect, &opts)
				}
			}
		}
	}
	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(-g.pos.x*cellSize, -g.pos.y*cellSize)
	var f_1 ebiten.GeoM
	f_1.SetElement(0, 0, -g.dir.y)
	f_1.SetElement(0, 1, g.dir.x)
	f_1.SetElement(1, 0, g.dir.x)
	f_1.SetElement(1, 1, g.dir.y)
	f_1.Invert()
	opts.GeoM.Concat(f_1)
	opts.GeoM.Scale(1, -1)
	opts.GeoM.Translate(float64(g.minimap.Bounds().Dx())/2+4, float64(g.minimap.Bounds().Dy())/2+4)
	g.DrawRays()
	screen.DrawImage(g.minimap, &opts)
}

func (g *Game) DrawRays() {
	for i := 0.0; i < rayNum; i++ {
		cameraX := 2*i/float64(rayNum) - 1
		rayDir := Point{g.dir.x + g.plane.x*cameraX, g.dir.y + g.plane.y*cameraX}
		p := Point{math.Floor(g.pos.x * cellSize), math.Floor(g.pos.y * cellSize)}
		deltaDist := Point{math.Abs(1 / rayDir.x), math.Abs(1 / rayDir.y)}
		var sideDist, step Point
		if rayDir.x < 0 {
			step.x = -1
			sideDist.x = (g.pos.x*cellSize - p.x) * deltaDist.x
		} else {
			step.x = 1
			sideDist.x = (p.x + 1.0 - g.pos.x*cellSize) * deltaDist.x
		}
		if rayDir.y < 0 {
			step.y = -1
			sideDist.y = (g.pos.y*cellSize - p.y) * deltaDist.y
		} else {
			step.y = 1
			sideDist.y = (p.y + 1.0 - g.pos.y*cellSize) * deltaDist.y
		}
		for worldMap[int(p.y)/cellSize][int(p.x)/cellSize] == 0 {
			if sideDist.x < sideDist.y {
				sideDist.x += deltaDist.x
				p.x += step.x
			} else {
				sideDist.y += deltaDist.y
				p.y += step.y
			}
			g.minimap.Set(int(p.x), int(p.y), color.RGBA{255, 255, 0, 50})
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for x := 0; x < g.width; x++ {
		cameraX := 2*float64(x)/float64(g.width) - 1
		rayDir := Point{g.dir.x + g.plane.x*cameraX, g.dir.y + g.plane.y*cameraX}
		p := Point{math.Floor(g.pos.x), math.Floor(g.pos.y)}
		deltaDist := Point{math.Abs(1 / rayDir.x), math.Abs(1 / rayDir.y)}
		var sideDist, step Point
		var perpWallDist float64
		var side int
		if rayDir.x < 0 {
			step.x = -1
			sideDist.x = (g.pos.x - p.x) * deltaDist.x
		} else {
			step.x = 1
			sideDist.x = (p.x + 1.0 - g.pos.x) * deltaDist.x
		}
		if rayDir.y < 0 {
			step.y = -1
			sideDist.y = (g.pos.y - p.y) * deltaDist.y
		} else {
			step.y = 1
			sideDist.y = (p.y + 1.0 - g.pos.y) * deltaDist.y
		}
		for worldMap[int(p.y)][int(p.x)] == 0 {
			if sideDist.x < sideDist.y {
				sideDist.x += deltaDist.x
				p.x += step.x
				side = 0
			} else {
				sideDist.y += deltaDist.y
				p.y += step.y
				side = 1
			}
		}
		if side == 0 {
			perpWallDist = sideDist.x - deltaDist.x
		} else {
			perpWallDist = sideDist.y - deltaDist.y
		}
		if g.fisheye {
			perpWallDist *= math.Sqrt(1 + math.Pow(math.Sqrt(math.Pow(rayDir.x, 2)+math.Pow(rayDir.y, 2)), 2))
		}
		lineHeight := int(float64(g.height) / perpWallDist)
		drawStart := -lineHeight/2 + g.height/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + g.height/2
		if drawEnd >= g.height {
			drawEnd = g.height - 1
		}
		texNum := worldMap[int(p.y)][int(p.x)] - 1
		var wallX float64
		if side == 0 {
			wallX = g.pos.y + perpWallDist*rayDir.y
		} else {
			wallX = g.pos.x + perpWallDist*rayDir.x
		}
		wallX -= math.Floor(wallX)
		texX := int(wallX * float64(texWidth))
		if (side == 0 && rayDir.x > 0) || (side == 1 && rayDir.y < 0) {
			texX = texWidth - texX - 1
		}
		texStep := float64(texHeight) / float64(lineHeight)
		texPos := float64(drawStart-g.height/2+lineHeight/2) * texStep
		for y := drawStart; y < drawEnd; y++ {
			texY := int(texPos) & (texHeight - 1)
			texPos += texStep
			c := g.textures[texNum][texHeight*texY+texX]
			if side == 1 {
				r, g, b, a := c.RGBA()
				c = color.RGBA{uint8(float64(r>>8) * 0.5), uint8(float64(g>>8) * 0.5), uint8(float64(b>>8) * 0.5), uint8(float64(a >> 8))}
			}
			g.buffer[y][x] = c
		}
	}
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			if g.buffer[y][x] != nil {
				screen.Set(x, y, g.buffer[y][x])
				g.buffer[y][x] = nil
			}
		}
	}
	g.DrawMinimap(screen)
	vector.DrawFilledCircle(screen, float32(13)*cellSize, float32(13)*cellSize, 3, color.RGBA{255, 255, 0, 150}, false)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
