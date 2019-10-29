package world

// Block is an id representing a type of Block
type Block uint8

const (
	// Dirt block
	Dirt Block = 2
	// Stone block
	Stone Block = 1
	// Grass block
	Grass Block = 0
	// GrassSide block
	GrassSide Block = 3
	//Air block
	Air Block = 255
)

// blockFaces allows you to put different blocks on sides, top and
// bottom of blocks
var blockFaces = map[Block]map[Face]Block{
	Grass: {
		Side:   GrassSide,
		Bottom: Dirt,
		Top:    Grass,
	},
}

type Face byte
const (
	Top Face = 0
	Side Face = 1
	Bottom Face = 2
)

func (b Block) GetSide(f Face) Block {
	if sides, ok := blockFaces[b]; ok {
		if side, ok := sides[f]; ok {
			return side
		}
	}
	return b
}
