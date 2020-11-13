#version 400 core

in vec2 pass_textureCoords;
in vec3 surfaceNormal;
in vec3 pass_colors;

out vec4 out_Color;

uniform sampler2D textureSampler;

void main(void) {
    vec3 unitNormal = normalize(surfaceNormal);
    float l = 0.8;
    if(dot(unitNormal, vec3(0,1,0)) > 0.8)
    {
        l = 1;
    }
    else if(dot(unitNormal, vec3(0,-1,0)) > 0.8)
    {
        l = 0.4;
    }

    vec4 textureColor = texture(textureSampler, pass_textureCoords);
    if(textureColor.a < 0.5) {
        discard;
    }

    out_Color = l * textureColor * vec4(pass_colors, 1.0);
}
