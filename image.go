package trunic

import (
	"image"
	"image/color"
	"math"
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

type shifted struct {
	image.Image
	offset image.Point
}

func (img *shifted) Bounds() image.Rectangle {
	return img.Image.Bounds().Add(img.offset)
}

func (img *shifted) At(x, y int) color.Color {
	return img.Image.At(x-img.offset.X, y-img.offset.Y)
}

type stack struct {
	Color color.Color

	images []image.Image
	bounds image.Rectangle
}

func stackOf(images ...image.Image) *stack {
	var s stack
	for _, img := range images {
		s.Append(img)
	}
	return &s
}

func (s *stack) color() color.NRGBA64 {
	if s.Color == nil {
		return color.NRGBA64{0, 0, 0, math.MaxUint16}
	}
	return color.NRGBA64Model.Convert(s.Color).(color.NRGBA64)
}

func (s *stack) Append(img image.Image) {
	s.images = append(s.images, img)
	s.bounds = s.bounds.Union(img.Bounds())
}

func (s *stack) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (s *stack) Bounds() image.Rectangle {
	return s.bounds
}

func (s *stack) At(x, y int) color.Color {
	var a uint32
	for _, img := range s.images {
		_, _, _, ia := img.At(x, y).RGBA()
		a = max(ia, a)
	}
	return s.fade(uint16(a * math.MaxUint16 / 0xFFFF))
}

func (s *stack) fade(to uint16) color.Color {
	c := s.color()
	c.A = to
	return c
}

type row struct {
	images []image.Image
	bounds image.Rectangle
}

func rowOf(images ...image.Image) *row {
	var r row
	for _, img := range images {
		r.Append(img)
	}
	return &r
}

func (r *row) Append(img image.Image) {
	shift := &shifted{
		Image:  img,
		offset: image.Pt(r.bounds.Max.X, 0),
	}

	r.images = append(r.images, shift)
	r.bounds = r.bounds.Union(shift.Bounds())
}

func (r *row) ColorModel() color.Model {
	return r.images[0].ColorModel()
}

func (r *row) Bounds() image.Rectangle {
	return r.bounds
}

func (r *row) At(x, y int) color.Color {
	return r.imageAt(x, y).At(x, y)
}

func (r *row) imageAt(x, y int) image.Image {
	i := x / letterWidth
	if (i < 0) || (i >= len(r.images)) {
		return image.Transparent
	}
	return r.images[i]
}
