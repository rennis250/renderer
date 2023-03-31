package renderer

import (
	"encoding/json"
	"image"
	"os"
	"time"

	"github.com/cstegel/opengl-samples-golang/colors/gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {
	WindowWidth, WindowHeight               float64
	pathProgram, filterProgram, fxaaProgram *gfx.Program
	VAO                                     uint32
	startTime                               time.Time
	frame                                   float64
	bounces, samples                        int32
	bluentex, multispectex                  *gfx.Texture
	pttex, filttex                          *RobTexture
	scene                                   string
	camera                                  Camera
}

func NewRenderer(windowWidth, windowHeight float64, cam Camera) (*Renderer, error) {
	var r Renderer

	r.WindowWidth = windowWidth
	r.WindowHeight = windowHeight
	r.startTime = time.Now()
	r.frame = 0
	r.scene = ""
	r.bounces = 10
	r.samples = 300
	r.camera = cam

	return &r, nil
}

func (r *Renderer) SetBounces(b int32) {
	r.bounces = b
}

func (r *Renderer) SetSamples(s int32) {
	r.samples = s
}

func NewRendererFromJSON(jsonFile string) (*Renderer, error) {
	dat, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	res := Scene{}
	if err := json.Unmarshal(dat, &res); err != nil {
		panic(err)
	}

	r, err := NewRenderer(res.Width, res.Height, res.Camera)
	if err != nil {
		return nil, err
	}

	r.SetBounces(res.Bounces)
	r.SetSamples(res.Samples)

	for _, v := range res.Shapes {
		for k, v2 := range v.(map[string]interface{}) {
			r.AddShapeFromJSON(k, v2.(map[string]interface{}))
		}
	}

	return r, nil
}

func NewRendererFromBytes(bytes []byte) (*Renderer, error) {
	res := Scene{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		panic(err)
	}

	r, err := NewRenderer(res.Width, res.Height, res.Camera)
	if err != nil {
		return nil, err
	}

	r.SetBounces(res.Bounces)
	r.SetSamples(res.Samples)

	for _, v := range res.Shapes {
		for k, v2 := range v.(map[string]interface{}) {
			r.AddShapeFromJSON(k, v2.(map[string]interface{}))
		}
	}

	return r, nil
}

func (r *Renderer) Render() (image.Image, error) {
	// patht
	r.pathBlit()

	ss := r.screenshot(false)
	tex, err := RobNewTexture(ss, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	r.pttex = tex
	if err != nil {
		return nil, err
	}

	// mitchell
	r.mitchellBlit()

	ss = r.screenshot(false)
	tex, err = RobNewTexture(ss, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	r.filttex = tex
	if err != nil {
		return nil, err
	}

	// fxaa
	r.fxaaBlit()

	return r.screenshot(true), nil
}

func (r *Renderer) Close() {
	r.pathProgram.Delete()
	r.filterProgram.Delete()
	r.fxaaProgram.Delete()

	r.bluentex.UnBind()
	r.multispectex.UnBind()

	r.pttex.UnBind()
	r.filttex.UnBind()
}
