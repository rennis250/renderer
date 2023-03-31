#version 410 core

precision highp float;

in vec2 TexCoord;

out vec4 color;

uniform sampler2D tex;
uniform vec3 iResolution;

// from: https://www.shadertoy.com/view/4sGcRW

//note: range [-2;2]
//note: not normalized
float Mitchell1D(float x) {
  //const float B = 0.0; //Catmull-Rom?
  const float B = 1.0/3.0; //Mitchell
  const float C = 0.5 * (1.0-B);
  x = abs( 2.0 * x );
    if ( x > 2.0 )
      return 0.0;
  if (x > 1.0)
    return ((-B - 6.0*C) * x*x*x + (6.0*B + 30.0*C) * x*x + (-12.0*B - 48.0*C) * x + (8.0*B + 24.0*C)) * (1.0/6.0);
  else
    return ((12.0 - 9.0*B - 6.0*C) * x*x*x + (-18.0 + 12.0*B + 6.0*C) * x*x + (6.0 - 2.0*B)) * (1.0/6.0);
}

float FilterMitchell(vec2 p, vec2 r) {
  p /= r; //TODO: fails at radius0
  return Mitchell1D(length(p));
}

vec3 blur(sampler2D smpl, vec2 p, float filtersiz_nm_ch0, vec2 ch0siz) {
  float filtersiz = 1.6;

  vec2 pq = (floor(p*ch0siz) + vec2(0.5, 0.5)) / ch0siz;
  
  vec4 bb_nm = vec4(pq - vec2(filtersiz_nm_ch0),
                    pq + vec2(filtersiz_nm_ch0));
  vec4 bb_px_q = vec4(floor(bb_nm.xy * ch0siz.xy), ceil(bb_nm.zw * ch0siz.xy));
  vec4 bb_nm_q = bb_px_q / ch0siz.xyxy;
  ivec2 bb_px_siz = ivec2(bb_px_q.zw - bb_px_q.xy);

  vec3 sumc = vec3(0.0);
  float sumw = 0.0;
  for (int y = 0; y < bb_px_siz.y; ++y) {
      for (int x = 0; x < bb_px_siz.x; ++x) {
          vec2 xy_f = (vec2(x,y) + vec2(0.5)) / vec2(bb_px_siz);
          vec2 sp = bb_nm_q.xy + (bb_nm_q.zw - bb_nm_q.xy)*xy_f;

          float w = FilterMitchell(sp - p, vec2(filtersiz_nm_ch0));
          
          sumc += w*texture(tex, sp, -10.0).rgb;
          
          sumw += w;
      }
  }

  return sumc / sumw;
}

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
}