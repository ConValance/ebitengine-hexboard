package main

import (
	_ "image/png"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)


const columns = 10
const rows = 6
const tilewidth = 128
const tilesizex = 110
const tilesizey = 94
const floor1start = 4


var terrainimages []*ebiten.Image
var terrainimagenames [8]string

var terrainmap0 = [rows][columns]int{
	{ 0, 0, 0, 3, 1, 2, 2, 0, 0, 0},
	{ 1, 1, 2, 3, 0, 2, 2, 2, 0, 0},
	{ 1, 1, 1, 0, 0, 0, 2, 0, 0, 0},
	{ 0, 0, 0, 1, 1, 0, 0, 0, 0, 3},
	{ 0, 0, 0, 2, 2, 1, 0, 3, 3, 3},
	{ 0, 0, 0, 0, 0, 1, 1, 0, 3, 3},
}
var flip0 = [rows][columns]int{ // flipx=1, flipy=2, both=3
	{ 0, 1, 0, 0, 0, 1, 0, 0, 1, 0},
	{ 1, 2, 0, 1, 0, 0, 0, 1, 0, 0},
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{ 0, 1, 0, 1, 2, 0, 0, 1, 0, 0},
	{ 1, 0, 0, 0, 1, 0, 0, 0, 1, 0},
	{ 0, 0, 1, 0, 0, 1, 2, 0, 0, 1},
}

var terrainmap1 = [rows][columns]int{
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{ 0, 0, 0, 0, 5, 6, 0, 0, 0, 0},
	{ 0, 0, 5, 5, 0, 5, 0, 7, 0, 0},
	{ 0, 4, 0, 0, 0, 0, 0, 0, 0, 0},
	{ 0, 4, 0, 0, 0, 0, 0, 0, 0, 0},
}
var flip1 = [rows][columns]int{ // flipx=1, flipy=2, both=3
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{ 0, 0, 0, 3, 0, 0, 0, 0, 0, 0},
	{ 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func init() {

	terrainimagenames[0] = "grass.png"
	terrainimagenames[1] = "water.png"
	terrainimagenames[2] = "mountain.png"
	terrainimagenames[3] = "desert.png"
	terrainimagenames[4] = "road1.png"
	terrainimagenames[5] = "road2.png"
	terrainimagenames[6] = "companygreen.png"
	terrainimagenames[7] = "companyred.png"

	for i := 0; i < 8; i++ {
		var err error
		var tmpimage *ebiten.Image
		var tmpstring string
		tmpstring = "resources/terrain/" + terrainimagenames[i]
		tmpimage, _, err = ebitenutil.NewImageFromFile(tmpstring)
		if err != nil {
			log.Fatal(err)
		}
		terrainimages = append(terrainimages, tmpimage)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func drawHex(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// floor 0
	for y := 0; y < (rows); y++ {
		for x := 0; x < (columns); x++ {
			op.GeoM.Reset()

			if flip0[y][x] > 0 {
				if flip0[y][x]== 1 {
					op.GeoM.Scale(-1, 1)
					op.GeoM.Translate(tilewidth, 0)
				} else {
					if flip0[y][x]== 2 {
						op.GeoM.Scale(1, -1)
						op.GeoM.Translate(0 , tilewidth)
					} else {
						if flip0[y][x]== 3 {
							op.GeoM.Scale(-1, -1)
							op.GeoM.Translate(tilewidth , tilewidth)
						}
					}
				}
			}
			op.GeoM.Translate(float64(tilesizex*3/4)*float64(x), float64(y)*tilesizey)
			if x%2 == 0 {
				op.GeoM.Translate(0, float64(tilesizey/2))
			}

			screen.DrawImage(terrainimages[terrainmap0[y][x]], op)
		}
	}

	// floor 1
	for y := 0; y < (rows); y++ {
		for x := 0; x < (columns); x++ {
			op.GeoM.Reset()

			if flip1[y][x] > 0 {
				if flip1[y][x]== 1 {
					op.GeoM.Scale(-1, 1)
					op.GeoM.Translate(tilewidth, 0)
				} else {
					if flip1[y][x]== 2 {
						op.GeoM.Scale(1, -1)
						op.GeoM.Translate(0 , tilewidth)
					} else {
						if flip1[y][x]== 3 {
							op.GeoM.Scale(-1, -1)
							op.GeoM.Translate(tilewidth , tilewidth)
						}
					}
				}
			}
			op.GeoM.Translate(float64(tilesizex*3/4)*float64(x), float64(y)*tilesizey)
			if x%2 == 0 {
				op.GeoM.Translate(0, float64(tilesizey/2))
			}
			if terrainmap1[y][x]>=floor1start {
				screen.DrawImage(terrainimages[terrainmap1[y][x]], op)
			}
		}
	}


}

func (g *Game) Draw(screen *ebiten.Image) {
	drawHex(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1200, 800
}

func main() {
	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("hexboard ebitengine")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
