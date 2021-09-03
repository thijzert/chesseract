package engine

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thijzert/chesseract/internal/assets"
)

type rawModel struct {
	VAO       uint32
	Material  material
	NVertices int32
}

type modelManager struct {
	vaos []uint32
	vbos []uint32
}

func (m *modelManager) LoadModel(vertices []float32, uvmap []float32, normals []float32, tangents []float32, indices []int32) rawModel {
	vao := m.newVAO()

	m.bindIndicesBuffer(indices)
	m.storeDataInAttributeList(0, 3, vertices)
	m.storeDataInAttributeList(1, 2, uvmap)
	m.storeDataInAttributeList(2, 3, normals)
	m.storeDataInAttributeList(3, 3, tangents)

	gl.BindVertexArray(0)

	return rawModel{
		VAO:       vao,
		NVertices: int32(len(indices)),
	}
}

func (m *modelManager) DefaultQuad() rawModel {
	vao := m.newVAO()

	m.storeDataInAttributeList(0, 2, []float32{-1, 1, -1, -1, 1, 1, 1, -1})

	gl.BindVertexArray(0)

	return rawModel{
		VAO:       vao,
		NVertices: 4,
	}
}

func (m *modelManager) Cleanup() {
	if m.vaos != nil {
		gl.DeleteVertexArrays(int32(len(m.vaos)), &m.vaos[0])
	}
	if m.vaos != nil {
		gl.DeleteBuffers(int32(len(m.vbos)), &m.vbos[0])
	}
}

func (m *modelManager) newVBO() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	m.vbos = append(m.vbos, vbo)

	return vbo
}

func (m *modelManager) newVAO() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	m.vaos = append(m.vaos, vao)

	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)

	return vao
}

func (m *modelManager) bindIndicesBuffer(indices []int32) {
	vbo := m.newVBO()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
}

func (m *modelManager) storeDataInAttributeList(attributeNumber uint32, coordinateSize int32, data []float32) {
	vbo := m.newVBO()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	gl.VertexAttribPointer(attributeNumber, coordinateSize, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (m *modelManager) LoadModelAsset(name string) (rawModel, error) {

	obj, err := assets.GetAsset("models/" + name + ".obj")
	if err != nil {
		return rawModel{}, err
	}

	indices := []int32{}
	wfCount := 0
	waveFront := make(map[string]int32)

	vertexLib := [][]float32{}
	uvLib := [][]float32{}
	normalLib := [][]float32{}

	lines := strings.Split(string(obj), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 3 {
		} else if line[0] == '#' {
		} else if line[0:2] == "v " {
			// parse vertex
			vtx := []float32{0, 0, 0}
			_, err := fmt.Sscanf(line, "v %f %f %f", &vtx[0], &vtx[1], &vtx[2])
			if err != nil {
				return rawModel{}, fmt.Errorf("error at line %d: %v", i+1, err)
			}
			vertexLib = append(vertexLib, vtx)
		} else if line[0:3] == "vt " {
			// parse UV-coordinate
			uv := []float32{0, 0}
			_, err := fmt.Sscanf(line, "vt %f %f", &uv[0], &uv[1])
			if err != nil {
				return rawModel{}, fmt.Errorf("error at line %d: %v", i+1, err)
			}
			uvLib = append(uvLib, uv)
		} else if line[0:3] == "vn " {
			// parse normal
			norm := []float32{0, 0, 0}
			_, err := fmt.Sscanf(line, "vn %f %f %f", &norm[0], &norm[1], &norm[2])
			if err != nil {
				return rawModel{}, fmt.Errorf("error at line %d: %v", i+1, err)
			}
			normalLib = append(normalLib, norm)
		} else if line[0:2] == "f " {
			vtxs := strings.Split(line[2:], " ")
			if len(vtxs) < 3 {
				return rawModel{}, fmt.Errorf("invalid polygon at line %d: %v", i+1, line)
			}
			for _, vtx := range vtxs {
				if _, ok := waveFront[vtx]; !ok {
					waveFront[vtx] = int32(wfCount)
					wfCount++
				}
			}
			for j, vtx := range vtxs[2:] {
				indices = append(indices, waveFront[vtxs[0]], waveFront[vtxs[j+1]], waveFront[vtx])
			}
		} else if line[0:2] == "s " {
			// todo: ?
		} else if line[0:2] == "g " {
			// todo: groups
		} else if strings.HasPrefix(line, "mtllib ") || strings.HasPrefix(line, "usemtl ") {
			// todo: material library
		} else {
			return rawModel{}, fmt.Errorf("error at line %d: can't parse \"%s\"", i+1, line)
		}
	}

	if len(indices)%3 != 0 {
		return rawModel{}, fmt.Errorf("model error: have %d indices; this mesh is improperly triangulised", len(indices))
	}

	vertices := make([]float32, 3*wfCount)
	uvmap := make([]float32, 2*wfCount)
	normals := make([]float32, 3*wfCount)
	tangents := make([]float32, 3*wfCount)

	for vtx, idx := range waveFront {
		var i, j, k int
		if _, err := fmt.Sscanf(vtx, "%d/%d/%d", &i, &j, &k); err != nil {
			return rawModel{}, fmt.Errorf("parse error in face '%s': %v", vtx, err)
		}

		if i < 1 || i > len(vertexLib) {
			return rawModel{}, fmt.Errorf("parse error in face '%s': vertex index out of bounds", vtx)
		}
		if j < 1 || j > len(uvLib) {
			return rawModel{}, fmt.Errorf("parse error in face '%s': UV index out of bounds", vtx)
		}
		if k < 1 || k > len(normalLib) {
			return rawModel{}, fmt.Errorf("parse error in face '%s': normal index out of bounds", vtx)
		}

		copy(vertices[3*idx:], vertexLib[i-1])
		copy(uvmap[2*idx:], uvLib[j-1])
		copy(normals[3*idx:], normalLib[k-1])
	}

	for idx := range indices {
		if idx%3 != 0 {
			continue
		}

		i, j, k := indices[idx], indices[idx+1], indices[idx+2]

		dp1x, dp1y, dp1z := vertices[3*j+0]-vertices[3*i+0], vertices[3*j+1]-vertices[3*i+1], vertices[3*j+2]-vertices[3*i+2]
		dp2x, dp2y, dp2z := vertices[3*k+0]-vertices[3*i+0], vertices[3*k+1]-vertices[3*i+1], vertices[3*k+2]-vertices[3*i+2]

		duv1x, duv1y := uvmap[2*j+0]-uvmap[2*i+0], uvmap[2*j+1]-uvmap[2*i+1]
		duv2x, duv2y := uvmap[2*k+0]-uvmap[2*i+0], uvmap[2*k+1]-uvmap[2*i+1]

		r := 1.0 / (duv1x*duv2y - duv1y*duv2x)

		tanx := (dp1x*duv2y - dp2x*duv1y) * r
		tany := (dp1y*duv2y - dp2y*duv1y) * r
		tanz := (dp1z*duv2y - dp2z*duv1y) * r

		tangents[3*i], tangents[3*i+1], tangents[3*i+2] = tanx, tany, tanz
		tangents[3*j], tangents[3*j+1], tangents[3*j+2] = tanx, tany, tanz
		tangents[3*k], tangents[3*k+1], tangents[3*k+2] = tanx, tany, tanz
	}

	if len(indices) == 0 {
		return rawModel{}, fmt.Errorf("this model is empty")
	}

	return m.LoadModel(vertices, uvmap, normals, tangents, indices), nil
}

type entity struct {
	model    rawModel
	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    mgl32.Vec3
}

func (e entity) getTransformation() mgl32.Mat4 {
	rv := mgl32.Ident4()

	rv = rv.Mul4(mgl32.Translate3D(e.position[0], e.position[1], e.position[2]))

	rv = rv.Mul4(mgl32.HomogRotate3DX(e.rotation[0]))
	rv = rv.Mul4(mgl32.HomogRotate3DY(e.rotation[1]))
	rv = rv.Mul4(mgl32.HomogRotate3DZ(e.rotation[2]))

	rv = rv.Mul4(mgl32.Scale3D(e.scale[0], e.scale[1], e.scale[2]))

	return rv
}
