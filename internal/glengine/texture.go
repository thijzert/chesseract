package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/thijzert/chesseract/internal/assets"
)

type texture uint32

type material struct {
	Texture      texture
	NormalMap    texture
	SpecularMap  texture
	ShineDamper  float32
	Reflectivity float32
	Transparency bool
}

var identityTextures [3]*image.RGBA

func init() {
	defaultColours := []color.Color{
		color.RGBA{96, 96, 96, 255},
		color.RGBA{127, 127, 255, 255},
		color.RGBA{0, 0, 255, 255},
	}

	for i := range defaultColours {
		identityTextures[i] = image.NewRGBA(image.Rect(0, 0, 16, 16))
	}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			for i, c := range defaultColours {
				identityTextures[i].Set(x, y, c)
			}
		}
	}
}

func (eng *Engine) loadMaterialAsset(name string) material {
	var rv material
	var err error

	rv.Texture, err = eng.loadTexture(name + ".diffuse.jpeg")
	if err != nil {
		rv.Texture, _ = eng.loadRGBATexture(identityTextures[0])
	}

	rv.NormalMap, err = eng.loadTexture(name + ".normal.jpeg")
	if err != nil {
		rv.NormalMap, _ = eng.loadRGBATexture(identityTextures[1])
	}

	rv.SpecularMap, err = eng.loadTexture(name + ".specular.jpeg")
	if err != nil {
		rv.SpecularMap, _ = eng.loadRGBATexture(identityTextures[2])
	}

	return rv
}

func (eng *Engine) loadTexture(name string) (texture, error) {
	imgBytes, err := assets.GetAsset("textures/" + name)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", name, err)
	}
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return eng.loadRGBATexture(rgba)
}

func (eng *Engine) loadRGBATexture(rgba *image.RGBA) (texture, error) {
	var rv uint32
	gl.GenTextures(1, &rv)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, rv)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture(rv), nil
}
