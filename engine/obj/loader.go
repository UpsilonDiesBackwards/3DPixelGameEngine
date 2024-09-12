package obj

import (
	"bufio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Object struct {
	Name      string
	Faces     []Face
	materials []string
}

type DecodedObject struct {
	Objects    []Object
	Materials  map[string]*Material
	Vertices   []float32
	Indices    []uint32
	Normals    []float32
	UVs        []float32
	line       uint
	objCur     *Object
	matCur     *Material
	mtlDir     string
	matLib     string
	smoothCurr bool

	VAO uint32
	VBO uint32
	EBO uint32
	UBO uint32
	
	Texture uint32
}

type Face struct {
	Vertices []int
	Normals  []int
	UVs      []int
	Material string
}

type Material struct {
	Name     string
	Opacity  float32
	Metallic float32
	Ambient  mgl32.Vec3
	Diffuse  mgl32.Vec3
	Specular mgl32.Vec3
	Emission mgl32.Vec3

	Texture string
}

func LoadModel(objPath, mtlPath string) (*DecodedObject, error) {
	objFile, err := os.Open(objPath)
	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	mtlFile, err := os.Open(mtlPath)
	if err != nil {
		return nil, err
	}
	defer mtlFile.Close()

	decObj, err := DecodeObject(objFile, mtlFile)
	if err != nil {
		return nil, err
	}

	decObj.mtlDir = filepath.Dir(objPath)
	return decObj, nil
}

func DecodeObject(objReader, mtlReader io.Reader) (*DecodedObject, error) {
	dec := new(DecodedObject)
	dec.Objects = make([]Object, 0)
	dec.Materials = make(map[string]*Material)
	dec.Vertices = make([]float32, 0)
	dec.Normals = make([]float32, 0)
	dec.UVs = make([]float32, 0)
	dec.line = 1

	err := dec.parse(objReader, dec.parseObjLine)
	if err != nil {
		return nil, err
	}

	dec.matCur = nil
	dec.line = 1
	err = dec.parse(mtlReader, dec.parseMtlLine)
	if err != nil {
	}

	return dec, nil
}

func (dec *DecodedObject) parse(reader io.Reader, parseLine func(string) error) error {
	buf := bufio.NewReader(reader)
	dec.line = 1
	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		line = strings.Trim(line, "\r\n\t ")
		perr := parseLine(line)
		if perr != nil {
			return perr
		}

		if err == io.EOF {
			break
		}
		dec.line++
	}
	return nil
}

func (dec *DecodedObject) parseObjLine(line string) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}

	lType := fields[0]
	if strings.HasPrefix(lType, "#") {
		return nil
	}

	switch lType {
	case "mtllib":
		return dec.parseMatlib(fields[1:])
	case "g":
		return dec.parseObject(fields[1:])
	case "o":
		return dec.parseObject(fields[1:])
	case "v":
		return dec.parseVertex(fields[1:])
	case "vn":
		return dec.parseNormal(fields[1:])
	case "vt":
		return dec.parseTex(fields[1:])
	case "f":
		return dec.parseFace(fields[1:])
	case "usemtl":
		return dec.parseUsemtl(fields[1:])
	case "s":
		return dec.parseSmooth(fields[1:])
	default:
		fmt.Println("field not supported: " + lType + " OBJ")
	}
	return nil
}

func (dec *DecodedObject) parseVertex(i []string) error {
	if len(i) < 3 {
		fmt.Println("Less than 3 vertices in line")
	}

	for _, f := range i[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		dec.Vertices = append(dec.Vertices, float32(val))
	}
	return nil
}

func (dec *DecodedObject) parseNormal(i []string) error {
	if len(i) < 3 {
		fmt.Println("Less than 3 normals in line")
	}

	for _, f := range i[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		dec.Normals = append(dec.Normals, float32(val))
	}
	return nil
}

func (dec *DecodedObject) parseTex(i []string) error {
	if len(i) < 2 {
		fmt.Println("Less than 2 UVs in line")
	}

	for _, f := range i[:2] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		dec.UVs = append(dec.UVs, float32(val))
	}
	return nil
}

func (dec *DecodedObject) parseObject(i []string) error {
	if len(i) < 1 {
		fmt.Println("Object line with no fields")
	}

	dec.Objects = append(dec.Objects, makeObject(i[0]))
	dec.objCur = &dec.Objects[len(dec.Objects)-1]
	return nil
}

func makeObject(name string) Object {
	var o Object
	o.Name = name
	o.Faces = make([]Face, 0)
	o.materials = make([]string, 0)
	return o
}

func (dec *DecodedObject) parseFace(i []string) error {
	if dec.objCur == nil {
		err := dec.parseObject([]string{fmt.Sprintf("unnamed%d", dec.line)})
		if err != nil {
			return err
		}
	}

	if len(dec.objCur.materials) == 0 && dec.matCur != nil {
		dec.objCur.materials = append(dec.objCur.materials, dec.matCur.Name)
	}

	if len(i) < 3 {
		return fmt.Errorf("less than 3 vertices in face definition at line %d", dec.line)
	}

	var face Face
	face.Vertices = make([]int, len(i))
	face.Normals = make([]int, len(i))
	face.UVs = make([]int, len(i))

	if dec.matCur != nil {
		face.Material = dec.matCur.Name
	}

	for pos, f := range i {
		vFields := strings.Split(f, "/")
		if len(vFields) < 1 {
			return fmt.Errorf("face field with no parts at line %d", dec.line)
		}

		val, err := strconv.Atoi(vFields[0])
		if err != nil {
			return fmt.Errorf("invalid vertex index: %s at line %d", vFields[0], dec.line)
		}
		if val > 0 {
			face.Vertices[pos] = val - 1
		} else if val < 0 {
			cur := (len(dec.Vertices) / 3) - 1
			face.Vertices[pos] = cur + val + 1
		} else {
			return fmt.Errorf("vertex index cannot be 0 at line %d", dec.line)
		}

		if len(vFields) > 1 && vFields[1] != "" {
			val, err := strconv.Atoi(vFields[1])
			if err != nil {
				return fmt.Errorf("invalid UV index: %s at line %d", vFields[1], dec.line)
			}
			if val > 0 {
				face.UVs[pos] = val - 1
			} else if val < 0 {
				curr := (len(dec.UVs) / 2) - 1
				face.UVs[pos] = curr + val + 1
			} else {
				return fmt.Errorf("UV index cannot be 0 at line %d", dec.line)
			}
		} else {
			face.UVs[pos] = math.MaxUint32
		}

		if len(vFields) > 2 && vFields[2] != "" {
			val, err := strconv.Atoi(vFields[2])
			if err != nil {
				return fmt.Errorf("invalid normal index: %s at line %d", vFields[2], dec.line)
			}
			if len(dec.Normals) > 0 {
				if val > 0 {
					face.Normals[pos] = val - 1
				} else if val < 0 {
					curr := (len(dec.Normals) / 3) - 1
					face.Normals[pos] = curr + val + 1
				} else {
					return fmt.Errorf("normal index cannot be 0 at line %d", dec.line)
				}
			}
		} else {
			face.Normals[pos] = math.MaxUint32
		}
	}

	if len(face.Vertices) >= 3 {
		for j := 1; j < len(face.Vertices)-1; j++ {
			dec.Indices = append(dec.Indices, uint32(face.Vertices[0]), uint32(face.Vertices[j]), uint32(face.Vertices[j+1]))
		}
	}

	dec.objCur.Faces = append(dec.objCur.Faces, face)
	return nil
}

func (dec *DecodedObject) parseMatlib(i []string) error {
	if len(i) < 1 {
		fmt.Println("Material library line with no fields")
	}
	dec.matLib = i[0]
	return nil
}

func (dec *DecodedObject) parseSmooth(i []string) error {
	if len(i) < 1 {
		fmt.Println("Smooth line with no fields")
	}
	if i[0] == "0" || i[0] == "off" {
		dec.smoothCurr = false
		return nil
	}
	dec.smoothCurr = true
	return nil
}

func (dec *DecodedObject) parseUsemtl(i []string) error {
	if len(i) < 1 {
		fmt.Println("Usemtl line with no fields")
	}

	if dec.objCur == nil {
		dec.parseObject([]string{fmt.Sprintf("unnamed%d", dec.line)})
	}

	name := i[0]
	mat := dec.Materials[name]
	if mat == nil {
		mat = new(Material)
		mat.Name = name
		dec.Materials[name] = mat
	}
	dec.objCur.materials = append(dec.objCur.materials, name)
	dec.matCur = mat
	return nil
}

// MTL PARSE

func (dec *DecodedObject) parseMtlLine(l string) error {
	fields := strings.Fields(l)
	if len(fields) == 0 {
		return nil
	}

	lType := fields[0]
	if strings.HasPrefix(lType, "#") {
		return nil
	}

	switch lType {
	case "newmtl":
		return dec.parseNewmtl(fields[1:])
	case "d":
		return dec.parseDissolve(fields[1:])
	case "Ka":
		return dec.parseKa(fields[1:])
	case "Ke":
		return dec.parseKe(fields[1:])
	case "Ks":
		return dec.parseKs(fields[1:])
	default:
		fmt.Println("Field not supported: " + lType + " MTL")
	}
	return nil
}

func (dec *DecodedObject) parseNewmtl(i []string) error {
	if len(i) < 1 {
		fmt.Println("Newmtl line with no fields")
	}

	name := i[0]
	mat := dec.Materials[name]
	if mat == nil {
		mat = new(Material)
		mat.Name = name
		dec.Materials[name] = mat
	}
	dec.matCur = mat
	return nil
}

func (dec *DecodedObject) parseDissolve(i []string) error {
	if len(i) < 1 {
		fmt.Println("Dissolve line with no fields")
	}
	val, err := strconv.ParseFloat(i[0], 32)
	if err != nil {
		fmt.Println("Dissolve float error")
	}
	dec.matCur.Opacity = float32(val)
	return nil
}

func (dec *DecodedObject) parseKa(i []string) error {
	if len(i) < 3 {
		fmt.Println("Ka line with less than 3 fields")
	}

	var color mgl32.Vec3
	for pos, f := range i[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		color[pos] = float32(val)
	}
	dec.matCur.Ambient = color
	return nil
}

func (dec *DecodedObject) parseKe(fields []string) error {

	if len(fields) < 3 {
		fmt.Println("'Ke' with less than 3 fields")
	}
	var colors [3]float32
	for pos, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		colors[pos] = float32(val)
	}
	dec.matCur.Emission = colors
	return nil
}

func (dec *DecodedObject) parseKs(fields []string) error {

	if len(fields) < 3 {
		fmt.Println("'Ks' with less than 3 fields")
	}
	var colors [3]float32
	for pos, f := range fields[:3] {
		val, err := strconv.ParseFloat(f, 32)
		if err != nil {
			return err
		}
		colors[pos] = float32(val)
	}
	dec.matCur.Specular = colors
	return nil
}
