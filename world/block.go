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

// BlockSides allows you to put different blocks on sides, top and
// bottom of blocks
var BlockSides = map[Block]map[string]Block{
	Grass: {
		"side":   GrassSide,
		"bottom": Dirt,
		"top":    Grass,
	},
	Dirt: {
		"side":   Dirt,
		"bottom": Dirt,
		"top":    Dirt,
	},
}
