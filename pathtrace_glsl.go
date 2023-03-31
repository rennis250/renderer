package renderer

const pathtrace_glsl = `
Spectrum L_i(const Ray ray) {
    Spectrum L = ZEROSPECT;
    Spectrum beta = UNITY;
    float pdf;
    Spectrum fr;
    vec3 wo;
    float eta_for_RR;
    float etaScale;

    vec3 X = ray.origin;
    vec3 wi = ray.direction;
    for (int i = ZERO; i < NBOUNCES; i++) {
        Surfel surfelY = findFirstIntersection(X, wi);
        
        if(surfelY.hit) {
            wo = -wi;

            fr = finiteScatteringDensity(surfelY, wi, wo, pdf, eta_for_RR);
            X = surfelY.position;

            if(pdf == 0.0) {
                break;
            }

            float lamb_corr = abs(dot(wi, surfelY.shadingNormal));
            Spectrum lscaled = spect_mul(fr, lamb_corr);
            Spectrum pscaled = spect_div(lscaled, pdf);
            beta = spect_mul(beta, pscaled);

            // do absorption if we are hitting from inside the object
            if (isInside && surfelY.mat.type == GLASS) {
                beta = spect_mul(beta, spect_mul(surfelY.mat.albedo, exp(-0.01*surfelY.t)));
            }

            // Update the term that tracks radiance scaling for refraction
            // depending on whether the ray is entering or leaving the
            // medium.
            etaScale *= (dot(wo, surfelY.shadingNormal) > 0.0) ? (eta_for_RR * eta_for_RR) : 1.0 / (eta_for_RR * eta_for_RR);

            // Possibly terminate the path with Russian roulette.
            // Factor out radiance scaling due to refraction in rrBeta.
            Spectrum rrBeta = spect_mul(beta, etaScale);
            float maxRR = spect_max(rrBeta);
            if (maxRR < 1.0 && i > 3) {
                float q = max(0.05, 1.0 - maxRR);
                if (rand() < q) {
                    break;
                }
                beta = spect_div(beta, 1.0 - q);
            }
        } else {
            L = spect_add(L, spect_mul(beta, skyColor(wi)));
            break;
        }
    }

    return L;
}
`
