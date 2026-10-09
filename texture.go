package gosigl

import (
	opengl "github.com/go-gl/gl/v4.1-core/gl"
)

type TextureBindingId uint32
type TextureSlotId uint32
type PixelFormat uint32

const RGBA = PixelFormat(opengl.RGBA)
const RGB = PixelFormat(opengl.RGB)
const BGR = PixelFormat(opengl.BGR)
const BGRA = PixelFormat(opengl.BGRA)
const DXT1 = PixelFormat(opengl.COMPRESSED_RGB_S3TC_DXT1_EXT)
const DXT1A = PixelFormat(opengl.COMPRESSED_RGBA_S3TC_DXT1_EXT)
const DXT3 = PixelFormat(opengl.COMPRESSED_RGBA_S3TC_DXT3_EXT)
const DXT5 = PixelFormat(opengl.COMPRESSED_RGBA_S3TC_DXT5_EXT)

var textureSlots = [...]TextureSlotId{
	opengl.TEXTURE0,
	opengl.TEXTURE1,
	opengl.TEXTURE2,
	opengl.TEXTURE3,
	opengl.TEXTURE4,
	opengl.TEXTURE5,
	opengl.TEXTURE6,
	opengl.TEXTURE7,
	opengl.TEXTURE8,
	opengl.TEXTURE9,
	opengl.TEXTURE10,
	opengl.TEXTURE11,
	opengl.TEXTURE12,
	opengl.TEXTURE13,
	opengl.TEXTURE14,
	opengl.TEXTURE15,
	opengl.TEXTURE16,
	opengl.TEXTURE17,
	opengl.TEXTURE18,
	opengl.TEXTURE19,
	opengl.TEXTURE20,
	opengl.TEXTURE21,
	opengl.TEXTURE22,
	opengl.TEXTURE23,
	opengl.TEXTURE24,
	opengl.TEXTURE25,
	opengl.TEXTURE26,
	opengl.TEXTURE27,
	opengl.TEXTURE28,
	opengl.TEXTURE29,
	opengl.TEXTURE30,
	opengl.TEXTURE31,
}

var cubemapSlots = [...]TextureSlotId{
	opengl.TEXTURE_CUBE_MAP_POSITIVE_X,
	opengl.TEXTURE_CUBE_MAP_NEGATIVE_X,
	opengl.TEXTURE_CUBE_MAP_POSITIVE_Y,
	opengl.TEXTURE_CUBE_MAP_NEGATIVE_Y,
	opengl.TEXTURE_CUBE_MAP_POSITIVE_Z,
	opengl.TEXTURE_CUBE_MAP_NEGATIVE_Z,
}

func TextureSlot(index int) TextureSlotId {
	if index > 31 {
		return textureSlots[0]
	}
	return textureSlots[index]
}

func BindTexture2D(slot TextureSlotId, id TextureBindingId) {
	opengl.ActiveTexture(uint32(slot))
	opengl.BindTexture(opengl.TEXTURE_2D, uint32(id))
}

func BindTextureCubemap(slot TextureSlotId, id TextureBindingId) {
	opengl.ActiveTexture(uint32(slot))
	opengl.BindTexture(opengl.TEXTURE_CUBE_MAP, uint32(id))
}

func CreateTexture2D(slot TextureSlotId, width int, height int, pixelData []byte, pixelFormat PixelFormat, clampToEdge bool) TextureBindingId {
	return CreateTexture2DMipmaps(slot, width, height, [][]byte{pixelData}, pixelFormat, clampToEdge)
}

// CreateTexture2DMipmaps creates a texture from each of its mipmaps, from the largest, width x height, to the smallest.
// Each is half the size of the one before, to no smaller than 1x1. A texture with more than one mipmap is filtered
// trilinearly between them.
func CreateTexture2DMipmaps(slot TextureSlotId, width int, height int, mipmaps [][]byte, pixelFormat PixelFormat, clampToEdge bool) TextureBindingId {
	textureBuffer := uint32(0)
	opengl.GenTextures(1, &textureBuffer)
	BindTexture2D(slot, TextureBindingId(textureBuffer))

	setTextureParameters(opengl.TEXTURE_2D, len(mipmaps), clampToEdge)
	for level, pixelData := range mipmaps {
		uploadTextureImage(opengl.TEXTURE_2D, level, mipmapSize(width, level), mipmapSize(height, level), pixelData, pixelFormat)
	}

	return TextureBindingId(textureBuffer)
}

func CreateTextureCubemap(slot TextureSlotId, width int, height int, pixelData [6][]byte, pixelFormat PixelFormat, clampToEdge bool) TextureBindingId {
	textureBuffer := uint32(0)
	opengl.GenTextures(1, &textureBuffer)
	BindTextureCubemap(slot, TextureBindingId(textureBuffer))

	setTextureParameters(opengl.TEXTURE_CUBE_MAP, 1, clampToEdge)
	for face, slot := range cubemapSlots {
		uploadTextureImage(uint32(slot), 0, width, height, pixelData[face], pixelFormat)
	}

	return TextureBindingId(textureBuffer)
}

// mipmapSize returns the size of a mipmap level of a texture of size, which halves each level to no smaller than 1
func mipmapSize(size int, level int) int {
	if size>>level < 1 {
		return 1
	}
	return size >> level
}

// maxTextureAnisotropy is the most samples anisotropic filtering takes, or 0 before it is known
var maxTextureAnisotropy float32

// maxAnisotropy returns the most samples anisotropic filtering takes: as many as the hardware supports, up to 16. It is
// 1, no anisotropic filtering, without support for it.
func maxAnisotropy() float32 {
	if maxTextureAnisotropy == 0 {
		maxTextureAnisotropy = 1
		supported := float32(0)
		opengl.GetFloatv(opengl.MAX_TEXTURE_MAX_ANISOTROPY, &supported)
		if opengl.GetError() == opengl.NO_ERROR && supported > 1 {
			maxTextureAnisotropy = supported
			if maxTextureAnisotropy > 16 {
				maxTextureAnisotropy = 16
			}
		}
	}
	return maxTextureAnisotropy
}

// setTextureParameters sets how the bound texture of textureType, which has the given number of mipmaps, is sampled
func setTextureParameters(textureType uint32, mipmaps int, clampToEdge bool) {
	if mipmaps > 1 {
		opengl.TexParameteri(textureType, opengl.TEXTURE_MIN_FILTER, opengl.LINEAR_MIPMAP_LINEAR)
		// Sample more of a texture seen at an angle, so that it is not blurred along it
		if anisotropy := maxAnisotropy(); anisotropy > 1 {
			opengl.TexParameterf(textureType, opengl.TEXTURE_MAX_ANISOTROPY, anisotropy)
		}
	} else {
		opengl.TexParameteri(textureType, opengl.TEXTURE_MIN_FILTER, opengl.LINEAR)
	}
	opengl.TexParameteri(textureType, opengl.TEXTURE_MAG_FILTER, opengl.LINEAR)
	// A texture is only complete with every mipmap up to its max level
	opengl.TexParameteri(textureType, opengl.TEXTURE_MAX_LEVEL, int32(mipmaps-1))

	if clampToEdge == true {
		opengl.TexParameteri(textureType, opengl.TEXTURE_WRAP_S, opengl.CLAMP_TO_EDGE)
		opengl.TexParameteri(textureType, opengl.TEXTURE_WRAP_T, opengl.CLAMP_TO_EDGE)
	} else {
		opengl.TexParameteri(textureType, opengl.TEXTURE_WRAP_S, opengl.REPEAT)
		opengl.TexParameteri(textureType, opengl.TEXTURE_WRAP_T, opengl.REPEAT)
	}
}

// uploadTextureImage uploads one mipmap level of the bound texture's image for target
func uploadTextureImage(target uint32, level int, width int, height int, pixelData []byte, pixelFormat PixelFormat) {
	if isPixelFormatCompressed(pixelFormat) {
		opengl.CompressedTexImage2D(
			target,
			int32(level),
			uint32(pixelFormat),
			int32(width),
			int32(height),
			0,
			int32(len(pixelData)),
			opengl.Ptr(pixelData))
	} else {
		// Rows are packed tightly. The rows of small RGB mipmaps are not a multiple of the default 4 bytes long.
		opengl.PixelStorei(opengl.UNPACK_ALIGNMENT, 1)
		opengl.TexImage2D(
			target,
			int32(level),
			opengl.RGBA,
			int32(width),
			int32(height),
			0,
			uint32(pixelFormat),
			opengl.UNSIGNED_BYTE,
			opengl.Ptr(pixelData))
	}
}

func DeleteTextures(ids ...TextureBindingId) {
	rawIds := make([]uint32, len(ids))
	for idx, id := range ids {
		rawIds[idx] = uint32(id)
	}
	opengl.DeleteTextures(int32(len(rawIds)), &rawIds[0])
}

func isPixelFormatCompressed(format PixelFormat) bool {
	switch format {
	case DXT1, DXT1A, DXT3, DXT5:
		return true
	}

	return false
}