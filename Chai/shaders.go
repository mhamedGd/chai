package chai

import (
	"bufio"
	"net/http"
	"strings"

	"github.com/gowebapi/webapi/graphics/webgl"
)

type ShaderSource struct {
	vertexShader   string
	fragmentShader string
}

type ShaderProgram struct {
	ShaderSource     ShaderSource
	AttributesNumber int
	ShaderProgramID  *webgl.Program
}

func UseShader(_sp *ShaderProgram) {
	glRef.UseProgram(_sp.ShaderProgramID)
}

func UnuseShader() {
	//glRef.UseProgram(nil)
}

func (_sp *ShaderProgram) ParseShader(_vertexSource string, _fragmentSource string) {
	_sp.ShaderSource = ShaderSource{_vertexSource, _fragmentSource}
}

func (_sp *ShaderProgram) ParseShaderFromFile(_filePath string) {
	/*
		file, err := os.Open(_filePath)
		if err != nil {
			LogF(err.Error())
		}
	*/
	resp, err := http.Get(app_url + "/" + _filePath)
	if err != nil {
		LogF(err.Error())
	}
	defer resp.Body.Close()
	/*
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			LogF(err.Error())
		}
	*/
	const VERTEX = 0
	const FRAGMENT = 1
	current_type := -1

	shaders := []string{"", ""}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#shader") {
			if strings.Contains(scanner.Text(), "vertex") {
				current_type = VERTEX
			}
			if strings.Contains(scanner.Text(), "fragment") {
				current_type = FRAGMENT
			}
		} else {
			Assert(current_type != -1, "Parse Shader: shader should start with #shader vertex/fragment")
			shaders[current_type] += scanner.Text()
			shaders[current_type] += string('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		LogF(err.Error())
	}

	_sp.ShaderSource = ShaderSource{shaders[0], shaders[1]}
}

func (_sp *ShaderProgram) CreateShaderProgram() {
	_sp.AttributesNumber = 0
	_sp.ShaderProgramID = glRef.CreateProgram()
	vertex_shader := CompileShader(webgl.VERTEX_SHADER, _sp.ShaderSource.vertexShader)
	fragment_shader := CompileShader(webgl.FRAGMENT_SHADER, _sp.ShaderSource.fragmentShader)

	glRef.AttachShader(_sp.ShaderProgramID, vertex_shader)
	glRef.AttachShader(_sp.ShaderProgramID, fragment_shader)

	glRef.LinkProgram(_sp.ShaderProgramID)
	if !glRef.GetProgramParameter(_sp.ShaderProgramID, webgl.LINK_STATUS).Bool() {
		//return webgl.Program(js.Null()), errors.New("link failed: " + glRef.GetProgramInfoLog(program))
		LogF("[LINK FAILED]: " + *glRef.GetProgramInfoLog(_sp.ShaderProgramID))
	}
}

func (_sp *ShaderProgram) AddAttribute(_attributeName string) {
	//BindAttribLocation(_sp.ShaderProgramID, _sp.AttributesNumber, _attributeName)
	glRef.BindAttribLocation(_sp.ShaderProgramID, uint(_sp.AttributesNumber), _attributeName)
	_sp.AttributesNumber += 1
}

func CompileShader(_shaderType uint, _shaderSource string) *webgl.Shader {
	shader := glRef.CreateShader(_shaderType)

	glRef.ShaderSource(shader, _shaderSource)
	glRef.CompileShader(shader)

	if !glRef.GetShaderParameter(shader, webgl.COMPILE_STATUS).Bool() {
		if _shaderType == webgl.FRAGMENT_SHADER {
			LogF("[FRAGMENT SHADER] compile failure: " + *glRef.GetShaderInfoLog(shader))

		} else if _shaderType == webgl.VERTEX_SHADER {
			LogF("[VERTEX SHADER] compile failure: " + *glRef.GetShaderInfoLog(shader))
		}
	}

	return shader
}

func (_sp *ShaderProgram) GetUniformLocation(_uniformName string) *webgl.UniformLocation {
	return glRef.GetUniformLocation(_sp.ShaderProgramID, _uniformName)
}
