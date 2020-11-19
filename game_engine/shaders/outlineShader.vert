#version 400 core

in vec3 position;

uniform mat4 transformationMatrix;
uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;

void main(void) {
    vec4 worldPosition = transformationMatrix * vec4(position, 1);
    gl_Position = projectionMatrix * viewMatrix * worldPosition;
}
