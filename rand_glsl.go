package renderer

const rand_glsl = `
// internal RNG state
uvec4 s0, s1;
ivec2 pixel;

void rng_initialize(vec2 p, int frame) {
    pixel = ivec2(p);

    // white noise seed
    s0 = uvec4(p, uint(frame), uint(p.x) + uint(p.y));

    // blue noise seed
    s1 = uvec4(frame, frame * 15843, frame * 31 + 4566, frame * 2345 + 58585);
}

// https://www.pcg-random.org/
void pcg4d(inout uvec4 v) {
    v = v * 1664525u + 1013904223u;
    v.x += v.y * v.w;
    v.y += v.z * v.x;
    v.z += v.x * v.y;
    v.w += v.y * v.z;
    v = v ^ (v >> 16u);
    v.x += v.y * v.w;
    v.y += v.z * v.x;
    v.z += v.x * v.y;
    v.w += v.y * v.z;
}

float rand() {
    pcg4d(s0);
    return float(s0.x) / float(0xffffffffu);
}

vec2 rand2() {
    pcg4d(s0);
    return vec2(s0.xy) / float(0xffffffffu);
}

vec3 rand3() {
    pcg4d(s0);
    return vec3(s0.xyz) / float(0xffffffffu);
}

// random blue noise sampling pos
ivec2 shift2() {
    pcg4d(s1);
    return (pixel + ivec2(s1.xy % 0x0fffffffu)) % 1024;
}

// halton low discrepancy sequence, from https://www.shadertoy.com/view/wdXSW8
vec2 halton(int index) {
    const vec2 coprimes = vec2(2.0, 3.0);
    vec2 s = vec2(index, index);
    vec4 a = vec4(1, 1, 0, 0);
    while(s.x > 0.0 && s.y > 0.0) {
        a.xy = a.xy / coprimes;
        a.zw += a.xy * mod(s, coprimes);
        s = floor(s / coprimes);
    }
    return a.zw;
}
`
