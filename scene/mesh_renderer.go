/*
Copyright (c) 2018 HaakenLabs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package scene

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/haakenlabs/arc/graphics"
	"github.com/haakenlabs/arc/system/instance"
)

var _ Drawable = &MeshRenderer{}

type MeshRenderer struct {
	BaseComponent

	material   *Material
	cullFace   bool
	depthWrite bool
	wireframe  bool
}

func NewMeshRenderer() *MeshRenderer {
	c := &MeshRenderer{
		cullFace:   true,
		depthWrite: true,
	}

	c.SetName("MeshRenderer")
	instance.MustAssign(c)

	return c
}

// MeshRenderer Functions

func (m *MeshRenderer) SetMaterial(material *Material) {
	m.material = material
}

func (m *MeshRenderer) GetMaterial() *Material {
	return m.material
}

func (m *MeshRenderer) Draw(camera *Camera) {
	if m.material == nil {
		return
	}

	m.material.Bind()

	if m.material.SupportsDeferredPath() {
		if camera.ActiveRenderPath() == RenderPathForward {
			m.material.Shader().SetSubroutine(graphics.ShaderComponentFragment, "forward_pass")
		} else {
			m.material.Shader().SetSubroutine(graphics.ShaderComponentFragment, "deferred_pass_geometry")
		}
	}

	m.DrawShader(m.material.Shader(), camera)

	m.material.Unbind()
}

func (m *MeshRenderer) DrawShader(shader *graphics.Shader, camera *Camera) {
	if shader == nil && m.GameObject() == nil {
		return
	}

	// FIXME: Move this somewhere out of the render loop
	var meshes []*graphics.Mesh
	components := m.GameObject().Components()
	for i := range components {
		if meshFilter, ok := components[i].(*MeshFilter); ok {
			if mesh := meshFilter.Mesh(); mesh != nil {
				meshes = append(meshes, mesh)
			}
		}
	}

	if len(meshes) == 0 {
		return
	}

	shader.SetUniform("v_model_matrix", m.GetTransform().ActiveMatrix())
	shader.SetUniform("v_view_matrix", camera.ViewMatrix())
	shader.SetUniform("v_projection_matrix", camera.ProjectionMatrix())
	shader.SetUniform("v_normal_matrix", camera.NormalMatrix())
	shader.SetUniform("f_camera", camera.CameraPosition())

	if !m.cullFace {
		gl.Disable(gl.CULL_FACE)
	}
	if !m.depthWrite {
		gl.DepthMask(false)
	}
	if m.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	for i := range meshes {
		meshes[i].Bind()

		if meshes[i].Indexed() {
			gl.DrawElements(gl.TRIANGLES, int32(len(meshes[i].Triangles())), gl.UNSIGNED_INT, nil)
		} else {
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(meshes[i].Vertices())))
		}

		meshes[i].Unbind()

	}

	if m.wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	if !m.depthWrite {
		gl.DepthMask(true)
	}
	if !m.cullFace {
		gl.Enable(gl.CULL_FACE)
	}
}

func (m *MeshRenderer) CullFaceEnabled() bool {
	return m.cullFace
}

func (m *MeshRenderer) DepthWriteEnabled() bool {
	return m.depthWrite
}

func (m *MeshRenderer) WireframeEnabled() bool {
	return m.wireframe
}

func (m *MeshRenderer) SetCullFaceEnabled(enable bool) {
	m.cullFace = enable
}

func (m *MeshRenderer) SetDepthWriteEnabled(enable bool) {
	m.depthWrite = enable
}

func (m *MeshRenderer) SetWireframeEnabled(enable bool) {
	m.wireframe = enable
}

func (m *MeshRenderer) SupportsDeferred() bool {
	if m.material != nil {
		return m.material.SupportsDeferredPath()
	}

	return false
}
