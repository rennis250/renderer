package renderer

const filt_main_glsl = `
vec4 mainImage(in vec2 fragCoord) {
  // vec2 uv = fragCoord.xy / iResolution.xy;
  vec2 uv = TexCoord;
  // float radius = 0.2;
  // radius = pow(radius, 2.0);
  
  // vec2 ch0siz = iChannelResolution[0].xy; //TODO: input RT and cutout from that
  vec2 ch0siz = iResolution.xy;
  float filtersiz = 1.6;

  //note: radius
  const float MIN_BLUR_SIZ_PX = 1.0;
  const float MAX_BLUR_SIZ_PX = 5.0;

  float footprintsiz_t = 0.1;
  float footprintsiz_px = mix(MIN_BLUR_SIZ_PX, MAX_BLUR_SIZ_PX, footprintsiz_t);
  float footprintsiz_nm_ch0 = footprintsiz_px / ch0siz.y;
  
  vec2 uv_q = (floor(uv * ch0siz) + vec2(0.5)) / ch0siz;
  vec3 col = blur(tex, uv_q, footprintsiz_nm_ch0*filtersiz, ch0siz);
  return vec4(col, 1.0);
}

void main() {
  color = mainImage(gl_FragCoord.xy);
}`
