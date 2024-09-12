package rendering

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type Camera struct {
	Position    mgl64.Vec3
	Front       mgl64.Vec3
	Up          mgl64.Vec3
	Right       mgl64.Vec3
	WorldUp     mgl64.Vec3
	Yaw         float64
	Pitch       float64
	Speed       float64
	Sensitivity float64
	Fov         float32
}

func (c *Camera) UpdateDirection(dx, dy float64) {
	c.Pitch += dy
	if c.Pitch > 89 {
		c.Pitch = 89
	} else if c.Pitch < -89 {
		c.Pitch = -89
	}
	c.Yaw = math.Mod(c.Yaw+dx, 360)
	c.UpdateVec()
}

func (c *Camera) UpdateVec() {
	c.Front = mgl64.Vec3{
		math.Cos(mgl64.DegToRad(c.Pitch)) * math.Cos(mgl64.DegToRad(c.Yaw)),
		math.Sin(mgl64.DegToRad(c.Pitch)),
		math.Cos(mgl64.DegToRad(c.Pitch)) * math.Sin(mgl64.DegToRad(c.Yaw)),
	}.Normalize()
	c.Right = c.WorldUp.Cross(c.Front).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

func (camera *Camera) GetTransform() mgl32.Mat4 {
	cameraTarget := camera.Position.Add(camera.Front)
	return mgl32.LookAt(
		float32(camera.Position.X()), float32(camera.Position.Y()), float32(camera.Position.Z()),
		float32(cameraTarget.X()), float32(cameraTarget.Y()), float32(cameraTarget.Z()),
		float32(camera.Up.X()), float32(camera.Up.Y()), float32(camera.Up.Z()),
	)
}

func (c *Camera) GetFov() float32 {
	return c.Fov
}
