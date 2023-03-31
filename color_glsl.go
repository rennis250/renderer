package renderer

import (
	"fmt"
	"math"
)

type Color interface {
	colorToGLSLString() string
}

type ColorName int

const (
	RED ColorName = iota
	GREEN
	BLUE
	YELLOW
	REDSUR
	GREENSUR
	BLUESUR
	YELLOWSUR
	REDT
	GREENT
	BLUET
	YELLOWT
	WHITE
	UNITY
	HALF
	NONE
)

func (spec ColorName) colorToGLSLString() string {
	switch spec {
	case RED:
		return "RED"
	case GREEN:
		return "GREEN"
	case BLUE:
		return "BLUE"
	case YELLOW:
		return "YELLOW"
	case REDSUR:
		return "REDSUR"
	case GREENSUR:
		return "GREENSUR"
	case BLUESUR:
		return "BLUESUR"
	case YELLOWSUR:
		return "YELLOWSUR"
	case REDT:
		return "REDT"
	case GREENT:
		return "GREENT"
	case BLUET:
		return "BLUET"
	case YELLOWT:
		return "YELLOWT"
	case WHITE:
		return "WHITE"
	case UNITY:
		return "UNITY"
	case HALF:
		return "HALF"
	case NONE:
		return "NONE"
	default:
		return "NONE"
	}
}

func colorFromInterface(in interface{}) Color {
	switch v := in.(type) {
	case string:
		switch v {
		case "RED":
			return RED
		case "GREEN":
			return GREEN
		case "BLUE":
			return BLUE
		case "YELLOW":
			return YELLOW
		case "REDSUR":
			return REDSUR
		case "GREENSUR":
			return GREENSUR
		case "BLUESUR":
			return BLUESUR
		case "YELLOWSUR":
			return YELLOWSUR
		case "REDT":
			return REDT
		case "GREENT":
			return GREENT
		case "BLUET":
			return BLUET
		case "YELLOWT":
			return YELLOWT
		case "WHITE":
			return WHITE
		case "UNITY":
			return UNITY
		case "HALF":
			return HALF
		case "NONE":
			return NONE
		default:
			return NONE
		}

	default:
		vms := v.(map[string]interface{})
		dkl := ColorDKL{LD: vms["LD"].(float64), RG: vms["RG"].(float64), YV: vms["YV"].(float64)}
		return dkl
	}
}

type ColorDKL struct {
	LD, RG, YV float64
}

func (ldrgyv ColorDKL) colorToGLSLString() string {
	GREENSUR := [12]float64{0.092, 0.096683, 0.10227, 0.13182, 0.40227, 0.44626, 0.31621, 0.19174, 0.12782, 0.11574, 0.13324, 0.15868}
	REDSUR := [12]float64{0.04, 0.054636, 0.060382, 0.059549, 0.055277, 0.058049, 0.067844, 0.1817, 0.50124, 0.63632, 0.62978, 0.64319}
	BLUESUR := [12]float64{0.70507, 0.86353, 1, 0.92747, 0.66104, 0.39434, 0.24311, 0.19277, 0.18943, 0.19758, 0.20284, 0.20418}
	YELLOWSUR := [12]float64{0.091609, 0.093483, 0.097796, 0.14636, 0.31904, 0.60099, 0.82454, 0.92035, 0.94373, 0.95753, 0.97963, 1}

	rg_mix := 0.5*ldrgyv.RG + 0.5
	yv_mix := 0.5*ldrgyv.YV + 0.5

	c := 1
	match_col_s := "Spectrum(vec4("
	for x, _ := range GREENSUR {
		rg_s := (1.0-rg_mix)*REDSUR[x] + rg_mix*GREENSUR[x]
		yv_s := (1.0-yv_mix)*BLUESUR[x] + yv_mix*YELLOWSUR[x]

		rg_yv_s := 0.5*rg_s + 0.5*yv_s

		match_val := rg_yv_s * ldrgyv.LD

		if c == 12 {
			match_col_s += fmt.Sprintf("%f", match_val) + "))"
		} else if math.Mod(float64(c), 4.0) == 0.0 {
			match_col_s += fmt.Sprintf("%f", match_val) + "), vec4("
		} else {
			match_col_s += fmt.Sprintf("%f", match_val) + ", "
		}

		c += 1
	}

	return match_col_s
}

func ColorfromJSON(c interface{}) Color {
	switch v := c.(type) {
	case map[string]interface{}:
		return ColorDKL{
			LD: v["LD"].(float64),
			RG: v["RG"].(float64),
			YV: v["YV"].(float64),
		}
	default:
		return c.(ColorName)
	}
}
