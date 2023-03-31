package renderer

const fxaa_header_glsl = `
#version 410 core

precision highp float;

in vec2 TexCoord;

out vec4 color;

uniform sampler2D tex;
uniform vec3 iResolution;
            
#define FXAA_REDUCE_MIN   (1.0/ 128.0)
#define FXAA_REDUCE_MUL   (1.0 / 8.0)
#define FXAA_SPAN_MAX     8.0`
