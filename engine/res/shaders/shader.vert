#version 420
#extension GL_ARB_explicit_uniform_location : enable
#extension GL_ARB_enhanced_layouts : enable

layout(location = 0) in vec3 vp;
layout(location = 1) in vec2 aTexCoord;
layout(location = 2) in vec3 aNormal;

layout(location = 0) out vec2 TexCoord;

layout(binding = 1) uniform PerspectiveBlock {
    mat4 project;
    mat4 camera;
    mat4 model;
};

void main() {
    gl_Position = project * camera * model * vec4(vp, 1);
    TexCoord = aTexCoord;
}