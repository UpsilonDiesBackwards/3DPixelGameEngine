package io

import (
	"3DPixelGameEngine/engine/rendering"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type KeyAction int

// UserInput Types of user input
type UserInput struct {
	InitialAction  bool
	cursor         mgl64.Vec2
	cursorChange   mgl64.Vec2
	cursorLast     mgl64.Vec2
	bufferedChange mgl64.Vec2
}

const (
	NO_ACTION = iota
	VP_FORW
	VP_BACK
	VP_LEFT
	VP_RGHT
	VP_UP
	VP_DOWN
	ED_QUIT
)

var ActionState = make(map[KeyAction]bool)

var keyToActionMap = map[glfw.Key]KeyAction{
	glfw.KeyW:      VP_FORW,
	glfw.KeyS:      VP_BACK,
	glfw.KeyA:      VP_LEFT,
	glfw.KeyD:      VP_RGHT,
	glfw.KeySpace:  VP_UP,
	glfw.KeyC:      VP_DOWN,
	glfw.KeyEscape: ED_QUIT,
}

func InputManager(win *rendering.Window, uI *UserInput) {
	win.SetKeyCallback(KeyCallBack)
	win.SetCursorPosCallback(uI.MouseCallBack)
}

func KeyCallBack(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifier glfw.ModifierKey) {
	a, ok := keyToActionMap[key]
	if !ok {
		return
	}
	switch action {
	case glfw.Press:
		ActionState[a] = true
	case glfw.Release:
		ActionState[a] = false
	}
}

func (cInput *UserInput) Cursor() mgl64.Vec2       { return cInput.cursor }
func (cInput *UserInput) CursorChange() mgl64.Vec2 { return cInput.cursorChange }
func (cInput *UserInput) CheckpointCursorChange() {
	cInput.cursorChange = cInput.bufferedChange
	cInput.bufferedChange = mgl64.Vec2{0, 0}
}
func (cInput *UserInput) MouseCallBack(win *glfw.Window, xpos, ypos float64) {
	if cInput.InitialAction {
		cInput.cursorLast = mgl64.Vec2{xpos, ypos}
		cInput.InitialAction = false
	}
	cInput.bufferedChange = mgl64.Vec2{xpos - cInput.cursorLast.X(), ypos - cInput.cursorLast.Y()}
	cInput.cursorLast = mgl64.Vec2{xpos, ypos}
}
