package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

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
	{0, 0, 0, 3, 1, 2, 2, 0, 0, 0},
	{1, 1, 2, 3, 0, 2, 2, 2, 0, 0},
	{1, 1, 1, 0, 0, 0, 2, 0, 0, 0},
	{0, 0, 0, 1, 1, 0, 0, 0, 0, 3},
	{0, 0, 0, 2, 2, 1, 0, 3, 3, 3},
	{0, 0, 0, 0, 0, 1, 1, 0, 3, 3},
}
var flip0 = [rows][columns]int{ // flipx=1, flipy=2, both=3
	{0, 1, 0, 0, 0, 1, 0, 0, 1, 0},
	{1, 2, 0, 1, 0, 0, 0, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 1, 0, 1, 2, 0, 0, 1, 0, 0},
	{1, 0, 0, 0, 1, 0, 0, 0, 1, 0},
	{0, 0, 1, 0, 0, 1, 2, 0, 0, 1},
}

var terrainmap1 = [rows][columns]int{
	{0, 0, 0, 0, 0, 0, 0, 7, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 5, 6, 0, 0, 0, 0},
	{0, 0, 5, 5, 0, 5, 0, 7, 0, 0},
	{0, 4, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}
var flip1 = [rows][columns]int{ // flipx=1, flipy=2, both=3
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 3, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
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

type touch struct {
	id  ebiten.TouchID
	pos pos
}

type pos struct {
	x int
	y int
}

type Game struct {
	cursor  pos
	touches []touch
	count   int
}

type point struct {
	x int
	y int
}

type hex struct {
	q float64
	r float64
}

type cube struct {
	q float64
	r float64
	s float64
}



func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	//	g.paint(g.canvasImage, mx, my)
	//	drawn = true
	//}
	g.cursor = pos{
		x: mx,
		y: my,
	}

	return nil
}

func drawHex(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// flattop “even-r” horizontal layout shoves even rows right
	// floor 0
	for y := 0; y < (rows); y++ {
		for x := 0; x < (columns); x++ {
			op.GeoM.Reset()

			if flip0[y][x] > 0 {
				if flip0[y][x] == 1 {
					op.GeoM.Scale(-1, 1)
					op.GeoM.Translate(tilewidth, 0)
				} else {
					if flip0[y][x] == 2 {
						op.GeoM.Scale(1, -1)
						op.GeoM.Translate(0, tilewidth)
					} else {
						if flip0[y][x] == 3 {
							op.GeoM.Scale(-1, -1)
							op.GeoM.Translate(tilewidth, tilewidth)
						}
					}
				}
			}
			op.GeoM.Translate(float64(tilesizex*3/4)*float64(x), float64(y)*tilesizey)
			if x%2 == 0 {
				op.GeoM.Translate(0, float64(tilesizey/2))
			}
			screen.DrawImage(terrainimages[terrainmap0[y][x]], op)

			op.GeoM.Reset()

			if flip1[y][x] > 0 {
				if flip1[y][x] == 1 {
					op.GeoM.Scale(-1, 1)
					op.GeoM.Translate(tilewidth, 0)
				} else {
					if flip1[y][x] == 2 {
						op.GeoM.Scale(1, -1)
						op.GeoM.Translate(0, tilewidth)
					} else {
						if flip1[y][x] == 3 {
							op.GeoM.Scale(-1, -1)
							op.GeoM.Translate(tilewidth, tilewidth)
						}
					}
				}
			}
			op.GeoM.Translate(float64(tilesizex*3/4)*float64(x), float64(y)*tilesizey)
			if x%2 == 0 {
				op.GeoM.Translate(0, float64(tilesizey/2))
			}
			if terrainmap1[y][x] >= floor1start {
				screen.DrawImage(terrainimages[terrainmap1[y][x]], op)
			}
		}
	}

	// floor 1
	/*
		for y := 0; y < (rows); y++ {
			for x := 0; x < (columns); x++ {
				op.GeoM.Reset()

				if flip1[y][x] > 0 {
					if flip1[y][x] == 1 {
						op.GeoM.Scale(-1, 1)
						op.GeoM.Translate(tilewidth, 0)
					} else {
						if flip1[y][x] == 2 {
							op.GeoM.Scale(1, -1)
							op.GeoM.Translate(0, tilewidth)
						} else {
							if flip1[y][x] == 3 {
								op.GeoM.Scale(-1, -1)
								op.GeoM.Translate(tilewidth, tilewidth)
							}
						}
					}
				}
				op.GeoM.Translate(float64(tilesizex*3/4)*float64(x), float64(y)*tilesizey)
				if x%2 == 0 {
					op.GeoM.Translate(0, float64(tilesizey/2))
				}
				if terrainmap1[y][x] >= floor1start {
					screen.DrawImage(terrainimages[terrainmap1[y][x]], op)
				}
			}
		}
	*/

}

func axial_to_cube(h hex) cube {
    var q = h.q
    var r = h.r
    var s = -q-r
    return cube{q, r, s}
}

func cube_round(c cube) cube {
    var q = math.Round(c.q)
    var r = math.Round(c.r)
    var s = math.Round(c.s)

    var q_diff = math.Abs(q - c.q)
    var r_diff = math.Abs(r - c.r)
    var s_diff = math.Abs(s - c.s)

    if q_diff > r_diff && q_diff > s_diff {
        q = -r-s
	} else if r_diff > s_diff {
        r = -q-s
	} else {
        s = -q-r
	}

    return cube{q, r, s}
}

func cube_to_axial(c cube) hex {
    var q = c.q
    var r = c.r
    return hex{q, r}
}

func axial_round(h hex) hex {
    return cube_to_axial(cube_round(axial_to_cube(h)))
}

func pixel_to_flat_hex(p point) hex {
    var q = ( 2.0/3.0 * float64(p.x) ) / (tilesizex/2)
    var r = (-1.0/3.0 * float64(p.x)  +  math.Sqrt(3.0)/3.0 * float64(p.y)) / (tilesizex/2)
    return axial_round(hex{q, r})
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawHex(screen)
	var p point
	var h hex
	
	p.x=g.cursor.x
	p.y=g.cursor.y
	h = pixel_to_flat_hex(p)
	msg := fmt.Sprintf("mouseposition (%d, %d)= tile(%2.f, %2.f)", g.cursor.x, g.cursor.y, h.q, h.r)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1200, 800
}

func main() {
	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("hexboard ebitengine")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
