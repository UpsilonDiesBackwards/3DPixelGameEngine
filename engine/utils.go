package engine

import (
	"github.com/go-gl/gl/v4.2-compatibility/gl"
	"log"
	"unsafe"
)

func GetGLErrorVerbose() {
	if err := gl.GetError(); err != 0 {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(func(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
			log.Printf("OpenGL Debug Message (source: 0x%X, type: 0x%X, id: %d, severity: 0x%X): %s\n", source, gltype, id, severity, message)
		}, gl.Ptr(nil))
	}
}
