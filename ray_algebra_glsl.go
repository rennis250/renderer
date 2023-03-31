package renderer

const ray_algebra_glsl = `
vec3 ortho(vec3 v) {
    return abs(v.x) > abs(v.z) ? vec3(-v.y, v.x, 0.0) : vec3(0.0, -v.z, v.y);
}

mat3 ONB(vec3 n) {
    vec3 nn = normalize(n);
    vec3 o1 = normalize(ortho(nn));
    vec3 o2 = normalize(cross(nn, o1));
    return mat3(o1, nn, o2);
}

vec2 ConcentricSampleDisk() {
    vec2 u = rand2();
    vec2 uOffset = 2.0 * u - 1.0;
    if (uOffset.x == 0.0 && uOffset.y == 0.0) {
        return vec2(0.0);
    }

    float theta, r;
    if (abs(uOffset.x) > abs(uOffset.y)) {
        r = uOffset.x;
        theta = PiOver4 * (uOffset.y / uOffset.x);
    } else {
        r = uOffset.y;
        theta = PiOver2 - PiOver4 * (uOffset.x / uOffset.y);
    }

    return r * vec2(cos(theta), sin(theta));
}

vec3 cosWeightHemi() {
    vec2 d = ConcentricSampleDisk();
    float z = sqrt(max(0.0, 1.0 - d.x * d.x - d.y * d.y));
    return vec3(d.x, z, d.y);
}

struct Ray {
    vec3 origin;
    vec3 direction;
};
`
