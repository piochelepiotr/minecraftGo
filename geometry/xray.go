package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func distance(x, y, z, x2, y2, z2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x-x2), 2) + math.Pow(float64(y-y2), 2) + math.Pow(float64(z-z2), 2)))
}

type Xray struct {
	P        Point
	Previous Point
	dir      mgl32.Vec3
	in       mgl32.Vec3
	Distance float32
}

func NewXray(start mgl32.Vec3, dir mgl32.Vec3) *Xray {
	p := Point{
		int(math.Floor(float64(start.X()))),
		int(math.Floor(float64(start.Y()))),
		int(math.Floor(float64(start.Z()))),
	}
	return &Xray{
		P: p,
		Previous: p,
		dir: dir,
		in: start,
	}
}

func isInside(p mgl32.Vec3, blockStart mgl32.Vec3, blockEnd mgl32.Vec3) bool {
	return p.X() >= blockStart.X() && p.X() <= blockEnd.X() &&
		p.Y() >= blockStart.Y() && p.Y() <= blockEnd.Y() &&
		p.Z() >= blockStart.Z() && p.Z() <= blockEnd.Z()
}

// GoToNextBlock goes to the next block in the direction dir
func (x *Xray) GoToNextBlock() {
	x.Previous = x.P
	start := x.in
	blockStart := x.P.ToFloat32()
	blockEnd := blockStart.Add(mgl32.Vec3{1, 1, 1})
	//up face
	if x.dir.Y() > 0 {
		t := (blockEnd.Y() - start.Y()) / x.dir.Y()
		newIn := mgl32.Vec3{start.X() + t*x.dir.X(), blockEnd.Y(), start.Z() + t*x.dir.Z()}
		if isInside(newIn, blockStart, blockEnd) {
			//fmt.Println("BACK")
			x.in = newIn
			x.P = x.P.Add(0, 1, 0)
			x.Distance += distance(x.in.X(), x.in.Y(), x.in.Z(), start.X(), start.Y(), start.Z())
		}
	}
	//bottom face
	if x.dir.Y() < 0 {
		t := (blockStart.Y() - start.Y()) / x.dir.Y()
		newIn := mgl32.Vec3{start.X() + t*x.dir.X(), blockStart.Y(), start.Z() + t*x.dir.Z()}
		if isInside(newIn, blockStart, blockEnd) {
			//fmt.Println("BACK")
			x.in = newIn
			x.P = x.P.Add(0, -1, 0)
			x.Distance += distance(x.in.X(), x.in.Y(), x.in.Z(), start.X(), start.Y(), start.Z())
		}
	}
	//right face
	if x.dir.X() > 0 {
		t := (blockEnd.X() - start.X()) / x.dir.X()
		newIn := mgl32.Vec3{blockEnd.X(), start.Y() + t*x.dir.Y(), start.Z() + t*x.dir.Z()}
		if isInside(newIn, blockStart, blockEnd) {
			//fmt.Println("BACK")
			x.in = newIn
			x.P = x.P.Add(1, 0, 0)
			x.Distance += distance(x.in.X(), x.in.Y(), x.in.Z(), start.X(), start.Y(), start.Z())
		}
	}
	//left face
	if x.dir.X() < 0 {
		t := (blockStart.X() - start.X()) / x.dir.X()
		newIn := mgl32.Vec3{blockStart.X(), start.Y() + t*x.dir.Y(), start.Z() + t*x.dir.Z()}
		if isInside(newIn, blockStart, blockEnd) {
			//fmt.Println("BACK")
			x.in = newIn
			x.P = x.P.Add(-1, 0, 0)
			x.Distance += distance(x.in.X(), x.in.Y(), x.in.Z(), start.X(), start.Y(), start.Z())
		}
	}
	//front face
	if x.dir.Z() > 0 {
		t := (blockEnd.Z() - start.Z()) / x.dir.Z()
		newIn := mgl32.Vec3{start.X() + t*x.dir.X(), start.Y() + t*x.dir.Y(), blockEnd.Z()}
		if isInside(newIn, blockStart, blockEnd) {
			//fmt.Println("BACK")
			x.in = newIn
			x.P = x.P.Add(0, 0, 1)
			x.Distance += distance(x.in.X(), x.in.Y(), x.in.Z(), start.X(), start.Y(), start.Z())
		}
	}
	//back face
	if x.dir.Z() < 0 {
		t := (blockStart.Z() - start.Z()) / x.dir.Z()
		newIn := mgl32.Vec3{start.X() + t*x.dir.X(), start.Y() + t*x.dir.Y(), blockStart.Z()}
		if isInside(newIn, blockStart, blockEnd) {
			//fmt.Println("BACK")
			x.in = newIn
			x.P = x.P.Add(0, 0, -1)
			x.Distance += distance(x.in.X(), x.in.Y(), x.in.Z(), start.X(), start.Y(), start.Z())
		}
	}
	//fmt.Println("problem : no face to get out of the cube")
}

//ComputeCameraRay computes the camera ray
func ComputeCameraRay(cameraRotation mgl32.Vec3) mgl32.Vec3 {
	rotY := mgl32.HomogRotate3DY(cameraRotation.Y())
	rotX := mgl32.HomogRotate3DX(cameraRotation.X())
	xray := mgl32.Vec4{0, 0, -1, 1}
	xray = rotX.Mul4x1(xray)
	xray = rotY.Mul4x1(xray)
	return mgl32.Vec3{-xray.X(), -xray.Y(), xray.Z()}
}
