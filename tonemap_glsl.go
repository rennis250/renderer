package renderer

const tonemap_glsl = `
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
`
