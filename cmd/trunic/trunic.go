package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"deedles.dev/trunic"
)

type imageLines []image.Image

func (img imageLines) ColorModel() color.Model {
	return img[0].ColorModel()
}

func (img imageLines) Bounds() (r image.Rectangle) {
	for _, s := range img {
		bounds := s.Bounds().Canon()
		bounds = image.Rect(bounds.Min.X, r.Max.Y, bounds.Max.X, r.Max.Y+bounds.Dy())
		r = r.Union(bounds)
	}
	return r
}

func (img imageLines) At(x, y int) color.Color {
	p := image.Pt(x, y)

	var prev image.Rectangle
	for _, s := range img {
		ibounds := s.Bounds().Canon()
		bounds := image.Rect(ibounds.Min.X, prev.Max.Y, ibounds.Max.X, prev.Max.Y+ibounds.Dy())
		if p.In(bounds) {
			p = p.Add(ibounds.Min.Sub(bounds.Min))
			return s.At(p.X, p.Y)
		}
		prev = bounds
	}
	return img.ColorModel().Convert(color.White)
}

func writeImage(output string, img image.Image) error {
	file := io.Writer(os.Stdout)
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return err
		}
		defer f.Close()
		file = f
	}

	err := png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("encode image: %w", err)
	}

	return nil
}

func run(ctx context.Context) error {
	output := flag.String("o", "", "output filename (empty for stdout)")
	flag.Parse()

	var lines imageLines

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var r trunic.Renderer
		r.Append(s.Text())
		img := image.NewRGBA(r.Bounds().Inset(-20))
		draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)
		r.DrawTo(img, 0, 0)

		lines = append(lines, img)
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("read line from input: %w", err)
	}
	if len(lines) == 0 {
		return nil
	}

	err := writeImage(*output, lines)
	if err != nil {
		return fmt.Errorf("write image: %w", err)
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := run(ctx)
	if err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}
