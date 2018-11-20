package fontMeshCreator;

type FontType struct {
	TextureAtlas uint32
	Loader TextMeshCreator
}

func CreateFontType(textureAtlas uint32, fontFile string) FontType {
    return FontType {
        TextureAtlas : textureAtlas,
        Loader : CreateTextMeshCreator(fontFile),
    }
}

func (ft *FontType)  LoadText(text string) TextMeshData {
	return ft.Loader.createTextMesh(text);
}


