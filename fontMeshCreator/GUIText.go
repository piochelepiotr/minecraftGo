package fontMeshCreator

import (
	"github.com/go-gl/mathgl/mgl32"
)

type GUIText struct {
	TextString    string
	FontSize      float32
	TextMeshVao   uint32
	VertexCount   int32
	Colour        mgl32.Vec3
	Position      mgl32.Vec2
	MaxLineSize   float32
	MaxTextHeight float32
	NumberOfLines int
	Font          *FontType
	CenterText    bool
	VCenterText   bool
}

func CreateGUIText(text string, fontSize float32, font *FontType, position mgl32.Vec2, maxLineLength float32, centered bool, maxTextHeight float32, vCentered bool) GUIText {
	return GUIText{
		TextString:    text,
		FontSize:      fontSize,
		Font:          font,
		Position:      position,
		MaxLineSize:   maxLineLength,
		CenterText:    centered,
		MaxTextHeight: maxTextHeight,
		VCenterText:   vCentered,
	}
}

func (text *GUIText) GetLineHeight() float32 {
	return text.FontSize * LINE_HEIGHT
}
