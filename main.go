package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	angle    float64
	speed    float64
	spinning bool
}

func (g *Game) Update() error {
	// Starten mit SPACE
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.angle = 0.0
		g.speed = 0.3 // Startwert, nicht +=
		g.spinning = true
	}

	// Drehen + Bremsen
	if g.spinning {
		g.angle += g.speed
		g.speed *= 0.98
		if g.speed < 0.001 {
			g.speed = 0
			g.spinning = false
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	data := fmt.Sprintf("Angle: %.2f (SPACE drücken)", g.angle)
	ebitenutil.DebugPrint(screen, data)

	screen.Fill(color.RGBA{30, 30, 30, 255})

	cx := 160.0
	cy := 120.0
	radius := 80.0

	x := cx + radius*math.Cos(g.angle)
	y := cy + radius*math.Sin(g.angle)

	zeiger := ebitenutil.DrawLine(screen, cx, cy, x, y, white)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
