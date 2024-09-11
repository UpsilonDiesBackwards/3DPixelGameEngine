package rendering

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"

	"os"
)

type Shader struct {
	Program uint32
}

func NewShader(vPath, fPath string) (*Shader, error) {
	vSource, err := loadSource(vPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load vertex source code: %v", err)
	}

	fSource, err := loadSource(fPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load fragment source code: %v", err)
	}

	program, err := createProgram(vSource, fSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create program: %v", err)
	}

	return &Shader{Program: program}, nil
}

func loadSource(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return string(data), nil
}

func createProgram(vSource, fSource string) (uint32, error) {
	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	vVertSource, free := gl.Strs(vSource + "\x00")
	gl.ShaderSource(vShader, 1, vVertSource, nil)
	gl.CompileShader(vShader)
	free()
	if err := verifyCompilation(vShader); err != nil {
		return 0, err
	}

	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fVertSource, free := gl.Strs(fSource + "\x00")
	gl.ShaderSource(fShader, 1, fVertSource, nil)
	gl.CompileShader(fShader)
	free()
	if err := verifyCompilation(fShader); err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vShader)
	gl.AttachShader(program, fShader)
	gl.LinkProgram(program)
	if err := verifyProgramLink(program); err != nil {
		return 0, err
	}

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	return program, nil
}

func verifyCompilation(shader uint32) error {
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		log := string(make([]byte, length))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))
		return fmt.Errorf("shader compilation failed: %s", log)
	}
	return nil
}

func verifyProgramLink(program uint32) error {
	var success int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		log := string(make([]byte, length))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))
		return fmt.Errorf("program linking failed: %s", log)
	}
	return nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.Program)
}

func (s *Shader) Delete() {
	gl.DeleteShader(s.Program)
}
