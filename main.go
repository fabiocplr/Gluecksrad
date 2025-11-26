package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var segmentTexts = []string{
	"Socken",
	"Sauna",
	"Suppe",
	"Saft",
	"Spiegelei",
	"Spiele Abend",
	"Süßigkeiten",
	"Saufen",
	"Schnitzel",
	"Staubsaugen",
	"Streaming",
}

var segmentCount = len(segmentTexts)

type Game struct {
	// Kreis Variabeln
	angle    float64
	speed    float64
	spinning bool

	// Segment Ergebnis
	result    int
	hasResult bool

	// Tutorial Game Text
	tutorialText string
}

func (g *Game) Update() error {
	// Starten mit SPACE
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.angle = 0
		g.speed = 0.2 + rand.Float64()*1.0
		g.spinning = true
		g.hasResult = false
		g.tutorialText = ""
	}

	// Drehen + Bremsen
	if g.spinning {
		g.angle += g.speed
		g.speed *= 0.98
		if g.speed < 0.001 {
			g.speed = 0
			g.spinning = false

			segmentAngle := 2 * math.Pi / float64(segmentCount)
			// Winkel normalisieren
			a := math.Mod(g.angle, 2*math.Pi)
			if a < 0 {
				a += 2 * math.Pi
			}

			// Segment bestimmen
			g.result = int(a / segmentAngle)

			g.hasResult = true

		}
	}

	return nil
}

func drawCircle(screen *ebiten.Image, cx, cy, r float64, col color.Color) {
	const steps = 64
	for i := 0; i < steps; i++ {
		a1 := 2 * math.Pi * float64(i) / float64(steps)
		a2 := 2 * math.Pi * float64(i+1) / float64(steps)

		x1 := cx + r*math.Cos(a1)
		y1 := cy + r*math.Sin(a1)

		x2 := cx + r*math.Cos(a2)
		y2 := cy + r*math.Sin(a2)

		ebitenutil.DrawLine(screen, x1, y1, x2, y2, col)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{30, 30, 30, 255})

	cx := 160.0
	cy := 120.0
	radius := 80.0

	// Kreis zeichnen
	drawCircle(screen, cx, cy, radius, color.RGBA{139, 90, 43, 255})

	// Segmente
	for i := 0; i < segmentCount; i++ {
		angle := 2 * math.Pi * float64(i) / float64(segmentCount)

		sx := cx + radius*math.Cos(angle)
		sy := cy + radius*math.Sin(angle)
		ebitenutil.DrawLine(screen, cx, cy, sx, sy, color.RGBA{139, 90, 43, 255})
	}

	// Zeiger
	zx := cx + radius*math.Cos(g.angle)
	zy := cy + radius*math.Sin(g.angle)
	ebitenutil.DrawLine(screen, cx, cy, zx, zy, color.RGBA{245, 222, 156, 255})

	// Winkel in Grad
	deg := math.Mod((g.angle+math.Pi/2)*180/math.Pi, 360)
	if deg < 0 {
		deg += 360
	}

	if g.tutorialText != "" {
		tutorialX := int(cx) - len(g.tutorialText)*4
		tutorialY := 5
		ebitenutil.DebugPrintAt(screen, g.tutorialText, tutorialX, tutorialY)
	}

	// Daten
	dataAngle := fmt.Sprintf("Winkel: %.0f°", deg)
	ebitenutil.DebugPrint(screen, dataAngle)

	// Gewinner-Text unter dem Rad
	winnerText := ""
	if g.hasResult && g.result >= 0 && g.result < len(segmentTexts) {
		winnerText = fmt.Sprintf("%s", segmentTexts[g.result])
	}
	// Position: unter dem Rad (cx, cy+radius+…)
	textX := int(cx) - 60
	textY := int(cy+radius) + 20
	ebitenutil.DebugPrintAt(screen, winnerText, textX, textY)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {

	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Glücksrad Edith")
	if err := ebiten.RunGame(&Game{
		tutorialText: "SPACE drücken!",
	}); err != nil {
		log.Fatal(err)
	}
}
