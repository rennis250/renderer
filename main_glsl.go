package renderer

const main_glsl = `
void main() {
    vec2 uv = -1.0 + 2.0 * gl_FragCoord.xy / iResolution.xy;
    uv.x *= iResolution.x / iResolution.y;
                
    int seed = int(gl_FragCoord.x + gl_FragCoord.y);

    vec4 b2 = texture(blue_noise_tex, uv);
    vec4 b1 = texture(blue_noise_tex, uv);
    float t = b1.r*b2.g;

    rng_initialize(gl_FragCoord.xy, seed + int(t) - int(t));
    
    // sample blue noise texture
    float bn1 = texelFetch(blue_noise_tex, shift2(), 0)[3];
    float bn2 = texelFetch(blue_noise_tex, shift2(), 0)[3];

    // decorrelate pixel seeds with blue noise
    pixel += ivec2(bn1, bn2);

    // Compute the average radiance into this pixel
    Spectrum L = ZEROSPECT;
    int nsamps = int(1.0 / (1.0 + 0.5 * length(uv - 0.5)) * float(NSAMPLES));
    float oneOverSPP = 1.0 / float(nsamps);
    float strataSize = oneOverSPP;
    for (int i = ZERO; i < nsamps; ++i) {
        // from demofox
        // calculate sub pixel jitter for anti aliasing with halton dist and stratified sampling
        vec2 subPixelJitter = halton(i % nsamps + 1) - 0.5;
        vec2 p = uv + 0.003 * vec2(strataSize * (float(i) + subPixelJitter.x), subPixelJitter.y);

        Ray r = getPrimaryRay(p);
        L = spect_add(L, L_i(r));
    }

    vec3 col = spect_to_rgb(L);
    col *= oneOverSPP;
    col = clamp(col, 0.0, 1000.0); // prevent NaN and Inf

    // vignetting
    vec2 p = gl_FragCoord.xy / iResolution.xy;
    col *= 0.5 + 0.5*pow(16.0*p.x*p.y*(1.0-p.x)*(1.0-p.y), 0.3);

    // aperture
    if (with_aperture == 1) {
        if (length(uv) > 0.9) {
            discard;
        }
    }

    // fragColor = vec4(col, 1.0);
    fragColor = vec4(reinhardTonemap(col, 0.8, 0.1), 1.0);
}`
