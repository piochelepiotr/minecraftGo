package world

// Point defines a point in 3D
type Point struct {
	X int
	Y int
	Z int
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
