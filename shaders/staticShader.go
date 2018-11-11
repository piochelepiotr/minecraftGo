package shaders

func CreateStaticShader() ShaderProgram {
    bindAttributes := func(s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
        s.bindAttribute(1, "textureCoords\x00")
    }
    s:= CreateShader("shaders/vertexShader.txt", "shaders/fragmentShader.txt", bindAttributes)
    return s
}

