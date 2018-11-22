package fontMeshCreator

type FontType struct {
	TextureAtlas uint32
	Loader       TextMeshCreator
}

func CreateFontType(textureAtlas uint32, fontFile string, aspectRatio float32) FontType {
	return FontType{
		TextureAtlas: textureAtlas,
		Loader:       CreateTextMeshCreator(fontFile, aspectRatio),
	}
}

func (ft *FontType) LoadText(text GUIText) TextMeshData {
	return ft.Loader.createTextMesh(text)
}
