#version 330

in vec2 position;
in vec2 textureCoords;

out vec2 pass_textureCoords;

uniform vec2 translation;

void main(void){
    gl_Position = vec4(position + vec2(translation.x, translation.y) * vec2(1,-1), 0.0, 1.0);
    pass_textureCoords = textureCoords;
}
