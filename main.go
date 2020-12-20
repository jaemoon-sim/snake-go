package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const WindowW = 640
const WindowH = 480

const tileSize = 20

var DeadChan chan interface{}

type Coord struct {
	X, Y int
}

func (c *Coord) Draw(board *ebiten.Image) {
	coordImage := ebiten.NewImage(tileSize, tileSize)
	coordImage.Fill(color.RGBA{0x88, 0x00, 0x88, 0xff})

	op := &ebiten.DrawImageOptions{}
	// bw, bh := board.Size()
	//
	// if c.X*tileSize > bw-tileSize {
	// c.X = 0
	// }
	// if c.X*tileSize < 0 {
	// c.X = bw / tileSize
	// }
	// if c.Y*tileSize > bh-tileSize {
	// c.Y = 0
	// }
	// if c.Y*tileSize < 0 {
	// c.Y = bh / tileSize
	// }
	//
	x := c.X * tileSize
	y := c.Y * tileSize

	op.GeoM.Translate(float64(x), float64(y))
	board.DrawImage(coordImage, op)
}

func (c *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

type Direction string

var (
	U Direction = "u"
	R Direction = "r"
	D Direction = "d"
	L Direction = "l"
)

type Snake struct {
	Direction
	body    []*Coord
	growing bool
}

func createSnake() *Snake {
	c0 := &Coord{2, 10}
	c1 := &Coord{1, 10}
	c2 := &Coord{0, 10}

	return &Snake{R, []*Coord{c0, c1, c2}, false}
}

func (s *Snake) Draw(board *ebiten.Image) {
	for _, c := range s.body {
		c.Draw(board)
	}
}

func (s *Snake) String() string {
	res := ""
	for _, c := range s.body {
		res += c.String()
		res += " "
	}
	return res
}

func (s *Snake) Move() {
	var prevX, prevY = s.body[0].X, s.body[0].Y
	for _, c := range s.body[1:] {
		x, y := c.X, c.Y
		c.X, c.Y = prevX, prevY
		prevX, prevY = x, y
	}
	if s.growing {
		c := &Coord{prevX, prevY}
		s.body = append(s.body, c)
		s.growing = false
	}
	switch s.Direction {
	case U:
		s.body[0].X = s.body[0].X
		s.body[0].Y = s.body[0].Y - 1
	case R:
		s.body[0].X = s.body[0].X + 1
		s.body[0].Y = s.body[0].Y
	case D:
		s.body[0].X = s.body[0].X
		s.body[0].Y = s.body[0].Y + 1
	case L:
		s.body[0].X = s.body[0].X - 1
		s.body[0].Y = s.body[0].Y
	}
}

func (s *Snake) CheckFood(f *Food) bool {
	if s.body[0].X == f.X && s.body[0].Y == f.Y {
		return true
	}
	return false
}

func (s *Snake) Grow() {
	s.growing = true
}

type Food struct {
	*Coord
}

func CreateFood() *Food {
	x := rand.Intn(WindowW/tileSize - 1)
	y := rand.Intn(WindowH/tileSize - 1)

	return &Food{&Coord{x, y}}
}

// Game implements ebiten.Game interface.
type Game struct {
	*Snake
	*Food
}

func (s *Snake) GoMove() {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {
			s.Move()
		}
	}()
	<-DeadChan
	ticker.Stop()
}

func (s *Snake) CheckDead() bool {
	// Check wall
	if s.body[0].X >= WindowW/tileSize ||
		s.body[0].X < 0 {
		return true
	}
	if s.body[0].Y > WindowH/tileSize ||
		s.body[0].Y < 0 {
		return true
	}

	// Check body
	for i := 0; i < len(s.body); i++ {
		for j := i + 1; j < len(s.body); j++ {
			if s.body[i].X == s.body[j].X &&
				s.body[i].Y == s.body[j].Y {
				fmt.Println("DeadBody")
				return true
			}
		}
	}
	return false
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	if g.Snake == nil {
		g.Snake = createSnake()
		go g.Snake.GoMove()
	}
	if g.Food == nil {
		g.Food = CreateFood()
	}

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		if g.Snake.Direction != D {
			g.Snake.Direction = U
		}
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		if g.Snake.Direction != L {
			g.Snake.Direction = R
		}
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		if g.Snake.Direction != U {
			g.Snake.Direction = D
		}
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		if g.Snake.Direction != R {
			g.Snake.Direction = L
		}
	}

	if g.Snake.CheckFood(g.Food) {
		g.Food = nil
		g.Snake.Grow()
	}
	if g.Snake.CheckDead() {
		g.Snake = nil
		DeadChan <- true
	}

	// select {
	// case <-DeadChan:
	// g.Snake = createSnake()
	// default:
	// fmt.Println("keep going")
	// }
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	screen.Fill(color.RGBA{0x77, 0x77, 0x77, 0xff})
	// if g.boardImage == nil {
	// w, h := 480, 480
	// g.boardImage = ebiten.NewImage(w, h)
	// g.boardImage.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
	// }
	//
	// op := &ebiten.DrawImageOptions{}
	// sw, sh := screen.Size()
	// bw, bh := g.boardImage.Size()
	// x := (sw - bw) / 2
	// y := (sh - bh) / 2
	// op.GeoM.Translate(float64(x), float64(y))
	// screen.DrawImage(g.boardImage, op)
	if g.Snake != nil {
		g.Snake.Draw(screen)
	}
	if g.Food != nil {
		g.Food.Draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	fmt.Println("Hello world")
	DeadChan = make(chan interface{}, 1)
	game := &Game{}
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(WindowW, WindowH)
	ebiten.SetWindowTitle("snake-go")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
