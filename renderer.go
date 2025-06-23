package trunic

import (
	"image"
	"image/draw"
)

type Renderer struct {
	Kerning int

	ph [][]image.Image
}

func (r *Renderer) AppendRune(ph ...string) {
	imgs := make([]image.Image, 0, len(ph))
	for _, ph := range ph {
		imgs = append(imgs, fontMap[ph])
	}

	r.ph = append(r.ph, imgs)
}

func (r *Renderer) Draw(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point) {
	dr = dr.Canon()

	for i, ph := range r.ph {
		dp := dr.Min.Add(image.Pt(i*(letterWidth+r.Kerning), 0))
		dr := image.Rectangle{Min: dp, Max: dp.Add(image.Pt(letterWidth, letterHeight))}.Intersect(dr)

		for _, img := range ph {
			draw.DrawMask(
				dst,
				dr,
				src,
				sp,
				img,
				image.Point{},
				draw.Over,
			)
		}
	}
}

func (r *Renderer) Bounds() image.Rectangle {
	return image.Rectangle{Max: r.Size()}
}

func (r *Renderer) Size() image.Point {
	return image.Pt(
		len(r.ph)*(letterWidth+r.Kerning),
		letterHeight,
	)
}
