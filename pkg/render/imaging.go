package render

import (
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
)

// separate file because  we want to import x/image/draw
// instead of image/draw.

func resize(i image.Image, r image.Rectangle) image.Image {
	dst := image.NewRGBA(r)
	s := draw.BiLinear
	s.Scale(dst, r, i, i.Bounds(), draw.Over, nil)
	return dst
}

// create a mask image by using the gray value of the given image as the
// value for the mask alpha channel. Returns the mask image.
func createMask(i image.Image) image.Image {
	rect := i.Bounds()
	mask := image.NewRGBA(rect)
	x0 := rect.Min.X
	x1 := rect.Max.X
	y0 := rect.Min.Y
	y1 := rect.Max.Y

	var gray uint8
	var r, g, b uint32
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			r, g, b, _ = i.At(x, y).RGBA()
			gray = uint8((r + g + b) / 3)
			mask.Set(x, y, color.RGBA{0, 0, 0, gray})
		}
	}

	return mask
}

// applyOpacity applies the given opacity (0.0..1.0) to the given image.
// This method returns a new image where the alpha channel is a combination
// of the source alpha and the opacity.
func applyOpacity(i image.Image, opacity float64) image.Image {
	alpha := uint8(math.Round(255 * opacity))
	mask := image.NewUniform(color.Alpha{alpha})

	rect := i.Bounds()
	dst := image.NewRGBA(rect)
	p := image.ZP
	draw.DrawMask(dst, rect, i, p, mask, p, draw.Over)
	return dst
}

// Rotate the given image counter-clockwise by angle (radians) degrees.
// Rotation is around the center of the source image.
// Returns an image with the rotated pixels.
func rotate(angle float64, i image.Image) image.Image {
	// Size of the source image
	box := i.Bounds()
	xMax := box.Max.X
	yMax := box.Max.Y

	// Create the destination image.
	// The dst size is the diagonal accross the source Rectangle.
	a := float64(box.Max.X - box.Min.X)
	b := float64(box.Max.Y - box.Min.Y)
	c := math.Sqrt(math.Pow(a, 2) + math.Pow(b, 2))
	size := int(math.Ceil(c))
	dst := image.NewRGBA(image.Rect(0, 0, size, size))

	// Rotation around center instead of origin
	// means: Translate - Rotate - Translate
	t0 := translation(-a/2, -b/2)
	rot := rotation(angle)
	t1 := translation(a/2, b/2)

	// Transform each pixel and set it on the destination image.
	var tx, ty float64
	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			tx, ty = float64(x), float64(y)
			tx, ty = transform(t0, tx, ty)
			tx, ty = transform(rot, tx, ty)
			tx, ty = transform(t1, tx, ty)

			tx = math.Round(tx)
			ty = math.Round(ty)
			dst.Set(int(tx), int(ty), i.At(x, y))
		}
	}

	return dst
}
