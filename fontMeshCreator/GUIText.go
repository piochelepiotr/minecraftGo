package fontMeshCreator;

import (
	"github.com/go-gl/mathgl/mgl32"
)

type GUIText struct {
	TextString string
	FontSize float32
	TextMeshVao int
	VertexCount int
	Colour mgl32.Vec3
	Position mgl32.Vec2
	LineMaxSize float32
	NumberOfLines int
	Font FontType
	CenterText bool
}

func CreateGUIText(text string, fontSize float32, font FontType, position mgl32.Vec2, maxLineLength float32, centered bool) GUIText {
    return GUIText{
        TextString : text,
        FontSize : fontSize,
        Font : font,
        Position : position,
        LineMaxSize : maxLineLength,
        CenterText : centered,
    }
}

