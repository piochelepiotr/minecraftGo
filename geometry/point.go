package geometry

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

// Point defines a point in 3D
type Point struct {
	X int
	Y int
	Z int
}

func (p Point) Equal(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y && p.Z == p2.Z
}

func (p Point) Print() {
	fmt.Printf("Point {x:%d,y:%d,z:%d}\n", p.X, p.Y, p.Z)
}

// DistanceTo computes the Distance from Point to an other point in space
func (p Point) DistanceTo(p2 mgl32.Vec3) float32 {
	x := math.Pow(float64(p2.X())-float64(p.X), 2)
	// y := math.Pow(float64(p2.Y())-float64(p.Y), 2)
	z := math.Pow(float64(p2.Z())-float64(p.Z), 2)
	return float32(math.Sqrt(x + z))
}

func (p Point) ToFloat32() mgl32.Vec3 {
	return mgl32.Vec3{float32(p.X), float32(p.Y), float32(p.Z)}
}

func (p Point) Add(x, y, z int) Point {
	return Point{p.X + x, p.Y + y, p.Z + z}
}

func (p Point) GetKey() string {
	return fmt.Sprintf("%d_%d_%d", p.X, p.Y, p.Z)
}

