package models

import "github.com/go-gl/mathgl/mgl32"

type RawModel struct {
    VaoID uint32
    VertexCount int32
}

type OutlineModel struct {
    RawModel
    Position mgl32.Vec3
}