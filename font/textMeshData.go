package font

type TextMeshData struct {
	VertexPositions []float32
	TextureCoords   []float32
}

func (t *TextMeshData) GetVertexCount() int32 {
	return int32(len(t.VertexPositions) / 2)
}
