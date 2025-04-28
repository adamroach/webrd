package imageconvert

import (
	"errors"
	"image"
	"image/color"
)

func ToYCbCr(dst *image.YCbCr, src image.Image) error {
	if dst.Bounds() != src.Bounds() {
		return errors.New("images must be the same size")
	}
	if dst.SubsampleRatio != image.YCbCrSubsampleRatio420 {
		return errors.New("only 4:2:0 subsampling is currently supported")
	}

	// Currently, we don't handle the final pixel for images that are an odd number
	// of pixels wide and/or high
	width := src.Bounds().Dx() & (^1)
	height := src.Bounds().Dy() & (^1)
	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x += 2 {
			r, g, b, _ := src.At(x, y).RGBA()
			yy, cb1, cr1 := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			dst.Y[x+y*dst.YStride] = yy

			r, g, b, _ = src.At(x+1, y).RGBA()
			yy, cb2, cr2 := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			dst.Y[(x+1)+y*dst.YStride] = yy

			r, g, b, _ = src.At(x, y+1).RGBA()
			yy, cb3, cr3 := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			dst.Y[x+(y+1)*dst.YStride] = yy

			r, g, b, _ = src.At(x+1, y+1).RGBA()
			yy, cb4, cr4 := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			dst.Y[(x+1)+(y+1)*dst.YStride] = yy

			cb := uint8((int(cb1) + int(cb2) + int(cb3) + int(cb4)) / 4)
			cr := uint8((int(cr1) + int(cr2) + int(cr3) + int(cr4)) / 4)
			dst.Cb[x/2+(y/2)*dst.CStride] = cb
			dst.Cr[x/2+(y/2)*dst.CStride] = cr
		}
	}
	return nil
}
