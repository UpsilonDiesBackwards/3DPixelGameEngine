#version 420
#extension GL_ARB_explicit_uniform_location : enable

layout (location = 0) out vec4 frag_colour;

layout(location = 0) in vec2 TexCoord;

uniform sampler2D texture1;

void main() {
    frag_colour = texture(texture1, TexCoord);
}
