package main

import (
	"fmt"
	"image/color"
	"image/color/palette"
	_ "image/png"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 500
	screenHeight = 500
	maxIt        = 100
)

type Mandelbrot struct {
	duration     time.Duration
	offscreen    *ebiten.Image
	offscreenPix []byte
	palette      []color.RGBA
	zoom         float64
}

func NewMandelbrot() *Mandelbrot {
	result := &Mandelbrot{}
	result.offscreen, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	result.offscreenPix = make([]byte, screenWidth*screenHeight*4)
	result.palette = make([]color.RGBA, len(palette.Plan9))
	for i := range palette.Plan9 {
		r, g, b, a := palette.Plan9[i].RGBA()
		result.palette[i] = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	result.zoom = 2
	return result
}

func iter(cx, cy float64) int {
	var x, y, xx, yy float64 = 0.0, 0.0, 0.0, 0.0

	for i := 0; i < maxIt; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}
	return maxIt
}

func (m *Mandelbrot) updateOffscreen(centerX, centerY, size float64) {
	start := time.Now()
	var wg sync.WaitGroup
	for j := 0; j < screenHeight; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			for i := 0; i < screenHeight; i++ {
				it := iter(float64(i)*size/screenWidth-size/2+centerX, (screenHeight-float64(j))*size/screenHeight-size/2+centerY)
				p := 4 * (i + j*screenWidth)
				if it < maxIt {
					rgba := m.palette[it]
					m.offscreenPix[p] = rgba.R
					m.offscreenPix[p+1] = rgba.G
					m.offscreenPix[p+2] = rgba.B
					m.offscreenPix[p+3] = rgba.A
				} else {
					m.offscreenPix[p] = 0
					m.offscreenPix[p+1] = 0
					m.offscreenPix[p+2] = 0
					m.offscreenPix[p+3] = 0
				}
			}
		}(j)
	}
	wg.Wait()
	m.offscreen.ReplacePixels(m.offscreenPix)
	m.duration = time.Since(start)
}

func (m *Mandelbrot) Update(screen *ebiten.Image) error {
	m.updateOffscreen(-0.70, 0.25, m.zoom)
	m.zoom *= 0.995
	return nil
}

func (m *Mandelbrot) Draw(screen *ebiten.Image) {
	screen.DrawImage(m.offscreen, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d ms %.2f fps", m.duration/time.Millisecond, ebiten.CurrentFPS()))
}

func (m *Mandelbrot) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mandelbrot")
	if err := ebiten.RunGame(NewMandelbrot()); err != nil {
		log.Fatal(err)
	}
}
