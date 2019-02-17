package gosigl

import (
	opengl "github.com/go-gl/gl/v4.1-core/gl"
)

type VertexBufferObject uint32
type VertexArrayObject uint32
type VertexObject struct {
	Id VertexBufferObject
	AttribId [32]VertexArrayObject
	attribStride [32]int
	numAttributes int
}

const Line = uint32(opengl.LINE)
const Triangles = uint32(opengl.TRIANGLES)

var vertexDrawMode = Triangles

func SetLineWidth(width float32) {
	opengl.LineWidth(width)
}

func SetVertexDrawMode(drawMode uint32) {
	vertexDrawMode = drawMode
}

func DrawArray(offset int, length int) {
	opengl.DrawArrays(vertexDrawMode, int32(offset), int32(length))
}

func NewMesh(vertices []float32) (mesh *VertexObject) {
	mesh = &VertexObject{
		numAttributes: 0,
	}
	vbo := uint32(0)
	opengl.GenBuffers(1, &vbo)
	opengl.BindBuffer(opengl.ARRAY_BUFFER, vbo)
	opengl.BufferData(opengl.ARRAY_BUFFER, 4*len(vertices), opengl.Ptr(vertices), opengl.STATIC_DRAW)
	mesh.Id = VertexBufferObject(vbo)

	// vao
	vao := uint32(0)
	opengl.GenVertexArrays(1, &vao)
	opengl.BindVertexArray(vao)
	opengl.EnableVertexAttribArray(0)
	opengl.BindBuffer(opengl.ARRAY_BUFFER, vbo)
	opengl.VertexAttribPointer(0, 3, opengl.FLOAT, false, 0, nil)

	mesh.AttribId[mesh.numAttributes] = VertexArrayObject(vao)
	mesh.attribStride[mesh.numAttributes] = 3

	mesh.numAttributes++

	return mesh
}

func BindMesh(mesh *VertexObject) {
	BindVertexAttributes(mesh)
}

// CreateVertexAttribute creates a vertex array attribute
func CreateVertexAttribute(mesh *VertexObject, bufferData []float32, stride int) {
	buffer := uint32(0)
	opengl.GenBuffers(1, &buffer)
	opengl.BindBuffer(opengl.ARRAY_BUFFER, buffer)
	opengl.BufferData(opengl.ARRAY_BUFFER, len(bufferData)*4, opengl.Ptr(bufferData), opengl.STATIC_DRAW)

	mesh.AttribId[mesh.numAttributes] = VertexArrayObject(buffer)
	mesh.attribStride[mesh.numAttributes] = stride
	mesh.numAttributes++
}

func CreateVertexAttributeArrayBuffer(mesh *VertexObject, bufferData []float32, stride int) {
	CreateVertexAttribute(mesh, bufferData, stride)
}

func BindVertexAttributes(mesh *VertexObject) {
	opengl.EnableVertexAttribArray(0)
	opengl.BindVertexArray(uint32(mesh.AttribId[0]))

	for i := 1; i < mesh.numAttributes; i++ {
		opengl.EnableVertexAttribArray(uint32(i))
		opengl.BindBuffer(opengl.ARRAY_BUFFER, uint32(mesh.AttribId[i]))
		opengl.VertexAttribPointer(uint32(i), int32(mesh.attribStride[i]), opengl.FLOAT, false, 0, nil)
	}
}

func FinishMesh() {
	opengl.BindVertexArray(0)
}

func DeleteMesh(mesh *VertexObject) {
	vao := uint32(mesh.AttribId[0])
	opengl.DeleteVertexArrays(1, &vao)

	for i := 1; i < mesh.numAttributes; i++ {
		buf := uint32(mesh.AttribId[1])
		opengl.DeleteBuffers(1, &buf)
	}



}