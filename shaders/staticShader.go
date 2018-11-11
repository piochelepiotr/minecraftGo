package shaders

func CreateStaticShader() ShaderProgram {
    s:= CreateShader("shaders/vertexShader.txt", "shaders/fragmentShader.txt")
    s.bindAttribute(0, "position\x00")
    s.bindAttribute(1, "textureCoords\x00")
    return s
}
