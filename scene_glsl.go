package renderer

type Scene struct {
	Width, Height    float64
	Bounces, Samples int32
	Camera           Camera
	Shapes           []interface{}
}

const scene_header_glsl = `
Surfel findFirstIntersection(const vec3 X, const vec3 wi) {
    Surfel s;
    s.hit = false;
    s.t = MAX_DIST;

    float t;
    vec3 nor;
`

const scene_tail_glsl = `   
    return s;
}`

func (r *Renderer) AssembleScene() error {
	r.createVAO(cubeVertices, nil)

	err := r.loadTextures()
	if err != nil {
		return err
	}

	err = r.buildShaders(scene_header_glsl + r.scene + scene_tail_glsl)
	if err != nil {
		return err
	}

	return nil
}

func (r *Renderer) ClearScene() {
	r.scene = ""
}

func (r *Renderer) AddShape(s Shape) {
	r.scene += s.ShapeToGLSLString()
}

func (r *Renderer) AddShapeFromJSON(shapeType string, shape map[string]interface{}) {
	mat := materialFromString(shape["Material"].(string))
	col := colorFromInterface(shape["Color"])

	switch shapeType {
	case "Box":
		pos := shape["Position"].([]interface{})
		size := shape["Size"].([]interface{})

		r.AddShape(&Box{
			Position: [3]float64{pos[0].(float64), pos[1].(float64), pos[2].(float64)},
			Size:     [3]float64{size[0].(float64), size[1].(float64), size[2].(float64)},
			Material: mat,
			Color:    col,
		})

	case "Plane":
		ori := shape["Orientation"].([]interface{})

		r.AddShape(&Plane{
			Orientation: [3]float64{ori[0].(float64), ori[1].(float64), ori[2].(float64)},
			Offset:      shape["Offset"].(float64),
			Material:    mat,
			Color:       col,
		})

	case "Sphere":
		pos := shape["Position"].([]interface{})

		r.AddShape(&Sphere{
			Position: [3]float64{pos[0].(float64), pos[1].(float64), pos[2].(float64)},
			Radius:   shape["Radius"].(float64),
			Material: mat,
			Color:    col,
		})

	case "Triangle":
		cent := shape["Center"].([]interface{})
		p1 := shape["Point1"].([]interface{})
		p2 := shape["Point2"].([]interface{})
		p3 := shape["Point3"].([]interface{})

		r.AddShape(&Triangle{
			Center:   [3]float64{cent[0].(float64), cent[1].(float64), cent[2].(float64)},
			Point1:   [3]float64{p1[0].(float64), p1[1].(float64), p1[2].(float64)},
			Point2:   [3]float64{p2[0].(float64), p2[1].(float64), p2[2].(float64)},
			Point3:   [3]float64{p3[0].(float64), p3[1].(float64), p3[2].(float64)},
			Material: mat,
			Color:    col,
		})

	case "Cone":
		pos := shape["Position"].([]interface{})
		pa := shape["PointA"].([]interface{})
		pb := shape["PointB"].([]interface{})

		r.AddShape(&Cone{
			Position: [3]float64{pos[0].(float64), pos[1].(float64), pos[2].(float64)},
			PointA:   [3]float64{pa[0].(float64), pa[1].(float64), pa[2].(float64)},
			PointB:   [3]float64{pb[0].(float64), pb[1].(float64), pb[2].(float64)},
			RadiusA:  shape["RadiusA"].(float64),
			RadiusB:  shape["RadiusB"].(float64),
			Material: mat,
			Color:    col,
		})

	default:
		panic("The requested shape is not available. Please check for typos in scene description")
	}
}
