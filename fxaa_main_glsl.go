package renderer

const fxaa_main_glsl = `
vec4 mainImage(in vec2 fragCoord) {
    vec4 outcol = apply(tex, gl_FragCoord.xy, iResolution.xy);
    // vec4 outcol = texture(tex, uv);
    return vec4(sqrt(outcol.rgb), 1.0);
}

void main() {
    color = mainImage(TexCoord);
}`
