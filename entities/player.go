package entities

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	//TopSpeed in m/s
	TopSpeed float32 = 20
)

var acceleration = mgl32.Vec3{5, 0, 5}

//Player is the player of the game
type Player struct {
	Entity   Entity
	speed    mgl32.Vec3
	lastMove int64
}

func limitSpeed(speed mgl32.Vec3) mgl32.Vec3 {
	v2 := math.Pow(float64(speed.X()), 2) + math.Pow(float64(speed.Z()), 2)
	x := speed.X()
	z := speed.Z()
	if v2 > 0 {
		v := float32(math.Sqrt(v2))
		if v > TopSpeed {
			d := float32(math.Sqrt(v2 / float64(TopSpeed)))
			x = x / d
			z = z / d
		}
	}
	return mgl32.Vec3{
		x,
		speed.Y(),
		z,
	}
}

// Accelerate increases the speed forward
func (p *Player) Accelerate(dist float32) {
	mat := mgl32.HomogRotate3DY(p.Entity.Rotation.Y())
	forward := mgl32.Vec4{0, 0, -dist, 1}
	forward = mat.Mul4x1(forward)
	//p.Entity.Position = p.Entity.Position.Add(forward.Vec3())
	p.speed = limitSpeed(p.speed.Add(forward.Vec3()))
}

// speed in m/s, acc in m/s2, t in s
func friction(x float32, z float32, acc float32, t float32) float32 {
	v2 := math.Pow(float64(x), 2) + math.Pow(float64(z), 2)
	speed := float32(math.Sqrt(v2))
	speed -= acc * t
	if speed <= 0 {
		return 0
	}
	return float32(math.Sqrt(math.Pow(float64(speed), 2) / v2))
}

// speed in m/s, acc in m/s2, t in s
func gravity(speed float32, acc float32, t float32) float32 {
	speed += acc * t
	if speed < -TopSpeed {
		speed = -TopSpeed
	}
	return speed
}

func forces(speed mgl32.Vec3, t float32) mgl32.Vec3 {
	d := friction(speed.X(), speed.Z(), acceleration.X(), t)
	return mgl32.Vec3{
		speed.X() * d,
		gravity(speed.Y(), acceleration.Y(), t),
		speed.Z() * d,
	}
}

// Move moves player at its speed
func (p *Player) Move(moves bool) {
	now := time.Now().UnixNano()
	if p.lastMove != 0 {
		diff := now - p.lastMove
		secDiff := float32(diff) / 1e9
		p.Entity.Position = p.Entity.Position.Add(
			mgl32.Vec3{
				p.speed.X() * secDiff,
				p.speed.Y() * secDiff,
				p.speed.Z() * secDiff,
			})
		if !moves {
			p.speed = forces(p.speed, secDiff)
		}
	}
	p.lastMove = now
}
