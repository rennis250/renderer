package renderer

import "fmt"

type Shape interface {
	ShapeToGLSLString() string
}

type Sphere struct {
	Position [3]float64
	Radius   float64
	Material Material
	Color    Color
}

func (s *Sphere) ShapeToGLSLString() string {
	str := `vec4 sph = vec4(` + fmt.Sprintf("%f", s.Position[0]) + ", " + fmt.Sprintf("%f", s.Position[1]) + ", " + fmt.Sprintf("%f", s.Position[2]) + ", " + fmt.Sprintf("%f", s.Radius)

	str += `);
    t = iSphere(X, wi, sph);
    if (t > EPS && t < s.t) {
        s.hit = true;
        s.t = t;
        s.position = X + t * wi;
        s.shadingNormal = nSphere(s.position, sph);
        s.mat = getMaterial(` + materialString(s.Material) + `, s.position, s.shadingNormal);
        s.mat.albedo = ` + s.Color.colorToGLSLString() + `;
    }`

	return str
}

type Triangle struct {
	Center, Point1, Point2, Point3 [3]float64
	Material                       Material
	Color                          Color
}

func (t *Triangle) ShapeToGLSLString() string {
	cent := `vec3(` + fmt.Sprintf("%f", t.Center[0]) + ", " + fmt.Sprintf("%f", t.Center[1]) + ", " + fmt.Sprintf("%f", t.Center[2]) + `)`
	tri0 := `vec3(` + fmt.Sprintf("%f", t.Point1[0]) + ", " + fmt.Sprintf("%f", t.Point1[1]) + ", " + fmt.Sprintf("%f", t.Point1[2]) + ")"
	tri1 := `vec3(` + fmt.Sprintf("%f", t.Point2[0]) + ", " + fmt.Sprintf("%f", t.Point2[1]) + ", " + fmt.Sprintf("%f", t.Point2[2]) + ")"
	tri2 := `vec3(` + fmt.Sprintf("%f", t.Point3[0]) + ", " + fmt.Sprintf("%f", t.Point3[1]) + ", " + fmt.Sprintf("%f", t.Point3[2]) + ")"

	str := `t = iTriangle(X - ` + cent + `, wi, vec2(EPS, s.t), nor, ` + tri0 + `, ` + tri1 + `, ` + tri2 + `);
	if (t > EPS && t < s.t) {
		s.hit = true;
		s.t = t;
		s.position = X + t*wi;
		s.shadingNormal = nor;
        s.mat = getMaterial(` + materialString(t.Material) + `, s.position, s.shadingNormal);
        // s.mat.albedo = ` + t.Color.colorToGLSLString() + `;
    }`

	return str
}

type Plane struct {
	Orientation [3]float64
	Offset      float64
	Material    Material
	Color       Color
}

func (p *Plane) ShapeToGLSLString() string {
	str := `vec4 pl = vec4(` + fmt.Sprintf("%f", p.Orientation[0]) + ", " + fmt.Sprintf("%f", p.Orientation[1]) + ", " + fmt.Sprintf("%f", p.Orientation[2]) + ", " + fmt.Sprintf("%f", p.Offset)

	str += `);
    t = iPlane(X, wi, pl);
    if (t > EPS && t < s.t) {
        s.hit = true;
        s.t = t;
        s.position = X + t * wi;
        s.shadingNormal = nPlane(pl);
        s.mat = getMaterial(` + materialString(p.Material) + `, s.position, s.shadingNormal);
        // s.mat.albedo = ` + p.Color.colorToGLSLString() + `;
    }`

	return str
}

type Cone struct {
	Position, PointA, PointB [3]float64
	RadiusA, RadiusB         float64
	Material                 Material
	Color                    Color
}

func (c *Cone) ShapeToGLSLString() string {
	pos := `vec3(` + fmt.Sprintf("%f", c.Position[0]) + ", " + fmt.Sprintf("%f", c.Position[1]) + ", " + fmt.Sprintf("%f", c.Position[2]) + `)`

	pa := `vec3(` + fmt.Sprintf("%f", c.PointA[0]) + ", " + fmt.Sprintf("%f", c.PointA[1]) + ", " + fmt.Sprintf("%f", c.PointA[2]) + `)`
	pb := `vec3(` + fmt.Sprintf("%f", c.PointB[0]) + ", " + fmt.Sprintf("%f", c.PointB[1]) + ", " + fmt.Sprintf("%f", c.PointB[2]) + `)`

	ra := fmt.Sprintf("%f", c.RadiusA)
	rb := fmt.Sprintf("%f", c.RadiusB)

	str := `t = iRoundedCone(X - ` + pos + `, wi, vec2(EPS, s.t), nor, ` + pa + ` , ` + pb + `, ` + ra + `, ` + rb + `);
    if (t > EPS && t < s.t) {
        s.hit = true;
        s.t = t;
        s.position = X + t * wi;
        s.shadingNormal = nor;
        s.mat = getMaterial(` + materialString(c.Material) + `, s.position, s.shadingNormal);
        s.mat.albedo = ` + c.Color.colorToGLSLString() + `;
    }`

	return str
}

type Box struct {
	Position [3]float64
	Size     [3]float64
	Material Material
	Color    Color
}

func (b *Box) ShapeToGLSLString() string {
	pos := `vec3(` + fmt.Sprintf("%f", b.Position[0]) + ", " + fmt.Sprintf("%f", b.Position[1]) + ", " + fmt.Sprintf("%f", b.Position[2]) + `)`
	sz := `vec3(` + fmt.Sprintf("%f", b.Size[0]) + ", " + fmt.Sprintf("%f", b.Size[1]) + ", " + fmt.Sprintf("%f", b.Size[2]) + `)`

	str := `t = iBox(X - ` + pos + `, wi, vec2(EPS, s.t), nor, ` + sz + `);
    if (t > EPS && t < s.t) {
        s.hit = true;
        s.t = t;
        s.position = X + t * wi;
        s.shadingNormal = nor;
        s.mat = getMaterial(` + materialString(b.Material) + `, s.position, s.shadingNormal);
        s.mat.albedo = ` + b.Color.colorToGLSLString() + `;
    }`

	return str
}

const shape_intersect_glsl = `
bool isInside = false;

float iPlane(vec3 ro, vec3 rd, vec4 pl) {
    float a = dot(rd, pl.xyz);
    float d = -(dot(ro, pl.xyz) + pl.w) / a;
    if (a > 0.0 || d < EPS || d > MAX_SCENE_DIST) {
        return MAX_DIST;
    } else {
        return d;
    }
}

vec3 nPlane(vec4 pl) {
    return pl.xyz;
}

float iSphere(vec3 ro, vec3 rd, vec4 sph) {
    vec3 oc = ro - sph.xyz;
    float a = dot(rd, rd);
    float b = dot(oc, rd);
    float c = dot(oc, oc) - sph.w * sph.w;
    float d = b*b - a*c;
    if (d < 0.0) return MAX_DIST;
    
    d = sqrt(d);
    float t0 = (-b - d) / a;
    float t1 = (-b + d) / a;
    
    if (t0 >= EPS && t0 <= MAX_SCENE_DIST) {
        isInside = false;
        return t0;
    } else if (t1 >= EPS && t1 <= MAX_SCENE_DIST) {
        isInside = true;
        return t1;
    } else {
        return MAX_DIST;
    }
}

vec3 nSphere(vec3 ro, vec4 sph) {
    return (ro - sph.xyz) / sph.w;
}

// Box: https://www.shadertoy.com/view/ld23DV
float iBox(in vec3 ro, in vec3 rd, in vec2 dbound, inout vec3 nor, in vec3 boxSize) {
    vec3 m = sign(rd) / max(abs(rd), 1e-8);
    vec3 n = m*ro;
    vec3 k = abs(m) * boxSize;
    
    vec3 t1 = -n - k;
    vec3 t2 = -n + k;
    
    float tN = max(max(t1.x, t1.y), t1.z);
    float tF = min(min(t2.x, t2.y), t2.z);
    
    if (tN > tF || tF <= 0.0) {
        return MAX_DIST;
    } else {
        if (tN >= dbound.x && tN <= dbound.y) {
            isInside = false;
            nor = -sign(rd) * step(t1.yzx, t1.xyz) * step(t1.zxy, t1.xyz);
            return tN;
        } else if (tF >= dbound.x && tF <= dbound.y) {
            isInside = true;
            nor = -sign(rd) * step(t2.yzx, t2.xyz) * step(t2.zxy, t2.xyz);
            return tF;
        } else {
            return MAX_DIST;
        }
    }
}

// Rounded Cone:    https://www.shadertoy.com/view/MlKfzm
float iRoundedCone(const in vec3 ro, const in vec3 rd, const in vec2 distBound, inout vec3 normal, const in vec3 pa, const in vec3 pb, const in float ra, const in float rb) {
    vec3 ba = pb - pa;
    vec3 oa = ro - pa;
    vec3 ob = ro - pb;
    float rr = ra - rb;
    float m0 = dot(ba, ba);
    float m1 = dot(ba, oa);
    float m2 = dot(ba, rd);
    float m3 = dot(rd, oa);
    float m5 = dot(oa, oa);
    float m6 = dot(ob, rd);
    float m7 = dot(ob, ob);
    
    float d2 = m0 - rr * rr;
    
    float k2 = d2 - m2 * m2;
    float k1 = d2 * m3 - m1 * m2 + m2 * rr * ra;
    float k0 = d2 * m5 - m1 * m1 + m1 * rr * ra * 2.0 - m0 * ra * ra;
    
    float h = k1 * k1 - k0 * k2;
    if (h < 0.0) {
        return MAX_DIST;
    }
    
    float t = (-sqrt(h) - k1) / k2;
    
    float y = m1 - ra * rr + t * m2;
    if (y > 0.0 && y < d2) {
        if (t >= distBound.x && t <= distBound.y) {
            normal = normalize(d2 * (oa + t * rd) - ba * y);
            return t;
        } else {
            return MAX_DIST;
        }
    } else {
        float h1 = m3 * m3 - m5 + ra * ra;
        float h2 = m6 * m6 - m7 + rb * rb;
        
        if (max(h1, h2) < 0.0) {
            return MAX_DIST;
        }
        
        vec3 n = vec3(0);
        float r = MAX_DIST;
        
        if (h1 > 0.0) {
            r = -m3 - sqrt(h1);
            n = (oa + r * rd) / ra;
        }
        if (h2 > 0.0) {
            t = -m6 - sqrt(h2);
            if (t < r) {
                n = (ob + t * rd) / rb;
                r = t;
            }
        }
        if (r >= distBound.x && r <= distBound.y) {
            normal = n;
            return r;
        } else {
            return MAX_DIST;
        }
    }
}

// Triangle:        https://www.shadertoy.com/view/MlGcDz
float iTriangle(const in vec3 ro, const in vec3 rd, const in vec2 distBound, inout vec3 nor, const in vec3 v0, const in vec3 v1, const in vec3 v2) {
    vec3 v1v0 = v1 - v0;
    vec3 v2v0 = v2 - v0;
    vec3 rov0 = ro - v0;
    
    vec3  n = cross(v1v0, v2v0);
    vec3  q = cross(rov0, rd);
    float d = 1.0/dot(rd, n);
    float u = d*dot(-q, v2v0);
    float v = d*dot(q, v1v0);
    float t = d*dot(-n, rov0);
    
    if( u < 0.0 || v < 0.0 || (u+v) > 1.0 || t < distBound.x || t > distBound.y) {
        return MAX_DIST;
    } else {
        nor = normalize(-n);
        return t;
    }
}
`
