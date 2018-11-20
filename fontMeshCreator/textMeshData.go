package fontMeshCreator

type TextMeshData struct {
	vertexPositions []float32
	textureCoords []float32
}


func (t *TextMeshData) GetVertexCount() int{
	return len(t.vertexPositions)/2
}

