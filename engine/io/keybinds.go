package io

import (
	"3DPixelGameEngine/engine/rendering"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

var ViewportTransform mgl32.Mat4
var u = &UserInput{}
var camera = &rendering.Camera{
	Position:    mgl64.Vec3{0, 0, 0},
	Up:          mgl64.Vec3{0, 1, 0},
	WorldUp:     mgl64.Vec3{0, 1, 0},
	Yaw:         -90,
	Pitch:       0,
	Speed:       12,
	Sensitivity: 0.075,
	Fov:         60,
}

func InputRunner(win *rendering.Window, deltaTime float64) error {
	adjCSpeed := deltaTime * float64(camera.Speed)
	if ActionState[VP_FORW] {
		camera.Position = camera.Position.Add(camera.Front.Mul(adjCSpeed))
	}
	if ActionState[VP_BACK] {
		camera.Position = camera.Position.Sub(camera.Front.Mul(adjCSpeed))
	}
	if ActionState[VP_LEFT] {
		camera.Position = camera.Position.Sub(camera.Front.Cross(camera.Up).Mul(adjCSpeed))
	}
	if ActionState[VP_RGHT] {
		camera.Position = camera.Position.Add(camera.Front.Cross(camera.Up).Mul(adjCSpeed))
	}
	if ActionState[VP_UP] {
		camera.Position = camera.Position.Add(camera.Up.Mul(adjCSpeed))
	}
	if ActionState[VP_DOWN] {
		camera.Position = camera.Position.Sub(camera.Up.Mul(adjCSpeed))
	}
	if ActionState[ED_QUIT] {
		fmt.Println("Exiting!")
		glfw.Terminate()
	}
	// Update direction based on mouse input
	camera.UpdateDirection(u.CursorChange().X(), u.CursorChange().Y())
	u.CheckpointCursorChange()
	ViewportTransform = camera.GetTransform()
	InputManager(win, u)
	return nil
}
