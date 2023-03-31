package renderer

const surfel_glsl = `
struct Surfel {
    float t;
    vec3 position;
    vec3 shadingNormal;
    Material mat;
    bool hit;
};

Spectrum finiteScatteringDensity(inout Surfel surfelX, inout vec3 wi, const vec3 woW, out float pdf, out float eta_for_RR) {
    Material mat = surfelX.mat;
    vec3 X = surfelX.position;
    vec3 n = surfelX.shadingNormal;

    mat3 ltow = ONB(n);
    mat3 wtol = transpose(ltow);

    if (mat.type == LAMB) {
        wi = ltow * cosWeightHemi();
        if (dot(wi, n) > 0.0 && dot(woW, n) > 0.0) {
            surfelX.position += n*EPS;
            pdf = abs(dot(wi, n)) * InvPi;
            return spect_mul(mat.albedo, InvPi);
        } else {
            pdf = 0.0;
            return ZEROSPECT;
        }
    } else if (mat.type == METAL) {
        surfelX.position += n * EPS;
        wi = reflect(-woW, n);
        pdf = 1.0;
        return spect_mul(mat.albedo, FrConductor(abs(dot(wi, n)), 1.0, 1.4, 1.6) / abs(dot(wi, n)));
    } else if (mat.type == GLASS) {
        float cos_theta = dot(woW, n);
        float extIor = 1.0;
        float intIor = 1.5;
        
        float eta;
        vec3 outward_normal;
        if (cos_theta < 0.0) {
            outward_normal = -n;
            eta = intIor / extIor;
        } else {
            outward_normal = n;
            eta = extIor / intIor;
        }
        eta_for_RR = eta;

        float F = FrDielectric(cos_theta, extIor, intIor);
        if (rand() <= F) {
            wi = reflect(-woW, n);
            pdf = F;
            return spect_mul(UNITY, F / abs(dot(wi, n)));
        } else {
            wi = refract(-woW, outward_normal, eta);
            pdf = 1.0 - F;
            return spect_mul(UNITY, (eta*eta) * (1.0 - F) / abs(dot(wi, n)));
        }
    }
}
`
