package font

import (
	"github.com/go-gl/mathgl/mgl32"
)

type GUIText struct {
	Width float32
	Height float32
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
	VCenterText   bool
}

func CreateGUIText(text string, fontSize float32, font *FontType, position mgl32.Vec2, maxLineLength float32, maxTextHeight float32, vCentered bool, color mgl32.Vec3) GUIText {
	return GUIText{
		TextString:    text,
		FontSize:      fontSize,
		Font:          font,
		Position:      position,
		MaxLineSize:   maxLineLength,
		MaxTextHeight: maxTextHeight,
		VCenterText:   vCentered,
		Colour: color,
	}
}

func (text *GUIText) GetLineHeight() float32 {
	return text.FontSize * LINE_HEIGHT
}
