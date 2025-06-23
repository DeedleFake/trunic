package main

import (
	"image"
	"image/png"
	"os"

	"deedles.dev/trunic"
)

func main() {
	r := trunic.Renderer{
		Kerning: -2,
	}
	r.AppendRune("t", "e")
	r.AppendRune("s")
	r.AppendRune("t")

	img := image.NewRGBA(r.Bounds())
	r.Draw(img, img.Bounds(), image.Black, image.Point{})
	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}
