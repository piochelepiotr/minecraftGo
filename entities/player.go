package entities

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	//topSpeed in m/s
	topSpeed             float32 = 5
	topFlightVSpeed      float32 = 15
	topVSpeed            float32 = 20
	jumpHeight           float32 = 1.1
	g                            = 40
	breakingAcceleration         = 100
	minJumpPeriod = time.Millisecond * 500
)

//jumpSpeed is the vertical speed when jumping
var jumpSpeed = float32(math.Sqrt(float64(2* g * jumpHeight)))

// acceleration is the player acceleration in m/s2
var acceleration = mgl32.Vec3{10, -g, 10}

//Player is the player of the game
type Player struct {
	previousJump time.Time
	Entity   Entity
	Speed    mgl32.Vec3
	LastMove time.Time
	InFlight bool
}

func (p *Player) HSpeed() mgl32.Vec3 {
	return mgl32.Vec3{p.Speed.X(), 0, p.Speed.Z()}
}

func limitHSpeed(speed mgl32.Vec3) mgl32.Vec3 {
	v2 := math.Pow(float64(speed.X()), 2) + math.Pow(float64(speed.Z()), 2)
	x := speed.X()
	z := speed.Z()
	if v2 > 0 {
		v := float32(math.Sqrt(v2))
		if v > topSpeed {
			d := topSpeed / v
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

func (p *Player) FacingDir(dist float32) mgl32.Vec3 {
	mat := mgl32.HomogRotate3DY(p.Entity.Rotation.Y())
	forward := mgl32.Vec4{0, 0, -dist, 1}
	forward = mat.Mul4x1(forward)
	return forward.Vec3()
}

func (p *Player) SideFacingDir(dist float32) mgl32.Vec3 {
	mat := mgl32.HomogRotate3DY(p.Entity.Rotation.Y())
	side := mgl32.Vec4{-dist, 0, 0, 1}
	side = mat.Mul4x1(side)
	return side.Vec3()
}

// Accelerate increases the speed forward
func (p *Player) Accelerate(dist float32) {
	p.Speed = limitHSpeed(p.Speed.Add(p.FacingDir(dist)))
}

// Accelerate increases the speed forward
func (p *Player) LeftAccelerate(dist float32) {
	p.Speed = limitHSpeed(p.Speed.Add(p.SideFacingDir(dist)))
}

// speed in m/s, acc in m/s2, t in s
func friction(x float32, z float32, t float32) float32 {
	v2 := math.Pow(float64(x), 2) + math.Pow(float64(z), 2)
	speed := float32(math.Sqrt(v2))
	speed -= breakingAcceleration * t
	if speed <= 0 {
		return 0
	}
	return float32(math.Sqrt(math.Pow(float64(speed), 2) / v2))
}

// Friction returns new speed after friction has been applied
func Friction(speed mgl32.Vec3, t float32) mgl32.Vec3 {
	d := friction(speed.X(), speed.Z(), t)
	return mgl32.Vec3{
		speed.X() * d,
		speed.Y(),
		speed.Z() * d,
	}
}

// Gravity returns new speed after gravity has been applied
//speed in m/s, acc in m/s2, t in s
func Gravity(speed mgl32.Vec3, t float32, touchGround bool) mgl32.Vec3 {
	vSpeed := float32(speed.Y())
	if !touchGround {
		vSpeed = speed.Y() + acceleration.Y()*t
		if vSpeed < -topVSpeed {
			vSpeed = -topVSpeed
		}
	}
	return mgl32.Vec3{
		speed.X(),
		vSpeed,
		speed.Z(),
	}
}

// Forward returns the forward direction of the player
func (p *Player) Forward() mgl32.Vec3 {
	return p.Speed
}

// Jump makes the player jump
func (p *Player) Jump() {
	now := time.Now()
	if now.Sub(p.previousJump) < minJumpPeriod {
		return
	}
	p.previousJump = now
	p.Speed = mgl32.Vec3{
		p.Speed.X(),
		jumpSpeed,
		p.Speed.Z(),
	}
}

//PosPlus returns a point a little bit forward of the player
func (p *Player) PosPlus(e float32) mgl32.Vec3 {
	return p.Entity.Position.Add(p.FacingDir(e))
}

//Move updates the player's speed according to pressed keys
func (p *Player) Move(forward, backward, up, down, onGround, right, left bool) {
	if forward {
		// this isn't good, it's not a function of time
		p.Accelerate(0.5)
		// fmt.Printf("speed is %f\n", p.Speed.Len())
	}
	if backward {
		p.Accelerate(-0.5)
	}
	if right {
		p.LeftAccelerate(-5)
	}
	if left {
		p.LeftAccelerate(5)
	}
	if p.InFlight {
		if up {
			p.Speed[1] += 0.5
			if p.Speed[1] > topFlightVSpeed {
				p.Speed[1] = topFlightVSpeed
			}
		} else if down {
			p.Speed[1] -= 0.5
			if p.Speed[1] < -topFlightVSpeed {
				p.Speed[1] = -topFlightVSpeed
			}
		} else {
			p.Speed[1] = 0
		}
	} else {
		if up && onGround {
			p.Jump()
		}
	}
}
