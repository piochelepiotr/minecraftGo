package toolbox

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func distance(x, y, z, x2, y2, z2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x-x2), 2) + math.Pow(float64(y-y2), 2) + math.Pow(float64(z-z2), 2)))
}

//GetNextBlock returns the next block in the direction dir
func GetNextBlock(xIn, yIn, zIn *float32, dir mgl32.Vec3, x, y, z *int) float32 {
	xS := *xIn
	yS := *yIn
	zS := *zIn
	xB := float32(*x)
	yB := float32(*y)
	zB := float32(*z)
	xE := xB + 1
	yE := yB + 1
	zE := zB + 1
	//up face
	if dir.Y() > 0 {
		*yIn = yE
		t := (*yIn - yS) / dir.Y()
		*xIn = xS + t*dir.X()
		*zIn = zS + t*dir.Z()
		if *zIn >= zB && *zIn <= zE && *xIn >= xB && *xIn <= xE {
			//fmt.Println("BACK")
			(*y)++
			return distance(*xIn, *yIn, *zIn, xS, yS, zS)
		}
	}
	//bottom face
	if dir.Y() < 0 {
		*yIn = yB
		t := (*yIn - yS) / dir.Y()
		*xIn = xS + t*dir.X()
		*zIn = zS + t*dir.Z()
		if *zIn >= zB && *zIn <= zE && *xIn >= xB && *xIn <= xE {
			//fmt.Println("BOTTOM")
			(*y)--
			return distance(*xIn, *yIn, *zIn, xS, yS, zS)
		}
	}
	//right face
	if dir.X() > 0 {
		*xIn = xE
		t := (*xIn - xS) / dir.X()
		*yIn = yS + t*dir.Y()
		*zIn = zS + t*dir.Z()
		if *zIn >= zB && *zIn <= zE && *yIn >= yB && *yIn <= yE {
			//fmt.Println("RIGHT")
			(*x)++
			return distance(*xIn, *yIn, *zIn, xS, yS, zS)
		}
	}
	//left face
	if dir.X() < 0 {
		*xIn = xB
		t := (*xIn - xS) / dir.X()
		*yIn = yS + t*dir.Y()
		*zIn = zS + t*dir.Z()
		if *zIn >= zB && *zIn <= zE && *yIn >= yB && *yIn <= yE {
			//fmt.Println("LEFT")
			(*x)--
			return distance(*xIn, *yIn, *zIn, xS, yS, zS)
		}
	}
	//front face
	if dir.Z() > 0 {
		*zIn = zE
		t := (*zIn - zS) / dir.Z()
		*xIn = xS + t*dir.X()
		*yIn = yS + t*dir.Y()
		if *xIn >= xB && *xIn <= xE && *yIn >= yB && *yIn <= yE {
			//fmt.Println("FRONT")
			(*z)++
			return distance(*xIn, *yIn, *zIn, xS, yS, zS)
		}
	}
	//back face
	if dir.Z() < 0 {
		*zIn = zB
		t := (*zIn - zS) / dir.Z()
		*xIn = xS + t*dir.X()
		*yIn = yS + t*dir.Y()
		if *xIn >= xB && *xIn <= xE && *yIn >= yB && *yIn <= yE {
			//fmt.Println("BACK")
			(*z)--
			return distance(*xIn, *yIn, *zIn, xS, yS, zS)
		}
	}
	//fmt.Println("problem : no face to get out of the cube")
	return 0
}
