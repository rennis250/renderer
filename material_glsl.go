package renderer

type Material int

const (
	SKY Material = iota
	LAMB
	FLOOR
	WALL
	TEXTURE
	DOTS
	GLASS
	METAL
)

func materialString(mat Material) string {
	switch mat {
	case SKY:
		return "SKY"
	case LAMB:
		return "LAMB"
	case FLOOR:
		return "FLOOR"
	case WALL:
		return "WALL"
	case TEXTURE:
		return "TEXTURE"
	case DOTS:
		return "DOTS"
	case GLASS:
		return "GLASS"
	case METAL:
		return "METAL"
	default:
		return "LAMB"
	}
}

func materialFromString(mat string) Material {
	switch mat {
	case "SKY":
		return SKY
	case "LAMB":
		return LAMB
	case "FLOOR":
		return FLOOR
	case "WALL":
		return WALL
	case "TEXTURE":
		return TEXTURE
	case "DOTS":
		return DOTS
	case "GLASS":
		return GLASS
	case "METAL":
		return METAL
	default:
		return LAMB
	}
}

const material_glsl = `
#define SKY 1.0
#define LAMB 2.0
#define FLOOR 3.0
#define WALL 4.0
#define TEXTURE 5.0
#define DOTS 6.0
#define GLASS 7.0
#define METAL 8.0

struct Material {
    float type;
    Spectrum albedo;
    Spectrum emitted;
    bool specularBounce;
};

Material getMaterial(float which_mat, vec3 X, vec3 n) {
    Material mat;

    if (which_mat == LAMB) {
        mat.type = LAMB;
        mat.albedo = WHITE;
        mat.emitted = ZEROSPECT;
        mat.specularBounce = false;
    } else if (which_mat == FLOOR) {
        // mat.type = LAMB;
        // mat.albedo = spect_smoothstep(BLUE, YELLOW, disk(mod(fract(0.7 * X.xz), vec2(0.5)), vec2(0.3), 0.1) + 0.5);
        // mat.emitted = ZEROSPECT;
        // mat.specularBounce = false;

        // -----------------------------------------------------------------------
        // Because we are in the GPU, we do have access to differentials directly
        // This wouldn't be the case in a regular raytrace.
        // It wouldn't work as well in shaders doing interleaved calculations in
        // pixels (such as some of the 3D/stereo shaders here in Shadertoy)
        // -----------------------------------------------------------------------
        
        // calc texture sampling footprint
        vec3 ddx_uv = dFdx(X);
        vec3 ddy_uv = dFdy(X);
        float c = checkersTextureGradTri(1.99*X, ddx_uv, ddy_uv);
        
        // mat.albedo = spect_mul(WHITESUR, (c - 0.2) * 4.0 - 0.9);
        mat.albedo = spect_smoothstep(BLUE, YELLOW, (c - 0.2) * 4.0 - 0.9);
        mat.type = LAMB;
        mat.emitted = ZEROSPECT;
        mat.specularBounce = false;
    } else if (which_mat == WALL) {
        mat.type = LAMB;
        mat.albedo = HALF;
        mat.emitted = ZEROSPECT;
        mat.specularBounce = false;
    } else if (which_mat == DOTS) {
        mat.type = LAMB;
        mat.albedo = spect_smoothstep(BLUESUR, YELLOWSUR, disk(mod(fract(0.7 * X.xz), vec2(0.5)), vec2(0.3), 0.1) + 0.5);
        mat.emitted = ZEROSPECT;
        mat.specularBounce = false;
    } else if (which_mat == TEXTURE) {
        vec2 uv = X.xy*0.11 + vec2(0.5, -0.3) - 0.8;
        float idx = texture(voronoi_indexed_tex, mod(0.4 * uv + 0.5, vec2(0.9, 0.9))).r;
        Spectrum alb = pal(idx, vec3(0.5), vec3(0.5), vec3(1.0), vec3(0.3, 0.20, 0.20));

        mat.type = LAMB;
        mat.albedo = alb;
        mat.emitted = ZEROSPECT;
        mat.specularBounce = false;
    } else if (which_mat == METAL) {
        mat.type = METAL;
        mat.albedo = UNITY;
        mat.emitted = ZEROSPECT;
        mat.specularBounce = true;
    } else if (which_mat == GLASS) {
        mat.type = GLASS;
        mat.albedo = UNITY;
        mat.emitted = ZEROSPECT;
        mat.specularBounce = true;
    }

    return mat;
}

float FrDielectric(float cosThetaI, float etaI, float etaT) {
    cosThetaI = clamp(cosThetaI, - 1.0, 1.0);
    
    // Potentially swap indices of refraction
    bool entering = cosThetaI > 0.0;
    if (!entering) {
        float tmp = etaT;
        etaT = etaI;
        etaI = tmp;
        cosThetaI = abs(cosThetaI);
    }
    
    // Compute _cosThetaT_ using Snell's law
    float sinThetaI = sqrt(max(0.0, 1.0 - cosThetaI * cosThetaI));
    float sinThetaT = etaI / etaT * sinThetaI;
    
    // Handle total internal reflection
    if (sinThetaT >= 1.0) return 1.0;

    float cosThetaT = sqrt(max(0.0, 1.0 - sinThetaT * sinThetaT));

    float Rparl = ((etaT * cosThetaI) - (etaI * cosThetaT)) /
            ((etaT * cosThetaI) + (etaI * cosThetaT));
    float Rperp = ((etaI * cosThetaI) - (etaT * cosThetaT)) /
            ((etaI * cosThetaI) + (etaT * cosThetaT));
    return (Rparl * Rparl + Rperp * Rperp) / 2.0;
}

float FrConductor(float cosThetaI, float etai, float etat, float k) {
    cosThetaI = clamp(cosThetaI, -1.0, 1.0);
    float eta = etat / etai;
    float etak = k / etai;
    
    float cosThetaI2 = cosThetaI * cosThetaI;
    float sinThetaI2 = 1.0 - cosThetaI2;
    float eta2 = eta * eta;
    float etak2 = etak * etak;
    
    float t0 = eta2 - etak2 - sinThetaI2;
    float a2plusb2 = sqrt(t0 * t0 + 4.0 * eta2 * etak2);
    float t1 = a2plusb2 + cosThetaI2;
    float a = sqrt(0.5 * (a2plusb2 + t0));
    float t2 = 2.0 * cosThetaI * a;
    float Rs = (t1 - t2) / (t1 + t2);
    
    float t3 = cosThetaI2 * a2plusb2 + sinThetaI2 * sinThetaI2;
    float t4 = t2 * sinThetaI2;
    float Rp = Rs * (t3 - t4) / (t3 + t4);
    
    return 0.5 * (Rp + Rs);
}
`
