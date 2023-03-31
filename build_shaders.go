package renderer

import (
	"github.com/cstegel/opengl-samples-golang/colors/gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func (r *Renderer) buildShaders(scene_glsl string) error {
	vert_shader, err := gfx.NewShader(vert_glsl, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	frag_src := header_uniforms_glsl + r.generateGLSLConstants() + rand_glsl + ray_algebra_glsl + shape_intersect_glsl + spectra_glsl + texture_glsl + material_glsl + surfel_glsl + sky_glsl + scene_glsl + pathtrace_glsl + r.camera.toGLSL() + tonemap_glsl + main_glsl
	frag_shader, err := gfx.NewShader(frag_src, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	prog, err := gfx.NewProgram(vert_shader, frag_shader)
	r.pathProgram = prog
	if err != nil {
		return err
	}

	frag_src = filt_header_glsl + filt_blur_glsl + filt_main_glsl
	frag_shader, err = gfx.NewShader(frag_src, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	prog, err = gfx.NewProgram(vert_shader, frag_shader)
	r.filterProgram = prog
	if err != nil {
		return err
	}

	frag_src = fxaa_header_glsl + fxaa_apply_glsl + fxaa_main_glsl
	frag_shader, err = gfx.NewShader(frag_src, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	prog, err = gfx.NewProgram(vert_shader, frag_shader)
	r.fxaaProgram = prog
	if err != nil {
		return err
	}

	return nil
}
