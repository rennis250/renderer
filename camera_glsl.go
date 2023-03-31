package renderer

import "fmt"

type Camera struct {
	Origin, Target [3]float64
	WithAperture   bool
}

func (r *Renderer) AttachCamera(c Camera) {
	r.camera = c
}

func (c *Camera) toGLSL() string {
	orig := `vec3(` + fmt.Sprintf("%f", c.Origin[0]) + ", " + fmt.Sprintf("%f", c.Origin[1]) + ", " + fmt.Sprintf("%f", c.Origin[2]) + ");"
	targ := `vec3(` + fmt.Sprintf("%f", c.Target[0]) + ", " + fmt.Sprintf("%f", c.Target[1]) + ", " + fmt.Sprintf("%f", c.Target[2]) + ");"

	s := `
    Ray getPrimaryRay(vec2 uv) {
        Ray r;

        // r.origin = vec3(-0.1, 0.8, 2.1);
        r.origin = ` + orig + `
        // r.origin = vec3(0.0, 0.0, -2.0);
        // r.origin = vec3(0.0, -0.4, -2.0);

        vec2 p = ConcentricSampleDisk();
        r.origin.xy += p*0.004;

        vec3 vup = vec3(0.0, 1.0, 0.0);

        // vec3 ta = vec3(-0.5, 0.8, 2.5);
        vec3 ta = ` + targ + `
        // vec3 ta = vec3(0.0, 0.0, 0.0);
        vec3 cw = normalize(ta - r.origin);
        vec3 cu = normalize(cross(cw, vup));
        vec3 cv = normalize(cross(cu, cw));

        mat3 ca = mat3(cu, cv, cw);

        r.direction = ca * normalize(vec3(-1.0*uv.x, uv.y, 1.0));

        return r;
    }
    `

	return s
}
