# Gosgl

Gosigl (Go SImple openGL) is an somewhat opinionated Go wrapper for OpenGL with emphasis on simple. It's purpose is to make creating simple OpenGL programs really goddamn
easy. It doesn't provide a whole lot of functionality, but it's really quick and simple to get a basic renderable
world up and running.

It provides simplified functions for creating textures, VAO and VBOs, and shader objects.


#### Texture
Provides functions for creating and binding Texture2D and Cube_Map_Textures.
Creating and binding a Texture2D can be done as follows:
```go
buf := gosgl.CreateTexture2D(
		gosgl.TextureSlot(0),   // Assign texture slot
		width,                  // texture width
		height,                 // texture height
		pixelData,              // []byte raw colour data
		gosgl.RGB,              // colour data format
		false)                  // true = clamp to edge, false = repeat
gosgl.BindTexture2D(gosgl.TextureSlot(0), buf)
```
Cubemaps are very similar:
```go
buf := gosgl.CreateTextureCubemap(
	gosgl.TextureSlot(0),   // Assign texture slot
	width,                  // texture width
    height,                 // texture height
    pixelData,              // [6][]byte raw colour data
    gosgl.RGB,              // colour data format
	true)
gosgl.BindTextureCubemap(gosgl.TextureSlot(0), buf)
```

#### Mesh
Simple methods are available for VBO and VAO generation, as follows:
```go
mesh := gosgl.NewMesh(vertices) // vertices = []float32
gosgl.CreateVertexAttribute(mesh, uvs, 2) // uvs = []float32, 2 = numPerVertex
gosgl.CreateVertexAttribute(mesh, normals, 3)
gosgl.FinishMesh()
```

To draw:
```go
gosgl.BindMesh(mesh)            // bind vbo
gosgl.DrawArray(offset, length) // draw from bound
```

To delete:
```go
gosgl.DeleteMesh(mesh)
```

#### Shader
@Incomplete