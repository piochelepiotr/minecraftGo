package models

type TexturedModel struct {
    RawModel    RawModel
    TextureID   uint32
    Transparent bool
}
