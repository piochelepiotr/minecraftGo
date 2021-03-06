package entities

import (
	"github.com/go-gl/mathgl/mgl32"
)

//Camera is the point of view of the game
type Camera struct {
	Position       mgl32.Vec3
	Rotation       mgl32.Vec3
	FollowDistance float32
	Height         float32
}

//CreateCamera creates a camera
func CreateCamera(x, y, z, followDistance, height float32) *Camera {
	return &Camera{
		Position:       mgl32.Vec3{x, y, z},
		Rotation:       mgl32.Vec3{0, 0, 0},
		FollowDistance: followDistance,
		Height:         height,
	}
}

//LockOnPlayer locks the camera on the player
func (c *Camera) LockOnPlayer(player *Player) {
	rotY := mgl32.Rotate3DY(player.Entity.Rotation.Y())
	cameraShift := mgl32.Vec3{0, 0, c.FollowDistance}
	movement := rotY.Mul3x1(cameraShift)
	c.Position = mgl32.Vec3{
		player.Entity.Position.X() + movement.X(),
		player.Entity.Position.Y() + c.Height,
		player.Entity.Position.Z() + movement.Z(),
	}
	c.Rotation = mgl32.Vec3{
		c.Rotation.X(),
		-player.Entity.Rotation.Y(),
		c.Rotation.Z(),
	}
}

//IncreaseRotation increases the rotation of the camera
func (c *Camera) IncreaseRotation(dx, dy, dz float32) {
	c.Rotation = c.Rotation.Add(mgl32.Vec3{dx, dy, dz})
}

//IncreasePosition increases the position of the camera
func (c *Camera) IncreasePosition(dx, dy, dz float32) {
	c.Position = c.Position.Add(mgl32.Vec3{dx, dy, dz})
}

// CursorVector returns the vector that represents what the player is looking at
func (c *Camera) CursorVector() mgl32.Vec3 {
	return mgl32.Vec3{0, 0, 0}
}
