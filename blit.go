package renderer

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func (r *Renderer) pathBlit() {
	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.pathProgram.Use()

	gl.BindVertexArray(r.VAO)

	if r.camera.WithAperture {
		gl.Uniform1i(r.pathProgram.GetUniformLocation("with_aperture"), int32(1))
	} else {
		gl.Uniform1i(r.pathProgram.GetUniformLocation("with_aperture"), int32(0))
	}

	gl.Uniform2f(r.pathProgram.GetUniformLocation("iResolution"), float32(r.WindowWidth), float32(r.WindowHeight))
	gl.Uniform1f(r.pathProgram.GetUniformLocation("iTime"), float32(time.Since(r.startTime).Seconds()))
	gl.Uniform1f(r.pathProgram.GetUniformLocation("iFrame"), float32(r.frame))

	gl.Uniform1i(r.pathProgram.GetUniformLocation("sphere_color"), int32(1))

	r.bluentex.Bind(gl.TEXTURE0)
	r.bluentex.SetUniform(r.pathProgram.GetUniformLocation("blue_noise_tex"))

	r.multispectex.Bind(gl.TEXTURE1)
	r.multispectex.SetUniform(r.pathProgram.GetUniformLocation("voronoi_indexed_tex"))

	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	r.bluentex.UnBind()
	r.multispectex.UnBind()
}

func (r *Renderer) mitchellBlit() {
	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.filterProgram.Use()

	gl.BindVertexArray(r.VAO)
	gl.Uniform3f(r.filterProgram.GetUniformLocation("iResolution"), float32(r.WindowWidth), float32(r.WindowHeight), 0)

	r.pttex.Bind(gl.TEXTURE0)
	r.pttex.SetUniform(r.filterProgram.GetUniformLocation("tex"))

	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	r.pttex.UnBind()
}

func (r *Renderer) fxaaBlit() {
	gl.ClearColor(0, 0, 0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.fxaaProgram.Use()

	gl.BindVertexArray(r.VAO)
	gl.Uniform3f(r.fxaaProgram.GetUniformLocation("iResolution"), float32(r.WindowWidth), float32(r.WindowHeight), 0)

	r.filttex.Bind(gl.TEXTURE0)
	r.filttex.SetUniform(r.fxaaProgram.GetUniformLocation("tex"))

	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	r.filttex.UnBind()
}
