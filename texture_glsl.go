package renderer

const texture_glsl = `
// modified from iq - Spectrum form
Spectrum pal(const in float t, const in vec3 a, const in vec3 b, const in vec3 c, const in vec3 d) {
    vec3 coords = a + b*cos(2.0 * Pi * (c * t + d));
    Spectrum rc = spect_mul(REDPRIM, coords.r);
    Spectrum gc = spect_mul(GREENPRIM, coords.g);
    Spectrum bc = spect_mul(BLUEPRIM, coords.b);
    return spect_add(rc, spect_add(gc, bc));
}

float disk(vec2 r, vec2 center, float radius) {
    return 1.0 - step(radius, length(r - center));
}

// from iq
vec3 pri(in vec3 x) {
    // see https://www.shadertoy.com/view/MtffWs
    vec3 h = fract(x / 2.0) - 0.5;
    return x * 0.5 + h*(1.0 - 2.0 * abs(h));
}

float checkersTextureGradTri(in vec3 p, in vec3 ddx, in vec3 ddy) {
    vec3 w = max(abs(ddx), abs(ddy)) + 0.01; // filter kernel
    vec3 i = (pri(p + w) - 2.0 * pri(p) + pri(p - w)) / (w * w); // analytical integral (box filter)
    return 0.5 - 0.5 * i.x * i.y * i.z; // xor pattern
}
`
