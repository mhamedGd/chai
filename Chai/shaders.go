package chai

import (
	"bufio"
	"net/http"
	"strings"
	"syscall/js"
)

type ShaderSource struct {
	vertexShader   string
	fragmentShader string
}

type ShaderProgram struct {
	ShaderSource     ShaderSource
	AttributesNumber int
	ShaderProgramID  js.Value
}

func UseShader(_sp *ShaderProgram) {
	canvasContext.Call("useProgram", _sp.ShaderProgramID)
}

func UnuseShader() {
	//glRef.UseProgram(nil)
	canvasContext.Call("useProgram", js.Null())
}

func (_sp *ShaderProgram) ParseShader(_vertexSource string, _fragmentSource string) {
	_sp.ShaderSource = ShaderSource{_vertexSource, _fragmentSource}
}

func (_sp *ShaderProgram) ParseShaderFromFile(_filePath string) {
	resp, err := http.Get(app_url + "/" + _filePath)
	if err != nil {
		LogF(err.Error())
	}
	defer resp.Body.Close()
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
	_sp.ShaderProgramID = canvasContext.Call("createProgram")
	vertex_shader := CompileShader(canvasContext.Get("VERTEX_SHADER"), _sp.ShaderSource.vertexShader)
	fragment_shader := CompileShader(canvasContext.Get("FRAGMENT_SHADER"), _sp.ShaderSource.fragmentShader)

	canvasContext.Call("attachShader", _sp.ShaderProgramID, vertex_shader)
	canvasContext.Call("attachShader", _sp.ShaderProgramID, fragment_shader)

	canvasContext.Call("linkProgram", _sp.ShaderProgramID)

	if canvasContext.Call("getProgramParameter", _sp.ShaderProgramID, canvasContext.Get("LINK_STATUS")).IsNull() {
		//return webgl.Program(js.Null()), errors.New("link failed: " + glRef.GetProgramInfoLog(program))
		LogF("[LINK FAILED]: " + canvasContext.Call("getProgramInfoLog", _sp.ShaderProgramID).String())
	}
}

func (_sp *ShaderProgram) AddAttribute(_attributeName string) {
	//BindAttribLocation(_sp.ShaderProgramID, _sp.AttributesNumber, _attributeName)
	canvasContext.Call("bindAttribLocation", _sp.ShaderProgramID, _sp.AttributesNumber, _attributeName)
	_sp.AttributesNumber += 1
}

func CompileShader(_shaderType js.Value, _shaderSource string) js.Value {
	shader := canvasContext.Call("createShader", _shaderType)

	canvasContext.Call("shaderSource", shader, _shaderSource)
	canvasContext.Call("compileShader", shader)

	if canvasContext.Call("getShaderParameter", shader, canvasContext.Get("COMPILE_STATUS")).IsNull() {
		if _shaderType.Equal(canvasContext.Get("FRAGMENT_SHADER")) {
			LogF("[FRAGMENT SHADER] compile failure: " + canvasContext.Call("getShaderInfoLog", shader).String())

		} else if _shaderType.Equal(canvasContext.Get("VERTEX_SHADER")) {
			LogF("[VERTEX SHADER] compile failure: " + canvasContext.Call("getShaderInfoLog", shader).String())
		}
	}

	return shader
}

func (_sp *ShaderProgram) GetUniformLocation(_uniformName string) js.Value {
	return canvasContext.Call("getUniformLocation", _sp.ShaderProgramID, _uniformName)
}
