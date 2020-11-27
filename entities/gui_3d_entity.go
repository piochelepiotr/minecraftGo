package entities

import "github.com/go-gl/mathgl/mgl32"

type Gui3dEntity struct {
	Entity Entity
	Translation mgl32.Vec2
	Scale float32
}