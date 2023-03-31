package renderer

import "fmt"

func (r *Renderer) generateGLSLConstants() string {
	s := `
#define Pi 3.14159265358979323846
#define InvPi 0.31830988618379067154
#define PiOver2 1.57079632679489661923
#define PiOver4 0.78539816339744830961

#define MAX_DIST 1e10
#define MAX_SCENE_DIST 100.0

#define EPS 0.0001
#define OneMinusEpsilon 0.99999994

#define NBOUNCES ` + fmt.Sprintf("%d", r.bounces) + ` // 10
#define NSAMPLES ` + fmt.Sprintf("%d", r.samples) + ` // 300

#define ZERO min(int(iFrame), 0)
`

	return s
}
