package renderer

const filt_header_glsl = `
#version 410 core

precision highp float;

in vec2 TexCoord;

out vec4 color;

uniform sampler2D tex;
uniform vec3 iResolution;

// from: https://www.shadertoy.com/view/4sGcRW
`
