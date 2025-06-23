package main

import (
	"image"
	"image/png"
	"os"

	"deedles.dev/trunic"
)

func main() {
	var r trunic.Renderer
	r.AppendRune("t", "e")
	r.AppendRune("s")
	r.AppendRune("t")

	img := image.NewRGBA(r.Bounds().Inset(-20))
	r.DrawTo(img, 0, 0)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}
