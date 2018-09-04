package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel/imdraw"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hero Gopher",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	heroImg := loadHero("gotham-3x.png")

	hero := pixel.NewSprite(heroImg, heroImg.Bounds())

	h := Hero{
		Sprite:   hero,
		Position: pixel.V(0, 0),
	}
	/*
		phys := &gophy{
			gravity:  -512,
			flySpeed: 64,
			lift: 256,
			rec:      pixel.R(0, 0, 100, 100),
		}
	*/
	camPos := pixel.ZV

	last := time.Now()
	canvas := pixelgl.NewCanvas(pixel.R(-512, -384, 512, 370))

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(0, 374), basicAtlas)
	imd := imdraw.New(hero.Picture())

	for !win.Closed() {
		basicTxt.Clear()
		fmt.Fprintf(basicTxt, "Score: %v", h.Position.X)

		dt := time.Since(last).Seconds()
		last = time.Now()

		canvas.Clear(colornames.Darkgreen)
		win.Clear(colornames.Black)

		camPos = pixel.Lerp(camPos, pixel.V(h.Position.X, 0), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos.Scaled(-1))
		canvas.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyUp) {
			h.Position.Y = h.Position.Y + 15
		} else {
			h.Position.Y = h.Position.Y - 10
		}
		h.Position.X++
		basicTxt.Draw(win, pixel.IM)
		imd.Draw(canvas)

		h.Sprite.Draw(canvas, pixel.IM.Moved(h.Position))
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		win.Update()
	}

}

type Hero struct {
	Sprite   *pixel.Sprite
	Position pixel.Vec
}

func loadHero(imgPath string) pixel.Picture {
	file, err := os.Open(imgPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
	}

	return pixel.PictureDataFromImage(img)

}
func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
