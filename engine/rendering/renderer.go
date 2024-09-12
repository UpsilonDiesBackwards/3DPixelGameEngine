package rendering

import (
	"3DPixelGameEngine/engine"
	"fmt"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type Renderer struct {
	window   *Window
	program  uint32
	ubo      uint32
	Objects  map[string]*RenderableObject
	camera   *Camera
	lastTime time.Time
}

func NewRenderer(window *Window) *Renderer {
	if err := gl.Init(); err != nil {
		fmt.Println("Error initializing OpenGL:", err)
		engine.GetGLErrorVerbose()
	} else {
		fmt.Println("Init OpenGl using version: ", gl.GoStr(gl.GetString(gl.VERSION)))
	}

	shader, err := NewShader("engine/res/shaders/shader.vert", "engine/res/shaders/shader.frag")
	if err != nil {
		fmt.Println("Failed to create shader: ", err)
	}

	var ubo uint32
	gl.GenBuffers(1, &ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferData(gl.UNIFORM_BUFFER, 3*16*4, nil, gl.DYNAMIC_DRAW)

	blockIndex := gl.GetUniformBlockIndex(shader.Program, gl.Str("PerspectiveBlock\x00"))
	gl.UniformBlockBinding(shader.Program, blockIndex, 1)

	return &Renderer{
		window:  window,
		program: shader.Program,
		ubo:     ubo,
		Objects: make(map[string]*RenderableObject),
		camera: &Camera{
			Position:    mgl64.Vec3{0, 0, 3},
			Front:       mgl64.Vec3{0, 0, -1},
			Up:          mgl64.Vec3{0, 1, 0},
			WorldUp:     mgl64.Vec3{0, 1, 0},
			Yaw:         -90,
			Pitch:       0,
			Speed:       12,
			Sensitivity: 0.075,
			Fov:         60,
		},
		lastTime: time.Now(),
	}
}

func (r *Renderer) Draw() {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	view := r.camera.GetTransform()
	projection := mgl32.Perspective(r.camera.GetFov(),
		float32(r.window.GetWidth())/float32(r.window.GetHeight()), 0.1, 100.0)

	for _, obj := range r.Objects {
		obj.Draw(r.program, view, projection)
	}

	r.window.SwapBuffers()
}

func (r *Renderer) CalculateDeltaTime() float64 {
	currentTime := time.Now()
	deltaTime := currentTime.Sub(r.lastTime).Seconds()
	r.lastTime = currentTime
	return deltaTime
}
