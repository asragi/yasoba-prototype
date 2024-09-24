package frontend

import "github.com/hajimehoshi/ebiten/v2"

type ShaderId int

const (
	ShaderNone ShaderId = iota
	ShaderDisappear
)

type Shader struct {
	frame    int
	shader   *ebiten.Shader
	uniforms map[string]interface{}
}

func NewShader(shader *ebiten.Shader) *Shader {
	return &Shader{
		frame:    0,
		shader:   shader,
		uniforms: make(map[string]interface{}),
	}
}

func (s *Shader) Reset() {
	s.frame = 0
}

func (s *Shader) Update() {
	s.frame++
	s.SetUniforms("Time", float32(s.frame)/60)
}

func (s *Shader) SetUniforms(name string, value interface{}) {
	s.uniforms[name] = value
}

func (s *Shader) GetShader() *ebiten.Shader {
	return s.shader
}

func (s *Shader) GetUniforms() map[string]interface{} {
	return s.uniforms
}
