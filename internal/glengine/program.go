package engine

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thijzert/chesseract/internal/assets"
)

type Uniform int

const (
	U_PROJECTION Uniform = iota
	U_CAMERA
	U_TRANSFORM
	U_MATERIAL_SHINE_DAMPER
	U_MATERIAL_REFLECTIVITY
	U_TEXTURE_DIFFUSE
	U_TEXTURE_NORMAL
	U_TEXTURE_SPECULAR
	U_TMP_LIGHT_POS
	U_TMP_LIGHT_COLOUR
	// ...
	U_MAX
)

func (u Uniform) String() string {
	if u == U_PROJECTION {
		return "projectionMatrix"
	} else if u == U_CAMERA {
		return "cameraMatrix"
	} else if u == U_TRANSFORM {
		return "transformationMatrix"
	} else if u == U_MATERIAL_SHINE_DAMPER {
		return "materialShineDamper"
	} else if u == U_MATERIAL_REFLECTIVITY {
		return "materialReflectivity"
	} else if u == U_TEXTURE_DIFFUSE {
		return "diffuseTexture"
	} else if u == U_TEXTURE_NORMAL {
		return "normalMap"
	} else if u == U_TEXTURE_SPECULAR {
		return "specularMap"
	} else if u == U_TMP_LIGHT_POS {
		return "lightPosition"
	} else if u == U_TMP_LIGHT_COLOUR {
		return "lightColour"
	} else {
		return fmt.Sprintf("unknown uniform variable 0x%02x", int(u))
	}
}

type glShader uint32

type glProgram struct {
	programID        uint32
	includedShaders  []glShader
	uniformLocations [U_MAX]int32
}

func (p glProgram) Start() {
	gl.UseProgram(p.programID)
}

func (p glProgram) Stop() {
	gl.UseProgram(0)
}

func (p glProgram) Destroy() {
	p.Stop()
	if p.includedShaders != nil {
		for _, sh := range p.includedShaders {
			gl.DetachShader(p.programID, uint32(sh))
			gl.DeleteShader(uint32(sh))
		}
	}
	gl.DeleteProgram(p.programID)
}

func (p glProgram) UniformInt(attribute Uniform, value int32) {
	location := p.uniformLocations[attribute]
	if location == -1 {
		return
	}

	gl.Uniform1i(location, value)
}

func (p glProgram) UniformFloat(attribute Uniform, value float32) {
	location := p.uniformLocations[attribute]
	if location == -1 {
		return
	}

	gl.Uniform1f(location, value)
}

func (p glProgram) UniformVec3(attribute Uniform, value mgl32.Vec3) {
	location := p.uniformLocations[attribute]
	if location == -1 {
		return
	}

	gl.Uniform3f(location, value[0], value[1], value[2])
}

func (p glProgram) UniformMatrix4(attribute Uniform, value mgl32.Mat4) {
	location := p.uniformLocations[attribute]
	if location == -1 {
		return
	}

	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

func (eng *Engine) loadSimpleProgram(vertexShaderName, fragmentShaderName string) (glProgram, error) {
	vertexShader, err := eng.compileShader(vertexShaderName, gl.VERTEX_SHADER)
	if err != nil {
		return glProgram{}, err
	}
	fragmentShader, err := eng.compileShader(fragmentShaderName, gl.FRAGMENT_SHADER)
	if err != nil {
		return glProgram{}, err
	}

	return eng.linkProgram(vertexShader, fragmentShader)
}

func (eng *Engine) linkProgram(shaders ...glShader) (glProgram, error) {
	var rv glProgram

	prog := gl.CreateProgram()
	for _, sh := range shaders {
		gl.AttachShader(prog, uint32(sh))
	}
	gl.LinkProgram(prog)

	err := eng.getProgramError(prog, gl.LINK_STATUS)
	if err != nil {
		return rv, fmt.Errorf("error linking program: %v", err)
	}

	// gl.ValidateProgram(prog)
	// err = eng.getProgramError(prog, gl.VALIDATE_STATUS)
	// if err != nil {
	// 	return rv, fmt.Errorf("error validating program: %v", err)
	// }

	rv.programID = prog
	rv.includedShaders = make([]glShader, len(shaders))
	copy(rv.includedShaders, shaders)

	for i := Uniform(0); i < U_MAX; i++ {
		rv.uniformLocations[i] = gl.GetUniformLocation(prog, gl.Str(i.String()+"\x00"))
	}

	// Connect texture units
	textureUniforms := []Uniform{
		U_TEXTURE_DIFFUSE,
		U_TEXTURE_NORMAL,
		U_TEXTURE_SPECULAR,
	}
	rv.Start()
	tx := int32(0)
	for _, u := range textureUniforms {
		if rv.uniformLocations[u] == 0 {
			continue
		}
		rv.UniformInt(u, tx)
		tx++
	}
	rv.Stop()

	return rv, nil
}

func (eng *Engine) compileShader(name string, shaderType uint32) (glShader, error) {
	source, err := assets.GetAsset("shaders/" + name + ".glsl")
	if err != nil {
		return 0, err
	}

	rv := gl.CreateShader(shaderType)

	source = append(source, 0)
	csources, free := gl.Strs(string(source))
	gl.ShaderSource(rv, 1, csources, nil)
	free()
	gl.CompileShader(rv)

	err = eng.getShaderError(rv, gl.COMPILE_STATUS)
	if err != nil {
		return 0, err
	}

	return glShader(rv), nil
}

func (eng *Engine) getShaderError(shader uint32, field uint32) error {
	return eng.getCompilationError(gl.GetShaderiv, gl.GetShaderInfoLog, shader, field)
}

func (eng *Engine) getProgramError(prog uint32, field uint32) error {
	return eng.getCompilationError(gl.GetProgramiv, gl.GetProgramInfoLog, prog, field)
}

type glInformationVectorFunc func(uint32, uint32, *int32)
type glInfoLogFunc func(uint32, int32, *int32, *uint8)

func (eng *Engine) getCompilationError(iv glInformationVectorFunc, ilog glInfoLogFunc, obj uint32, field uint32) error {
	var status int32
	iv(obj, field, &status)
	if status != gl.FALSE {
		return nil
	}

	var logLength int32
	iv(obj, gl.INFO_LOG_LENGTH, &logLength)
	logOutput := strings.Repeat("\x00", int(logLength)+1)

	ilog(obj, logLength, nil, gl.Str(logOutput))
	return fmt.Errorf("%s", logOutput)
}
