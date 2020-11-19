#version 400 core

in vec3 position;
in vec2 textureCoords;
in vec3 normal;
in vec3 colors;

out vec2 pass_textureCoords;
out vec3 pass_colors;
out vec3 surfaceNormal;

uniform mat4 transformationMatrix;
uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;

void main(void) {
    vec4 worldPosition = transformationMatrix * vec4(position, 1);
    gl_Position = projectionMatrix * viewMatrix * worldPosition;
    pass_textureCoords = textureCoords;
    pass_colors = colors;
    surfaceNormal = (transformationMatrix * vec4(normal, 0)).xyz;
}
