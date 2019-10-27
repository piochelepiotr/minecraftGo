package toolbox

import (
	"github.com/go-gl/mathgl/mgl32"
)

func CreateTransformationMatrix2D(translation, scale mgl32.Vec2) mgl32.Mat4 {
	m := mgl32.Ident4()
	translate := mgl32.Translate3D(translation.X(), -translation.Y(), 0)
	scaleM := mgl32.Scale3D(scale.X(), scale.Y(), 1)
	return translate.Mul4(scaleM.Mul4(m))
}

func CreateTransformationMatrix(translation mgl32.Vec3, rotation mgl32.Vec3, scale float32) mgl32.Mat4 {
	m := mgl32.Ident4()
	rotX := mgl32.HomogRotate3DX(rotation.X())
	rotY := mgl32.HomogRotate3DY(rotation.Y())
	rotZ := mgl32.HomogRotate3DZ(rotation.Z())
	translate := mgl32.Translate3D(translation.X(), translation.Y(), translation.Z())
	scaleM := mgl32.Scale3D(scale, scale, scale)
	return translate.Mul4(scaleM.Mul4(rotZ.Mul4(rotY.Mul4(rotX.Mul4(m)))))
}

func CreateViewMatrix(translation mgl32.Vec3, rotation mgl32.Vec3) mgl32.Mat4 {
	scale := float32(1)
	m := mgl32.Ident4()
	rotX := mgl32.HomogRotate3DX(rotation.X())
	rotY := mgl32.HomogRotate3DY(rotation.Y())
	rotZ := mgl32.HomogRotate3DZ(rotation.Z())
	translate := mgl32.Translate3D(-translation.X(), -translation.Y(), -translation.Z())
	scaleM := mgl32.Scale3D(scale, scale, scale)
	return rotX.Mul4(rotY.Mul4(rotZ.Mul4(scaleM.Mul4(translate.Mul4(m)))))
}

func Abs(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}
