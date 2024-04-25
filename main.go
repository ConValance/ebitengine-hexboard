package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)


const screenWidth  = 1200
const screenHeight = 800
const maxAngle     = 256

const columns = 10
const rows = 6
const tilewidth = 128
const tilesizex = 110
const sizex = tilesizex/2
const tilesizey = 94
const sizey = tilesizey/2
const floor1start = 4

var terrainimages []*ebiten.Image
var terrainimagenames [8]string
var spriteimages []*ebiten.Image
var spriteimagenames [8]string

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

	spriteimagenames[0] = "gopher.png"
	spriteimagenames[1] = "gopher1.png"
	for i := 0; i < 2; i++ {
		var err error
		var tmpimage *ebiten.Image
		var tmpstring string
		tmpstring = "resources/sprites/" + spriteimagenames[i]
		tmpimage, _, err = ebitenutil.NewImageFromFile(tmpstring)
		if err != nil {
			log.Fatal(err)
		}
		spriteimages = append(spriteimages, tmpimage)
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
	sprites  Sprites
	inited   bool
	op       ebiten.DrawImageOptions
}

type point struct {
	x float64
	y float64
}

type hex struct {
	q int
	r int
}

type cube struct {
	q int
	r int
	s int
}

type Orientation struct {
    f0, f1, f2, f3 float64
    b0, b1, b2, b3 float64
};


type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
	image		int
}

func (s *Sprite) Update() {
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if mx := screenWidth - s.imageWidth; mx <= s.x {
		s.x = 2*mx - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if my := screenHeight - s.imageHeight; my <= s.y {
		s.y = 2*my - s.y
		s.vy = -s.vy
	}
	//s.angle++
	//if s.angle == maxAngle {
	//	s.angle = 0
	//}
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for i := 0; i < s.num; i++ {
		s.sprites[i].Update()
	}
}

const (
	MinSprites = 0
	MaxSprites = 50000
)


func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	mx, my := ebiten.CursorPosition()
	//if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	//	g.paint(g.canvasImage, mx, my)
	//	drawn = true
	//}
	g.cursor = pos{
		x: mx,
		y: my,
	}

	g.sprites.Update()

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


func hex_round(hq, hr, hs float64) hex {
    var q int = int(math.Round(hq));
    var r int = int(math.Round(hr));
    var s int = int(math.Round(hs));
    var q_diff float64 = math.Abs(float64(q) - hq);
    var r_diff float64= math.Abs(float64(r) - hr);
    var s_diff float64 = math.Abs(float64(s) - hs);
    if (q_diff > r_diff && q_diff > s_diff) {
        q = -r - s;
    } else if (r_diff > s_diff) {
        r = -q - s;
    } else {
        s = -q - r;
    }
    return hex{q-1, r-1}
}


func pixel_to_hex(p point) hex {
	var layout_flat Orientation
	layout_flat=Orientation{3.0 / 2.0, 0.0, math.Sqrt(3.0) / 2.0, math.Sqrt(3.0),
		2.0 / 3.0, 0.0, -1.0 / 3.0, math.Sqrt(3.0) / 3.0}
	
	//var pt point = point{(p.x - layout.origin.x) / layout.size.x, (p.y - layout.origin.y) / layout.size.y}
	var pt point = point{(p.x)/sizex, (p.y)/sizey }
	var q float64 = layout_flat.b0 * pt.x + layout_flat.b1 * pt.y
	//var r float64 = layout_flat.b2 * pt.x  + layout_flat.b3 * pt.y
	var r float64 = layout_flat.b3 * pt.y
	if int(q)%2 == 0 {
		r=r*0.85
	}
	//return hex{int(q), int(r)}
	return hex_round(q, r, -q -r)
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	drawHex(screen)
	var p point
	var hx hex

	p.x = float64(g.cursor.x)
	p.y = float64(g.cursor.y)
	hx = pixel_to_hex(p)
	msg := fmt.Sprintf("mouseposition (%d, %d) = tile(%d, %d)", g.cursor.x, g.cursor.y, hx.q, hx.r)
	
	op.GeoM.Translate(float64(tilesizex*3/4)*float64(hx.q), float64(hx.r)*tilesizey)
	if hx.q%2 == 0 {
		op.GeoM.Translate(0, float64(tilesizey/2))
	}
	screen.DrawImage(terrainimages[6], op)

	//w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	var w,h int
	for i := 0; i < g.sprites.num; i++ {
		w = spriteimages[i].Bounds().Dx()
		h = spriteimages[i].Bounds().Dy()
		s := g.sprites.sprites[i]
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		g.op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / maxAngle)
		g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
		g.op.GeoM.Translate(float64(s.x), float64(s.y))
		screen.DrawImage(spriteimages[g.sprites.sprites[i].image], &g.op)
	}


	ebitenutil.DebugPrint(screen, msg)

}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	g.sprites.sprites = make([]*Sprite, 2)
	g.sprites.num = 2
	
	w, h := spriteimages[0].Bounds().Dx(), spriteimages[0].Bounds().Dy()
	x, y := 160, 310
	vx, vy := 0, 0
	a := 0
	g.sprites.sprites[0] = &Sprite{
		imageWidth:  w,
		imageHeight: h,
		x:           x,
		y:           y,
		vx:          vx,
		vy:          vy,
		angle:       a,
		image:		 0,
	}

	x, y = 410, 254
	vx, vy = 0, 0
	a = 0
	g.sprites.sprites[1] = &Sprite{
		imageWidth:  w,
		imageHeight: h,
		x:           x,
		y:           y,
		vx:          vx,
		vy:          vy,
		angle:       a,
		image: 		 1,
	}

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
