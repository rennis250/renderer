package renderer

const spectra_glsl = `
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

// mitsuba's
// const vec3 CMFs[NWLNS] = vec3[NWLNS](
    // vec3(0.0775,    0.0023,    0.3719),
    // vec3(0.3183,    0.0207,    1.5953),
    // vec3(0.2596,    0.0694,    1.5308),
    // vec3(0.0541,    0.1977,    0.5782),
    // vec3(0.0313,    0.5573,    0.1442),
    // vec3(0.2622,    0.9221,    0.0275),
    // vec3(0.6370,    0.9770,    0.0039),
    // vec3(0.9874,    0.7826,    0.0013),
    // vec3(0.9512,    0.4735,    0.0004),
    // vec3(0.5015,    0.2014,    0.0000),
    // vec3(0.1523,    0.0565,         0),
    // vec3(0.0309,    0.0112,         0)
// );

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

    // mitsuba's
    // float scale = 0.234096;

    float scale = 0.0094;
    return xyz_to_rgb*xyz * scale;
}

#define WHITE Spectrum(vec4(0.343, 0.72054, 0.77327, 0.75047), vec4(0.73394, 0.73291, 0.72675, 0.73457), vec4(0.74717, 0.7158, 0.75194, 0.7356))
#define UNITY Spectrum(vec4(1.0, 1.0, 1.0, 1.0), vec4(1.0, 1.0, 1.0, 1.0), vec4(1.0, 1.0, 1.0, 1.0))
#define HALF Spectrum(vec4(0.5, 0.5, 0.5, 0.5), vec4(0.5, 0.5, 0.5, 0.5), vec4(0.5, 0.5, 0.5, 0.5))
#define ZEROSPECT Spectrum(vec4(0, 0, 0, 0), vec4(0, 0, 0, 0), vec4(0, 0, 0, 0))

//#define GREEN Spectrum(vec4(0.092, 0.096683, 0.10227, 0.13182), vec4(0.40227, 0.44626, 0.31621, 0.19174), vec4(0.12782, 0.11574, 0.13324, 0.15868))
//#define RED Spectrum(vec4(0.04, 0.054636, 0.060382, 0.059549), vec4(0.055277, 0.058049, 0.067844, 0.1817), vec4(0.50124, 0.63632, 0.62978, 0.64319))
//#define BLUE Spectrum(vec4(0.70507, 0.86353, 1.0, 0.92747), vec4(0.66104, 0.39434, 0.24311, 0.19277), vec4(0.18943, 0.19758, 0.20284, 0.20418))
//#define YELLOW Spectrum(vec4(0.091609, 0.093483, 0.097796, 0.14636), vec4(0.31904, 0.60099, 0.82454, 0.92035), vec4(0.94373, 0.95753, 0.97963, 1))

#define RED spect_mul(Spectrum(vec4(0.1356, 0.3529, 0.62845, 0.64258), vec4(0.46056, 0.39138, 0.49095, 0.73642), vec4(1, 0.93571, 0.57224, 0.30199)), 0.7)
#define GREEN spect_mul(Spectrum(vec4(0.13602, 0.33801, 0.63546, 0.8576), vec4(1, 0.94693, 0.73657, 0.57353), vec4(0.48263, 0.39272, 0.30013, 0.23531)), 0.7)
#define BLUE spect_mul(Spectrum(vec4(0.17754, 0.47993, 0.88611, 1), vec4(0.82436, 0.57558, 0.37316, 0.30928), vec4(0.33312, 0.28798, 0.17666, 0.10117)), 0.7)
#define YELLOW spect_mul(Spectrum(vec4(0.10012, 0.23794, 0.48142, 0.73897), vec4(0.93542, 0.91019, 0.76972, 0.82155), vec4(1, 0.89736, 0.54427, 0.29565)), 0.7)

#define BLUEILLUM Spectrum(vec4(0.031549, 0.071507, 0.14412, 0.26378), vec4(0.44571, 0.70435, 1.0518, 1.4968), vec4(2.0439, 2.6932, 3.4406, 4.2785))
#define YELLOWILLUM Spectrum(vec4(2.0133, 1.9481, 2.0776, 1.7557), vec4(1.4298, 1.1795, 0.97182, 0.82488), vec4(0.73179, 0.64452, 0.60555, 0.52057))

#define GREENSUR Spectrum(vec4(0.092, 0.096683, 0.10227, 0.13182), vec4(0.40227, 0.44626, 0.31621, 0.19174), vec4(0.12782, 0.11574, 0.13324, 0.15868))
#define REDSUR Spectrum(vec4(0.04, 0.054636, 0.060382, 0.059549), vec4(0.055277, 0.058049, 0.067844, 0.1817), vec4(0.50124, 0.63632, 0.62978, 0.64319))
#define BLUESUR Spectrum(vec4(0.70507, 0.86353, 1, 0.92747), vec4(0.66104, 0.39434, 0.24311, 0.19277), vec4(0.18943, 0.19758, 0.20284, 0.20418))
#define YELLOWSUR Spectrum(vec4(0.091609, 0.093483, 0.097796, 0.14636), vec4(0.31904, 0.60099, 0.82454, 0.92035), vec4(0.94373, 0.95753, 0.97963, 1))

#define REDPRIM Spectrum(vec4(0.0012888, 0.0035377, 0.0060555, 0.0068225), vec4(0.0070749, 0.017122, 0.15541, 0.58273), vec4(1, 0.82653, 0.41649, 0.18382))
#define GREENPRIM Spectrum(vec4(0.0015031, 0.005478, 0.089139, 0.43421), vec4(0.93476, 1, 0.57713, 0.18899), vec4(0.040937, 0.010657, 0.0041689, 0.0020919))
#define BLUEPRIM Spectrum(vec4(0.1884, 0.5944, 1, 0.79649), vec4(0.31772, 0.064679, 0.010412, 0.0069007), vec4(0.0091953, 0.0074626, 0.0039193, 0.0019485))

#define REDT Spectrum(vec4(0.1356, 0.3529, 0.62845, 0.64258), vec4(0.46056, 0.39138, 0.49095, 0.73642), vec4(1, 0.93571, 0.57224, 0.30199))
#define GREENT Spectrum(vec4(0.13602, 0.33801, 0.63546, 0.8576), vec4(1, 0.94693, 0.73657, 0.57353), vec4(0.48263, 0.39272, 0.30013, 0.23531))
#define BLUET Spectrum(vec4(0.17754, 0.47993, 0.88611, 1), vec4(0.82436, 0.57558, 0.37316, 0.30928), vec4(0.33312, 0.28798, 0.17666, 0.10117))
#define YELLOWT Spectrum(vec4(0.10012, 0.23794, 0.48142, 0.73897), vec4(0.93542, 0.91019, 0.76972, 0.82155), vec4(1, 0.89736, 0.54427, 0.29565))
`
