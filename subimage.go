package trunic

import (
	"image"
	"image/color"
)

type subimage struct {
	image.Image
	rect image.Rectangle
}

func (img *subimage) Bounds() image.Rectangle {
	return img.rect.Sub(img.rect.Min)
}

func (img *subimage) At(x, y int) color.Color {
	return img.Image.At(x+img.rect.Min.X, y+img.rect.Min.Y)
}
