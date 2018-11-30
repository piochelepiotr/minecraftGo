package entities

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	//TopSpeed in m/s
	TopSpeed float32 = 20
	//JumpSpeed is the vertical speed when jumping
	JumpSpeed float32 = 5
)

// Acceleration is the player acceleration in m/s2
var acceleration = mgl32.Vec3{10, -10, 10}

//Player is the player of the game
type Player struct {
	Entity   Entity
	Speed    mgl32.Vec3
	LastMove int64
}

func limitSpeed(speed mgl32.Vec3) mgl32.Vec3 {
	v2 := math.Pow(float64(speed.X()), 2) + math.Pow(float64(speed.Z()), 2)
	x := speed.X()
	z := speed.Z()
	if v2 > 0 {
		v := float32(math.Sqrt(v2))
		if v > TopSpeed {
			d := TopSpeed / v
			x = x * d
			z = z * d
		}
	}
	return mgl32.Vec3{
		x,
		speed.Y(),
		z,
	}
}

func (p *Player) facingDir(dist float32) mgl32.Vec3 {
	mat := mgl32.HomogRotate3DY(p.Entity.Rotation.Y())
	forward := mgl32.Vec4{0, 0, -dist, 1}
	forward = mat.Mul4x1(forward)
	return forward.Vec3()
}

// Accelerate increases the speed forward
func (p *Player) Accelerate(dist float32) {
	p.Speed = limitSpeed(p.Speed.Add(p.facingDir(dist)))
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

// Forces returns the new speed after applying forces to player
func Forces(speed mgl32.Vec3, t float32, touchGround bool) mgl32.Vec3 {
	d := friction(speed.X(), speed.Z(), acceleration.X(), t)
	ySpeed := float32(0)
	if !touchGround {
		ySpeed = gravity(speed.Y(), acceleration.Y(), t)
	}
	return mgl32.Vec3{
		speed.X() * d,
		ySpeed,
		speed.Z() * d,
	}
}

// Forward returns the forward direction of the player
func (p *Player) Forward() mgl32.Vec3 {
	return p.Speed
}

// Jump makes the player jump
func (p *Player) Jump() {
	p.Speed = mgl32.Vec3{
		p.Speed.X(),
		JumpSpeed,
		p.Speed.Z(),
	}
}

//PosPlus returns a point a little bit forward of the player
func (p *Player) PosPlus(e float32) mgl32.Vec3 {
	return p.Entity.Position.Add(p.facingDir(e))
}

//Move updates the player's speed according to pressed keys
func (p *Player) Move(forward, backward, jump bool) {
	if forward {
		p.Accelerate(0.5)
	}
	if backward {
		p.Accelerate(-0.5)
	}
	if jump {
		p.Jump()
	}
}
