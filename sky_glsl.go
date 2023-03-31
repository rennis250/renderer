package renderer

const sky_glsl = `
Spectrum skyColor(vec3 wi) {
    Spectrum illum_col = UNITY;
    Spectrum col = spect_mix(BLUEILLUM, YELLOWILLUM, smoothstep(-0.001, 0.001, 1.0*wi.x));

    // vec3 sunDir1 = vec3(-0.2, 0.8, -0.2);
    // vec3 sunDir2 = vec3(0.2, 0.8, -0.3);

    vec3 sunDir1 = vec3(-0.7, 0.8, 0.2);
    vec3 sunDir2 = vec3(0.1, 0.3, 0.3);

    col = spect_mix(col, illum_col, 0.5 + 0.5 * -wi.y);

    float sun = clamp(dot(normalize(sunDir1), wi), 0.0, 1.0);
    col = spect_add(col, spect_mul(illum_col, pow(sun, 44.0) + 10.0 * pow(sun, 256.0)));

    sun = clamp(dot(normalize(sunDir2), wi), 0.0, 1.0);
    col = spect_add(col, spect_mul(illum_col, pow(sun, 44.0) + 10.0 * pow(sun, 256.0)));

    return col;
}
`
