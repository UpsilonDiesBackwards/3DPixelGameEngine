package rendering

import (
	"3DPixelGameEngine/engine/obj"
	"fmt"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
)

type RenderableObject struct {
	DecodedObject obj.DecodedObject
	ModelMatrix   mgl32.Mat4

	Position mgl32.Vec3
	Scale    mgl32.Vec3
	Rotation mgl32.Quat
}

func NewObject(decodedObject *obj.DecodedObject) *RenderableObject {
	object := &RenderableObject{
		DecodedObject: *decodedObject,
		ModelMatrix:   mgl32.Ident4(),

		Position: mgl32.Vec3{0, 0, 0},
		Scale:    mgl32.Vec3{1, 1, 1},
		Rotation: mgl32.QuatIdent(),
	}
	object.setup()
	return object
}

func (o *RenderableObject) setup() {
	var vao, vbo, ebo uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.BufferData(gl.ARRAY_BUFFER, len(o.DecodedObject.Vertices)*4, gl.Ptr(o.DecodedObject.Vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(5*4))
	gl.EnableVertexAttribArray(2)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(o.DecodedObject.Indices)*4, gl.Ptr(o.DecodedObject.Indices), gl.STATIC_DRAW)

	o.DecodedObject.VAO = vao
	o.DecodedObject.VBO = vbo
	o.DecodedObject.EBO = ebo

	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("OpenGL error during VAO/VBO setup: %v", err)
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (o *RenderableObject) Draw(program uint32, project, camera mgl32.Mat4) {
	gl.BindVertexArray(o.DecodedObject.VAO)
	if err := gl.GetError(); err != gl.NO_ERROR {
		log.Printf("Error after glBindVertexArray: %v", err)
	}

	o.UpdateModelMatrix()

	projData := project[:]
	camData := camera[:]
	modelData := o.ModelMatrix[:]

	projBuf := make([]float32, 16)
	camBuf := make([]float32, 16)
	modelBuf := make([]float32, 16)

	copy(projBuf, projData)
	copy(camBuf, camData)
	copy(modelBuf, modelData)

	gl.BindBuffer(gl.UNIFORM_BUFFER, o.DecodedObject.UBO)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 16*4, gl.Ptr(projBuf))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 16*4, 16*4, gl.Ptr(camBuf))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 32*4, 16*4, gl.Ptr(modelBuf))

	if err := gl.GetError(); err != gl.NO_ERROR {
		fmt.Println(program) // WHY IS PROGRAM 3?
		fmt.Println(o.DecodedObject.VAO)
		log.Printf("OpenGL error after uniform loc: %v", err)
	}

	gl.UseProgram(program)

	gl.DrawElements(gl.TRIANGLES, int32(len(o.DecodedObject.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func (o *RenderableObject) SetPosition(x, y, z float32) {
	o.Position = mgl32.Vec3{x, y, z}
}

func (o *RenderableObject) SetScale(x, y, z float32) {
	o.Scale = mgl32.Vec3{x, y, z}
}

func (o *RenderableObject) SetRotation(angle float32, axis mgl32.Vec3) {
	o.Rotation = mgl32.QuatRotate(angle, axis)
}

func (o *RenderableObject) UpdateModelMatrix() {
	o.ModelMatrix = mgl32.Translate3D(o.Position.X(), o.Position.Y(), o.Position.Z()).
		Mul4(o.Rotation.Mat4()).
		Mul4(mgl32.Scale3D(o.Scale.X(), o.Scale.Y(), o.Scale.Z()))
}
