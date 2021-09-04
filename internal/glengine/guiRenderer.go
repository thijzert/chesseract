package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/thijzert/chesseract/internal/gui"
)

type GUIManager struct {
	Program glProgram
	Quad    rawModel
	layers  []guiLayer
}

type LayerName uint32

type guiLayer struct {
	name    LayerName
	source  gui.GUIContext
	texture texture
	enabled bool
}

func (m *GUIManager) AddLayer(name LayerName, source gui.GUIContext) {
	layer := guiLayer{
		name:    name,
		source:  source,
		enabled: true,
	}
	m.layers = append(m.layers, layer)
}

func (m *GUIManager) assignTexture(layer *guiLayer) {
	if layer.texture != 0 {
		return
	}

	var rv uint32
	gl.GenTextures(1, &rv)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, rv)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	layer.texture = texture(rv)
}

func (m *GUIManager) RemoveNamedLayer(name LayerName) {
	for j := range m.layers {
		i := len(m.layers) - j - 1
		if m.layers[i].name == name {
			m.layers = append(m.layers[:i], m.layers[i+1:]...)
		}
	}
}

func (m *GUIManager) ToggleNamedLayer(name LayerName, enabled bool) {
	for i := range m.layers {
		if m.layers[i].name == name {
			m.layers[i].enabled = enabled
		}
	}
}

func (m *GUIManager) DrawGUI() {

	m.Program.Start()

	gl.BindVertexArray(m.Quad.VAO)
	gl.EnableVertexAttribArray(0)
	gl.Disable(gl.DEPTH_TEST)

	gl.ActiveTexture(gl.TEXTURE0)

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.MAX)
	gl.BlendFunc(gl.BLEND_SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for _, l := range m.layers {
		if !l.enabled {
			continue
		}

		m.assignTexture(&l)
		gl.BindTexture(gl.TEXTURE_2D, uint32(l.texture))

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(l.source.Pixels.Rect.Size().X),
			int32(l.source.Pixels.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(l.source.Pixels.Pix))

		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, m.Quad.NVertices)

		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

	gl.Disable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)

	m.Program.Stop()
}
