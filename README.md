# Gosigl

Gosigl (Go SImple openGL) is an somewhat opinionated Go wrapper for OpenGL with emphasis on simple. It's purpose is to make creating simple OpenGL programs really goddamn
easy. It doesn't provide a whole lot of functionality, but it's really quick and simple to get a basic renderable
world up and running.

It provides simplified functions for creating textures, VAO and VBOs, and shader objects.


#### Texture
Provides functions for creating and binding Texture2D and Cube_Map_Textures.
Creating and binding a Texture2D can be done as follows:
```go
buf := gosigl.CreateTexture2D(
		gosigl.TextureSlot(0),   // Assign texture slot
		width,                  // texture width
		height,                 // texture height
		pixelData,              // []byte raw colour data
		gosigl.RGB,              // colour data format
		false)                  // true = clamp to edge, false = repeat
gosigl.BindTexture2D(gosigl.TextureSlot(0), buf)
```
Cubemaps are very similar:
```go
buf := gosigl.CreateTextureCubemap(
	gosigl.TextureSlot(0),   // Assign texture slot
	width,                  // texture width
    height,                 // texture height
    pixelData,              // [6][]byte raw colour data
    gosigl.RGB,              // colour data format
	true)
gosigl.BindTextureCubemap(gosigl.TextureSlot(0), buf)
```

#### Mesh
Simple methods are available for VBO and VAO generation, as follows:
```go
mesh := gosigl.NewMesh(vertices) // vertices = []float32
gosigl.CreateVertexAttribute(mesh, uvs, 2) // uvs = []float32, 2 = numPerVertex
gosigl.CreateVertexAttribute(mesh, normals, 3)
gosigl.FinishMesh()
```

To draw:
```go
gosigl.BindMesh(mesh)            // bind vbo
gosigl.DrawArray(offset, length) // draw from bound
```

To delete:
```go
gosigl.DeleteMesh(mesh)
```

#### Shader
@Incomplete
