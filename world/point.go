package world

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Point defines a point in 3D
type Point struct {
	X int
	Y int
	Z int
}

// DistanceTo computes the distance from Point to an other point in space
func (p *Point) DistanceTo(p2 mgl32.Vec3) float32 {
	x := math.Pow(float64(p2.X())-float64(p.X), 2)
	// y := math.Pow(float64(p2.Y())-float64(p.Y), 2)
	z := math.Pow(float64(p2.Z())-float64(p.Z), 2)
	return float32(math.Sqrt(x + z))
}

//func n2ToN(n1, n2 int) int {
//	return ((n1 + n2) * (n1 + n2 + 1)) / 2
//}
//
//func n3ToN(n1, n2, n3 int) int {
//	return n2ToN(n2ToN(n1, n2), n3)
//}
//
//func ztoN(z int) int {
//	if z >= 0 {
//		return 2 * z
//	}
//	return z*(-2) - 1
//}
//
////GetID returns a uniq ID referencing the point
//func (p *Point) GetID() int {
//	xn := ztoN(p.X)
//	yn := ztoN(p.Y)
//	zn := ztoN(p.Z)
//	//return n3ToN(xn, yn, zn)
//	f := 1000000
//	return f*f*xn + f*yn + zn
//}
