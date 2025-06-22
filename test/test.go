package main

import (
	"image/png"
	"os"

	"deedles.dev/trunic"
)

func main() {
	img := trunic.Render("test")
	err := png.Encode(os.Stdout, img)
	if err != nil {
		panic(err)
	}
}
