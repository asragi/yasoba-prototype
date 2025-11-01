//go:build ignore

//kage:unit pixels
package image

var Time float

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	tex := imageSrc0At(srcPos)

	speed := 100.0
	t := Time
	redTex := imageSrc0At(srcPos + vec2(0, speed*t*t))
	redTex.gb = vec2(0)
	redTex.a = 0
	timer := t * t
	reda := min(1.0, timer*4)
	return tex*(1.0-reda) + redTex*(1.0-reda)
}
