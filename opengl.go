package renderer

import (
	"errors"
	"image"
	"image/draw"

	"github.com/cstegel/opengl-samples-golang/colors/gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Creates the Vertex Array Object for a triangle.
func (r *Renderer) createVAO(vertices []float32, indices []uint32) {
	// var VAO uint32
	gl.GenVertexArrays(1, &r.VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(r.VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 2*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// texture position
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)
}

func (r *Renderer) loadTextures() error {
	tex, err := gfx.NewTextureFromFile("/home/me/Documents/source_code/vrRender/renderer/images/LDR_RGBA_13.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	r.bluentex = tex
	if err != nil {
		return err
	}

	tex, err = gfx.NewTextureFromFile("/home/me/Documents/source_code/vrRender/renderer/images/multispectral_tex.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	r.multispectex = tex
	if err != nil {
		return err
	}

	return nil
}

type RobTexture struct {
	handle  uint32
	target  uint32 // same target as gl.BindTexture(<this param>, ...)
	texUnit uint32 // Texture unit that is currently bound to ex: gl.TEXTURE0
}

var errUnsupportedStride = errors.New("unsupported stride, only 32-bit colors supported")
var errTextureNotBound = errors.New("texture not bound")

func RobNewTexture(img image.Image, wrapR, wrapS int32) (*RobTexture, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return nil, errUnsupportedStride
	}

	var handle uint32
	gl.GenTextures(1, &handle)

	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.RGBA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := RobTexture{
		handle: handle,
		target: target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	// set the texture wrapping/filtering options (applies to current bound texture obj)
	// TODO-cs
	gl.TexParameteri(target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)

	gl.GenerateMipmap(handle)

	return &texture, nil
}

func (tex *RobTexture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.target, tex.handle)
	tex.texUnit = texUnit
}

func (tex *RobTexture) UnBind() {
	tex.texUnit = 0
	gl.BindTexture(tex.target, 0)
}

func (tex *RobTexture) SetUniform(uniformLoc int32) error {
	if tex.texUnit == 0 {
		return errTextureNotBound
	}
	gl.Uniform1i(uniformLoc, int32(tex.texUnit-gl.TEXTURE0))
	return nil
}
