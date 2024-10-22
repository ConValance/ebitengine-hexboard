package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"hexboard/astarhexlib"
)

const screenWidth = 800
const screenHeight = 600
const maxAngle = 256

const columns = 10
const rows = 6
const tilewidth = 128
const tilesizex = 110
const sizex = tilesizex / 2
const tilesizey = 94
const sizey = tilesizey / 2
const floor1start = 4

var ngrid [][]*astarhexlib.ANode

var path []*astarhexlib.ANode
var vstart, vend astarhexlib.Vector2
var astar *astarhexlib.AStar

var terrainimages []*ebiten.Image
var terrains [16]Terrain
var spriteimages []*ebiten.Image
var spriteimagenames [16]string

var terrainmap0 = [rows][columns]int{
	{0, 0, 0, 3, 1, 2, 2, 0, 0, 0},
	{1, 1, 2, 1, 0, 2, 2, 2, 0, 0},
	{0, 0, 1, 1, 0, 0, 2, 0, 0, 0},
	{0, 0, 0, 1, 1, 0, 0, 0, 0, 3},
	{0, 0, 2, 0, 1, 1, 0, 3, 3, 3},
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
	{0, 7, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 5, 5, 0, 0, 0, 0, 0},
	{0, 6, 5, 0, 0, 0, 0, 7, 7, 0},
	{0, 4, 0, 0, 0, 0, 0, 0, 7, 0},
	{7, 4, 0, 0, 0, 0, 0, 0, 0, 0},
}
var flip1 = [rows][columns]int{ // flipx=1, flipy=2, both=3
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 3, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func init() {

	terrains[0].name = "Grass"
	terrains[0].filename = "grass.png"
	terrains[0].walkable = true

	terrains[1].name = "Water"
	terrains[1].filename = "water.png"
	terrains[1].walkable = false

	terrains[2].name = "Mountain"
	terrains[2].filename = "mountain.png"
	terrains[2].walkable = false

	terrains[3].name = "Desert"
	terrains[3].filename = "desert.png"
	terrains[3].walkable = true

	terrains[4].name = "Road1"
	terrains[4].filename = "road1.png"
	terrains[4].walkable = true

	terrains[5].name = "Road2"
	terrains[5].filename = "road2.png"
	terrains[5].walkable = true

	terrains[6].name = "Road3"
	terrains[6].filename = "road3.png"
	terrains[6].walkable = true

	terrains[7].name = "Tree"
	terrains[7].filename = "tree.png"
	terrains[7].walkable = true

	terrains[8].name = "Selection"
	terrains[8].filename = "selection.png"
	terrains[8].walkable = true

	terrains[9].name = "Pathfinder"
	terrains[9].filename = "pathfinder.png"
	terrains[9].walkable = true

	for i := 0; i < 10; i++ {
		var err error
		var tmpimage *ebiten.Image
		var tmpstring string
		tmpstring = "resources/terrain/" + terrains[i].filename
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

	genGrid()
	astar = astarhexlib.NewAStar(ngrid)

	fmt.Println("AStar colums, rows:", astar.GridCols, astar.GridRows)

	vstart.X = 1
	vstart.Y = 2
	vend.X = 7
	vend.Y = 0
	path = astar.FindPath(vstart, vend)
	if path != nil {
		var step int
		for i := 0; i < len(path); i++ {
			fmt.Println("path Nr:", i, " x:", path[i].Position.X, " y:", path[i].Position.Y)
			step = i
		}
		fmt.Println("and last step Nr:", step+1, " x:", vstart.X, " y:", vstart.Y)
	} else {
		fmt.Println("no path found!")
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
	sprites Sprites
	inited  bool
	op      ebiten.DrawImageOptions
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
}

type Terrain struct {
	name     string
	filename string
	walkable bool
}

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
	image       int
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
	g.cursor = pos{
		x: mx,
		y: my,
	}

	var p point
	var hx, sx hex
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.x = float64(g.cursor.x)
		p.y = float64(g.cursor.y)
		hx = pixel_to_hex(p)
		for i := 0; i < g.sprites.num; i++ {
			p.x = float64(g.sprites.sprites[i].x + tilesizex)
			p.y = float64(g.sprites.sprites[i].y - tilesizey/2)
			sx = pixel_to_hex(p)
			if int(sx.q)%2 != 0 {
				sx.r++
			}
			if hx == sx {
				fmt.Println("sprite clicked ", hx.q, hx.r)
			}
			p = hex_to_pixel(hx)
			//fmt.Println("pixelpos: ", p.x," ", p.y)
			g.setspritepos(0, int(p.x), int(p.y))
			vstart.X = float64(hx.q)
			vstart.Y = float64(hx.r)
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		p.x = float64(g.cursor.x)
		p.y = float64(g.cursor.y)
		hx = pixel_to_hex(p)
		//vstart.X = 1
		//vstart.Y = 2
		vend.X = float64(hx.q)
		vend.Y = float64(hx.r)
		//astar := NewAStar(ngrid)

		path = astar.FindPath(vstart, vend)
		if path != nil {
			hx.q = int(path[0].Position.X)
			hx.r = int(path[0].Position.Y)
			p = hex_to_pixel(hx)
			g.setspritepos(1, int(p.x), int(p.y))
		} else {
			fmt.Println("no path found!")
		}
	}

	g.sprites.Update()

	return nil
}

func genGrid() {
	for x := 0; x < (columns); x++ {
		var anodes []*astarhexlib.ANode
		for y := 0; y < (rows); y++ {
			var tmpanode astarhexlib.ANode
			tmpanode.Position.X = float64(x)
			tmpanode.Position.Y = float64(y)
			tmpanode.Walkable = terrains[terrainmap0[y][x]].walkable
			if terrainmap1[y][x] > 0 {
				tmpanode.Walkable = terrains[terrainmap1[y][x]].walkable
			}
			anodes = append(anodes, &tmpanode)
		}
		ngrid = append(ngrid, anodes)
	}
}

func drawHex(screen *ebiten.Image, floor int) {
	op := &ebiten.DrawImageOptions{}
	var flip int
	// flattop horizontal layout shoves oddq
	// floor 0
	for y := 0; y < (rows); y++ {
		for x := 0; x < (columns); x++ {
			switch floor {
			case 0:
				flip = flip0[y][x]
			case 1:
				flip = flip1[y][x]
			}
			op.GeoM.Reset()
			if flip > 0 {
				if flip == 1 {
					op.GeoM.Scale(-1, 1)
					op.GeoM.Translate(tilewidth, 0)
				} else {
					if flip == 2 {
						op.GeoM.Scale(1, -1)
						op.GeoM.Translate(0, tilewidth)
					} else {
						if flip == 3 {
							op.GeoM.Scale(-1, -1)
							op.GeoM.Translate(tilewidth, tilewidth)
						}
					}
				}
			}
			op.GeoM.Translate(float64(tilesizex*3/4)*float64(x), float64(y)*tilesizey)
			if x%2 != 0 {
				op.GeoM.Translate(0, float64(tilesizey/2))
			}
			switch floor {
			case 0:
				screen.DrawImage(terrainimages[terrainmap0[y][x]], op)
			case 1:
				if terrainmap1[y][x] >= floor1start {
					screen.DrawImage(terrainimages[terrainmap1[y][x]], op)
				}
			}

		}
	}

}

func hex_round(hq, hr, hs float64) hex {
	var q int = int(math.Round(hq))
	var r int = int(math.Round(hr))
	var s int = int(math.Round(hs))
	var q_diff float64 = math.Abs(float64(q) - hq)
	var r_diff float64 = math.Abs(float64(r) - hr)
	var s_diff float64 = math.Abs(float64(s) - hs)
	if q_diff > r_diff && q_diff > s_diff {
		q = -r - s
	} else if r_diff > s_diff {
		r = -q - s
	} else {
		s = -q - r
	}
	return hex{q - 1, r - 1}
}

func pixel_to_hex(p point) hex {
	var layout_flat Orientation
	layout_flat = Orientation{3.0 / 2.0, 0.0, math.Sqrt(3.0) / 2.0, math.Sqrt(3.0),
		2.0 / 3.0, 0.0, -1.0 / 3.0, math.Sqrt(3.0) / 3.0}

	//var pt point = point{(p.x - layout.origin.x) / layout.size.x, (p.y - layout.origin.y) / layout.size.y}
	var pt point = point{(p.x) / sizex, (p.y) / sizey}
	var q float64 = layout_flat.b0*pt.x + layout_flat.b1*pt.y
	//var r float64 = layout_flat.b2 * pt.x  + layout_flat.b3 * pt.y
	var r float64 = layout_flat.b3 * pt.y * 0.85

	if int(q)%2 != 0 {
		r = r - 0.5
	}
	//return hex{int(q), int(r)}
	return hex_round(q, r, -q-r)
}

func hex_to_pixel(h hex) point {
	var x float64 = float64(tilesizex*3/4) * float64(h.q)
	var y float64 = float64(h.r) * float64(tilesizey)
	if h.q%2 != 0 {
		y = y + float64(tilesizey/2)
	}
	return point{x, y}
}

func get_terrain_nr(n string) int {
	for i := 0; i < 16; i++ {
		if terrains[i].name == n {
			return i
		}
	}
	return -1
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// draw the hexboard
	drawHex(screen, 0)

	// draw the green steps of the pathfinding
	for i := 0; i < len(path); i++ {
		//fmt.Println("path Nr:", i, " x:", path.items[i].Position.X, " y:", path.items[i].Position.Y)
		op.GeoM.Reset()
		op.GeoM.Translate(float64(tilesizex*3/4)*float64(path[i].Position.X), float64(path[i].Position.Y)*tilesizey)
		if int(path[i].Position.X)%2 != 0 {
			op.GeoM.Translate(0, float64(tilesizey/2))
		}
		screen.DrawImage(terrainimages[9], op)
		//step = i
	}
	//fmt.Println("and last step Nr:", step+1, " x:", vstart.X, " y:", vstart.Y)
	op.GeoM.Reset()
	op.GeoM.Translate(float64(tilesizex*3/4)*float64(vstart.X), float64(vstart.Y)*tilesizey)
	if int(vstart.X)%2 != 0 {
		op.GeoM.Translate(0, float64(tilesizey/2))
	}
	screen.DrawImage(terrainimages[9], op)

	// draw the mousepos and the hextile to screen
	var p point
	var hx hex
	p.x = float64(g.cursor.x)
	p.y = float64(g.cursor.y + tilesizey/2)
	hx = pixel_to_hex(p)
	msg := fmt.Sprintf("mouseposition (%d, %d) =tile(%d, %d) - left mouseclick for new startposition, right click for new targetposition.", g.cursor.x, g.cursor.y, hx.q, hx.r)

	// draw the selected hextile
	op.GeoM.Reset()
	op.GeoM.Translate(float64(tilesizex*3/4)*float64(hx.q), float64(hx.r)*tilesizey)
	if hx.q%2 != 0 {
		op.GeoM.Translate(0, float64(tilesizey/2))
	}
	screen.DrawImage(terrainimages[8], op)

	// draw the hexboard
	drawHex(screen, 1)

	//w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	var w, h int
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

func (g *Game) setspritepos(nr int, xp int, yp int) {
	w, h := spriteimages[0].Bounds().Dx(), spriteimages[0].Bounds().Dy()
	vx, vy := 0, 0
	a := 0
	g.sprites.sprites[nr] = &Sprite{
		imageWidth:  w,
		imageHeight: h,
		x:           xp,
		y:           yp - 25,
		vx:          vx,
		vy:          vy,
		angle:       a,
		image:       nr,
	}
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()

	g.sprites.sprites = make([]*Sprite, 2)
	g.sprites.num = 2

	w, h := spriteimages[0].Bounds().Dx(), spriteimages[0].Bounds().Dy()
	x, y := 80, 210
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
		image:       0,
	}

	x, y = 580, 20
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
		image:       1,
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1200, 800
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("ebitengine-hexboard")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
