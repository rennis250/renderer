package renderer

const header_uniforms_glsl = `
#version 410 core

out vec4 fragColor;

in vec2 TexCoord;

uniform sampler2D blue_noise_tex;
uniform sampler2D voronoi_indexed_tex;

uniform vec2 iResolution;
uniform float iTime;
uniform float iFrame;

uniform float ld;
uniform float rg;
uniform float by;

uniform int sphere_color;
uniform int with_aperture;
`
