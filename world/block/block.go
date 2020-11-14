package block

import "github.com/go-gl/mathgl/mgl32"

// Block is an id representing a type of Block
type Block uint8

const (
	Dirt Block = 2
	Stone Block = 1
	Grass Block = 0
	GrassSide Block = 3
	Tree Block = 21
	TreeSide Block = 20
	Leaves Block = 52
	Cactus Block = 69
	CactusSide Block = 70
	CactusBottom Block = 71
	Sand Block = 18
	Air Block = 255
)

// blockFaces allows you to put different blocks on sides, top and
// bottom of blocks
var blockFaces = map[Block]map[Face]Block{
	Grass: {
		Side:   GrassSide,
		Bottom: Dirt,
	},
	Tree: {
		Side:   TreeSide,
	},
	Cactus: {
		Side: CactusSide,
		Bottom: CactusBottom,
	},
}

var blockColors = map[Block]mgl32.Vec3{
	Leaves: {55.0/255.0, 97.0/255.0, 43.0/255.0},
}

type Face byte
const (
	Top Face = 0
	Side Face = 1
	Bottom Face = 2
)

func (b Block) IsTransparent() bool {
	return b == Air || b == Leaves
}

func (b Block) GetSide(f Face) Block {
	if sides, ok := blockFaces[b]; ok {
		if side, ok := sides[f]; ok {
			return side
		}
	}
	return b
}

func (b Block) GetColor() mgl32.Vec3 {
	if color, ok := blockColors[b]; ok {
		return color
	}
	return mgl32.Vec3{1, 1, 1}
}
