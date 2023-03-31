package renderer

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func (r *Renderer) screenshot(flip bool) image.Image {
	w, h := int32(r.WindowWidth), int32(r.WindowHeight)
	screenshot := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	stride := int32(screenshot.Stride)
	pixels := make([]byte, len(screenshot.Pix))
	gl.ReadPixels(0, 0, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	// OpenGL reads pixels from the lower-left. Let's fix that.
	if flip {
		for y := int32(0); y < h; y++ {
			i := (h - 1 - y) * stride
			copy(screenshot.Pix[y*stride:], pixels[i:i+w*4])
		}
	} else {
		for y := int32(0); y < h; y++ {
			i := y * stride
			copy(screenshot.Pix[y*stride:], pixels[i:i+w*4])
		}
	}

	return screenshot
}
