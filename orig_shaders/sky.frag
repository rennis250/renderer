#version 410 core

out vec4 fragColor;

in vec2 TexCoord;

uniform sampler2D blue_noise_tex;

uniform vec2 iResolution;
uniform float iTime;
uniform float iFrame;

uniform int whichillum;

#define NBOUNCES 10
#define NSAMPLES 300

#define ZERO min(int(iFrame), 0)

#define Pi 3.14159265358979323846
#define InvPi 0.31830988618379067154
#define Inv2Pi 0.15915494309189533577
#define Inv4Pi 0.07957747154594766788
#define PiOver2 1.57079632679489661923
#define PiOver4 0.78539816339744830961
#define Sqrt2 1.41421356237309504880

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

struct Ray {
    vec3 origin;
    vec3 direction;
};

struct Spectrum {
    vec4 p1;
    vec4 p2;
    vec4 p3;
};

Spectrum spect_add(in Spectrum sa, in Spectrum sb) {
    return Spectrum(sa.p1 + sb.p1, sa.p2 + sb.p2, sa.p3 + sb.p3);
}

Spectrum spect_add(in Spectrum sa, in float c) {
    return Spectrum(sa.p1 + c, sa.p2 + c, sa.p3 + c);
}

Spectrum spect_subtract(in Spectrum sa, in Spectrum sb) {
    return Spectrum(sa.p1 - sb.p1, sa.p2 - sb.p2, sa.p3 - sb.p3);
}

Spectrum spect_subtract(in Spectrum sa, in float c) {
    return Spectrum(sa.p1 - c, sa.p2 - c, sa.p3 - c);
}

Spectrum spect_mul(in Spectrum sa, in Spectrum sb) {
    return Spectrum(sa.p1 * sb.p1, sa.p2 * sb.p2, sa.p3 * sb.p3);
}

Spectrum spect_mul(in Spectrum sa, in float c) {
    return Spectrum(sa.p1 * c, sa.p2 * c, sa.p3 * c);
}

Spectrum spect_div(in Spectrum sa, in Spectrum sb) {
    return Spectrum(sa.p1 / sb.p1, sa.p2 / sb.p2, sa.p3 / sb.p3);
}

Spectrum spect_div(in Spectrum sa, in float c) {
    return Spectrum(sa.p1 / c, sa.p2 / c, sa.p3 / c);
}

Spectrum spect_mix(in Spectrum sa, in Spectrum sb, in float t) {
    return Spectrum(mix(sa.p1, sb.p1, t), mix(sa.p2, sb.p2, t), mix(sa.p3, sb.p3, t));
}

Spectrum spect_smoothstep(in Spectrum sa, in Spectrum sb, in float t) {
    return Spectrum(vec4(smoothstep(sa.p1.x, sb.p1.x, t),
        smoothstep(sa.p1.y, sb.p1.y, t),
        smoothstep(sa.p1.z, sb.p1.z, t),
        smoothstep(sa.p1.w, sb.p1.w, t)),
    
    vec4(smoothstep(sa.p2.x, sb.p2.x, t),
        smoothstep(sa.p2.y, sb.p2.y, t),
        smoothstep(sa.p2.z, sb.p2.z, t),
        smoothstep(sa.p2.w, sb.p2.w, t)),
    
    vec4(smoothstep(sa.p3.x, sb.p3.x, t),
        smoothstep(sa.p3.y, sb.p3.y, t),
        smoothstep(sa.p3.z, sb.p3.z, t),
        smoothstep(sa.p3.w, sb.p3.w, t)));
}

float spect_max(in Spectrum s) {
    float ma1 = max(s.p1.x, max(s.p1.y, max(s.p1.z, s.p1.w)));
    float ma2 = max(s.p2.x, max(s.p2.y, max(s.p2.z, s.p2.w)));
    float ma3 = max(s.p3.x, max(s.p3.y, max(s.p3.z, s.p3.w)));
    return max(ma1, max(ma2, ma3));
}

Spectrum spect_exp(in Spectrum s) {
    vec4 e1 = exp(s.p1);
    vec4 e2 = exp(s.p2);
    vec4 e3 = exp(s.p3);
    return Spectrum(e1, e2, e3);
}

Spectrum spect_sqrt(in Spectrum s) {
    vec4 sq1 = sqrt(s.p1);
    vec4 sq2 = sqrt(s.p2);
    vec4 sq3 = sqrt(s.p3);
    return Spectrum(sq1, sq2, sq3);
}

#define NWLNS 12

const vec3 CMFs[NWLNS] = vec3[NWLNS](
    vec3(0.3903, 0.0108, 1.8505),
    vec3(6.7842, 0.2493, 32.9414),
    vec3(8.7453, 1.2821, 47.6772),
    vec3(2.2018, 4.0746, 20.0946),
    vec3(0.1956, 13.2202, 4.5711),
    vec3(6.6255, 25.2773, 0.7354),
    vec3(17.8760, 26.8398, 0.0820),
    vec3(28.1812, 20.3419, 0.0292),
    vec3(24.2010, 10.9868, 0.0057),
    vec3(9.6361, 3.6853, 0.0002),
    vec3(1.9999, 0.7307, 0.0000),
    vec3(0, 0, 0)
);

vec3 spect_to_rgb(in Spectrum s) {
    vec3 xyz = vec3(0.0);

    xyz += CMFs[0]  * s.p1.x;
    xyz += CMFs[1]  * s.p1.y;
    xyz += CMFs[2]  * s.p1.z;
    xyz += CMFs[3]  * s.p1.w;

    xyz += CMFs[4]  * s.p2.x;
    xyz += CMFs[5]  * s.p2.y;
    xyz += CMFs[6]  * s.p2.z;
    xyz += CMFs[7]  * s.p2.w;
    
    xyz += CMFs[8]  * s.p3.x;
    xyz += CMFs[9]  * s.p3.y;
    xyz += CMFs[10] * s.p3.z;
    xyz += CMFs[11] * s.p3.w;

    mat3 xyz_to_rgb = mat3(
        vec3(3.2406, -0.9689, 0.0557),
        vec3(-1.5372, 1.8758, -0.2040),
        vec3(-0.4986, 0.0415, 1.0570)
    );

    return xyz_to_rgb*xyz;
}

#define UNITY Spectrum(vec4(1, 1, 1, 1), vec4(1, 1, 1, 1), vec4(1, 1, 1, 1))
#define UNITY_N 106.6988
#define ZEROSPECT Spectrum(vec4(0, 0, 0, 0), vec4(0, 0, 0, 0), vec4(0, 0, 0, 0))

#define WHITESUR Spectrum(vec4(0.343, 0.72054, 0.77327, 0.75047), vec4(0.73394, 0.73291, 0.72675, 0.73457), vec4(0.74717, 0.7158, 0.75194, 0.7356))
#define WHITESUR_N 78.3062
#define GREENSUR Spectrum(vec4(0.092, 0.096683, 0.10227, 0.13182), vec4(0.40227, 0.44626, 0.31621, 0.19174), vec4(0.12782, 0.11574, 0.13324, 0.15868))
#define GREENSUR_N 31.6074
#define REDSUR Spectrum(vec4(0.04, 0.054636, 0.060382, 0.059549), vec4(0.055277, 0.058049, 0.067844, 0.1817), vec4(0.50124, 0.63632, 0.62978, 0.64319))
#define REDSUR_N 16.3615
#define WHITEILLUM Spectrum(vec4(0, 2.1818, 4.3636, 6.5455), vec4(8.6909, 10.7636, 12.8364, 14.9091), vec4(16.1091, 16.8727, 17.6364, 18.4))
#define WHITEILLUM_N 1319.6407
#define UNITY Spectrum(vec4(1, 1, 1, 1), vec4(1, 1, 1, 1), vec4(1, 1, 1, 1))
#define UNITY_N 106.6988
#define ZEROSPECT Spectrum(vec4(0, 0, 0, 0), vec4(0, 0, 0, 0), vec4(0, 0, 0, 0))
#define ZEROSPECT_N 0
#define BLUEILLUM Spectrum(vec4(0.031549, 0.071507, 0.14412, 0.26378), vec4(0.44571, 0.70435, 1.0518, 1.4968), vec4(2.0439, 2.6932, 3.4406, 4.2785))
#define BLUEILLUM_N 118.5474
#define YELLOWILLUM Spectrum(vec4(2.0133, 1.9481, 2.0776, 1.7557), vec4(1.4298, 1.1795, 0.97182, 0.82488), vec4(0.73179, 0.64452, 0.60555, 0.52057))
#define YELLOWILLUM_N 112.7625
#define REDPRIM Spectrum(vec4(0.0012888, 0.0035377, 0.0060555, 0.0068225), vec4(0.0070749, 0.017122, 0.15541, 0.58273), vec4(1, 0.82653, 0.41649, 0.18382))
#define REDPRIM_N 30.9251
#define GREENPRIM Spectrum(vec4(0.0015031, 0.005478, 0.089139, 0.43421), vec4(0.93476, 1, 0.57713, 0.18899), vec4(0.040937, 0.010657, 0.0041689, 0.0020919))
#define GREENPRIM_N 59.3465
#define BLUEPRIM Spectrum(vec4(0.1884, 0.5944, 1, 0.79649), vec4(0.31772, 0.064679, 0.010412, 0.0069007), vec4(0.0091953, 0.0074626, 0.0039193, 0.0019485))
#define BLUEPRIM_N 11.0641
#define BLUESUR Spectrum(vec4(0.70507, 0.86353, 1, 0.92747), vec4(0.66104, 0.39434, 0.24311, 0.19277), vec4(0.18943, 0.19758, 0.20284, 0.20418))
#define BLUESUR_N 37.3949
#define YELLOWSUR Spectrum(vec4(0.091609, 0.093483, 0.097796, 0.14636), vec4(0.31904, 0.60099, 0.82454, 0.92035), vec4(0.94373, 0.95753, 0.97963, 1))
#define YELLOWSUR_N 75.6207
#define REDT Spectrum(vec4(0.1356, 0.3529, 0.62845, 0.64258), vec4(0.46056, 0.39138, 0.49095, 0.73642), vec4(1, 0.93571, 0.57224, 0.30199))
#define REDT_N 62.5055
#define GREENT Spectrum(vec4(0.13602, 0.33801, 0.63546, 0.8576), vec4(1, 0.94693, 0.73657, 0.57353), vec4(0.48263, 0.39272, 0.30013, 0.23531))
#define GREENT_N 79.9561
#define BLUET Spectrum(vec4(0.17754, 0.47993, 0.88611, 1), vec4(0.82436, 0.57558, 0.37316, 0.30928), vec4(0.33312, 0.28798, 0.17666, 0.10117))
#define BLUET_N 51.9367
#define YELLOWT Spectrum(vec4(0.10012, 0.23794, 0.48142, 0.73897), vec4(0.93542, 0.91019, 0.76972, 0.82155), vec4(1, 0.89736, 0.54427, 0.29565))
#define YELLOWT_N 91.1251

Spectrum skyColor(vec3 wi) {
    Spectrum illum_col = UNITY;
    float N = 1.0 / UNITY_N;
    illum_col = spect_mul(illum_col, 1.0*N);
    
    Spectrum col;
    float corrFactor;
    if (whichillum == 0) {
        col = BLUEILLUM;
        corrFactor = BLUEILLUM_N;
    } else {
        col = YELLOWILLUM;
        corrFactor = YELLOWILLUM_N;
    }

    vec3 sunDir = vec3(0.0, 0.0, 0.8);

    col = spect_mix(spect_mul(col, 1.0 / corrFactor), illum_col, 0.5 + 0.5 * -wi.y);
    float sun = clamp(dot(normalize(sunDir), wi), 0.0, 1.0);
    col = spect_add(col, spect_mul(illum_col, pow(sun, 44.0) + 10.0 * pow(sun, 256.0)));

    return col;
}

Ray getPrimaryRay(vec2 uv) {
    Ray r;

    r.origin = vec3(0.0, 0.0, 0.0);

    vec2 p = ConcentricSampleDisk();
    r.origin.xy += p*0.004;

    vec3 vup = vec3(0.0, 1.0, 0.0);

    vec3 cw = normalize(vec3(0.0, 0.0, 12.0));
    vec3 cu = normalize(cross(cw, vup));
    vec3 cv = normalize(cross(cu, cw));

    mat3 ca = mat3(cu, cv, cw);
        
    r.direction = ca * normalize(vec3(uv.x, uv.y, -2.0));

    return r;
}

// from Mitsuba
vec3 reinhardTonemap(vec3 p, float key, float burn) {
    burn = min(1.0, max(1e-8, 1.0 - burn));

    float logAvgLuminance = log(2.0);
    float maxLuminance = 50.0;
    float scale = key / logAvgLuminance;
    float Lwhite = maxLuminance * scale;

    /* Having the 'burn' parameter scale as 1/b^4 provides a nicely behaved knob */
    float invWp2 = 1.0 / (Lwhite * Lwhite * pow(burn, 4.0));

    /* Convert ITU-R Rec. BT.709 linear RGB to XYZ tristimulus values */
    float X = p.r * 0.412453 + p.g * 0.357580 + p.b * 0.180423;
    float Y = p.r * 0.212671 + p.g * 0.715160 + p.b * 0.072169;
    float Z = p.r * 0.019334 + p.g * 0.119193 + p.b * 0.950227;

    /* Convert to xyY */
    float normalization = 1.0 / (X + Y + Z);
    float x = X * normalization;
    float y = Y * normalization;
    float Lp = Y * scale;

    /* Apply the tonemapping transformation */
    Y = Lp * (1.0 + Lp * invWp2) / (1.0 + Lp);

    /* Convert back to XYZ */
    float ratio = Y / y;
    X = ratio * x;
    Z = ratio * (1.0 - x - y);

    /* Convert from XYZ tristimulus values to ITU-R Rec. BT.709 linear RGB */
    vec3 outc;
    outc.r = 3.240479 * X + -1.537150 * Y + -0.498535 * Z;
    outc.g = -0.969256 * X + 1.875991 * Y + 0.041556 * Z;
    outc.b = 0.055648 * X + -0.204043 * Y + 1.057311 * Z;

    return outc;
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

void main() {
    vec2 uv = -1.0 + 2.0 * gl_FragCoord.xy / iResolution.xy;
    uv.x *= iResolution.x / iResolution.y;
    
    int nsamps = int(1.0 / (1.0 + 0.5 * length(uv - 0.5)) * float(NSAMPLES));
    float oneOverSPP = 1.0 / float(nsamps);
    float strataSize = oneOverSPP;
    Spectrum L;
    for (int i = ZERO; i < nsamps; ++i) {
        vec2 subPixelJitter = halton(i % nsamps + 1) - 0.5;
        uv = TexCoord - vec2(0.5) + 0.003 * vec2(strataSize * (float(i) + subPixelJitter.x), subPixelJitter.y);
        vec3 rd = normalize(vec3(uv, 1.0));

        L = spect_add(L, skyColor(rd));
    }
    vec3 col = spect_to_rgb(L);
    col *= oneOverSPP;
    col = clamp(col, 0.0, 1000.0); // prevent NaN and Inf

    // vignetting
    vec2 p = gl_FragCoord.xy / iResolution.xy;
    col *= 0.5 + 0.5*pow(16.0*p.x*p.y*(1.0-p.x)*(1.0-p.y), 0.3);

    fragColor = vec4(reinhardTonemap(col, 0.8, 0.1), 1.0);
}