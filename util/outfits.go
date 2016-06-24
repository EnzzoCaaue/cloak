package util

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

var (
	red    = color.RGBA{255, 0, 0, 255}
	blue   = color.RGBA{0, 0, 255, 255}
	green  = color.RGBA{0, 255, 0, 255}
	yellow = color.RGBA{255, 255, 0, 255}
	colors = map[int]color.RGBA{
		0:   color.RGBA{255, 255, 255, 255},
		1:   color.RGBA{255, 212, 191, 255},
		2:   color.RGBA{255, 233, 191, 255},
		3:   color.RGBA{255, 255, 191, 255},
		4:   color.RGBA{233, 255, 191, 255},
		5:   color.RGBA{212, 255, 191, 255},
		6:   color.RGBA{191, 255, 191, 255},
		7:   color.RGBA{191, 255, 212, 255},
		8:   color.RGBA{191, 255, 233, 255},
		9:   color.RGBA{191, 255, 255, 255},
		10:  color.RGBA{191, 233, 255, 255},
		11:  color.RGBA{191, 212, 255, 255},
		12:  color.RGBA{191, 191, 255, 255},
		13:  color.RGBA{212, 191, 255, 255},
		14:  color.RGBA{233, 191, 255, 255},
		15:  color.RGBA{255, 191, 255, 255},
		16:  color.RGBA{255, 191, 233, 255},
		17:  color.RGBA{255, 191, 212, 255},
		18:  color.RGBA{255, 191, 191, 255},
		19:  color.RGBA{218, 218, 218, 255},
		20:  color.RGBA{191, 159, 143, 255},
		21:  color.RGBA{191, 175, 143, 255},
		22:  color.RGBA{191, 191, 143, 255},
		23:  color.RGBA{175, 191, 143, 255},
		24:  color.RGBA{159, 191, 143, 255},
		25:  color.RGBA{143, 191, 143, 255},
		26:  color.RGBA{143, 191, 159, 255},
		27:  color.RGBA{143, 191, 175, 255},
		28:  color.RGBA{143, 191, 191, 255},
		29:  color.RGBA{143, 175, 191, 255},
		30:  color.RGBA{143, 159, 191, 255},
		31:  color.RGBA{143, 143, 191, 255},
		32:  color.RGBA{159, 143, 191, 255},
		33:  color.RGBA{175, 143, 191, 255},
		34:  color.RGBA{191, 143, 191, 255},
		35:  color.RGBA{191, 143, 175, 255},
		36:  color.RGBA{191, 143, 159, 255},
		37:  color.RGBA{191, 143, 143, 255},
		38:  color.RGBA{182, 182, 182, 255},
		39:  color.RGBA{191, 127, 95, 255},
		40:  color.RGBA{191, 159, 95, 255},
		41:  color.RGBA{191, 191, 95, 255},
		42:  color.RGBA{159, 191, 95, 255},
		43:  color.RGBA{127, 191, 95, 255},
		44:  color.RGBA{95, 191, 95, 255},
		45:  color.RGBA{95, 191, 127, 255},
		46:  color.RGBA{95, 191, 159, 255},
		47:  color.RGBA{95, 191, 191, 255},
		48:  color.RGBA{95, 159, 191, 255},
		49:  color.RGBA{95, 127, 191, 255},
		50:  color.RGBA{95, 95, 191, 255},
		51:  color.RGBA{127, 95, 191, 255},
		52:  color.RGBA{159, 95, 191, 255},
		53:  color.RGBA{191, 95, 191, 255},
		54:  color.RGBA{191, 95, 159, 255},
		55:  color.RGBA{191, 95, 127, 255},
		56:  color.RGBA{191, 95, 95, 255},
		57:  color.RGBA{145, 145, 145, 255},
		58:  color.RGBA{191, 106, 63, 255},
		59:  color.RGBA{191, 148, 63, 255},
		60:  color.RGBA{191, 191, 63, 255},
		61:  color.RGBA{148, 191, 63, 255},
		62:  color.RGBA{106, 191, 63, 255},
		63:  color.RGBA{63, 191, 63, 255},
		64:  color.RGBA{63, 191, 106, 255},
		65:  color.RGBA{63, 191, 148, 255},
		66:  color.RGBA{63, 191, 191, 255},
		67:  color.RGBA{63, 148, 191, 255},
		68:  color.RGBA{63, 106, 191, 255},
		69:  color.RGBA{63, 63, 191, 255},
		70:  color.RGBA{106, 63, 191, 255},
		71:  color.RGBA{148, 63, 191, 255},
		72:  color.RGBA{191, 63, 191, 255},
		73:  color.RGBA{191, 63, 148, 255},
		74:  color.RGBA{191, 63, 106, 255},
		75:  color.RGBA{191, 63, 63, 255},
		76:  color.RGBA{109, 109, 109, 255},
		77:  color.RGBA{255, 85, 0, 255},
		78:  color.RGBA{255, 170, 0, 255},
		79:  color.RGBA{255, 255, 0, 255},
		80:  color.RGBA{169, 255, 0, 255},
		81:  color.RGBA{84, 255, 0, 255},
		82:  color.RGBA{0, 255, 0, 255},
		83:  color.RGBA{0, 255, 85, 255},
		84:  color.RGBA{0, 255, 170, 255},
		85:  color.RGBA{0, 255, 255, 255},
		86:  color.RGBA{0, 169, 255, 255},
		87:  color.RGBA{0, 85, 255, 255},
		88:  color.RGBA{0, 0, 255, 255},
		89:  color.RGBA{84, 0, 255, 255},
		90:  color.RGBA{170, 0, 255, 255},
		91:  color.RGBA{254, 0, 255, 255},
		92:  color.RGBA{255, 0, 169, 255},
		93:  color.RGBA{255, 0, 85, 255},
		94:  color.RGBA{255, 0, 0, 255},
		95:  color.RGBA{72, 72, 72, 255},
		96:  color.RGBA{191, 63, 0, 255},
		97:  color.RGBA{191, 127, 0, 255},
		98:  color.RGBA{191, 191, 0, 255},
		99:  color.RGBA{127, 191, 0, 255},
		100: color.RGBA{63, 191, 0, 255},
		101: color.RGBA{0, 191, 0, 255},
		102: color.RGBA{0, 191, 63, 255},
		103: color.RGBA{0, 191, 127, 255},
		104: color.RGBA{0, 191, 191, 255},
		105: color.RGBA{0, 127, 191, 255},
		106: color.RGBA{0, 63, 191, 255},
		107: color.RGBA{0, 0, 191, 255},
		108: color.RGBA{63, 0, 191, 255},
		109: color.RGBA{127, 0, 191, 255},
		110: color.RGBA{191, 0, 191, 255},
		111: color.RGBA{191, 0, 127, 255},
		112: color.RGBA{191, 0, 63, 255},
		113: color.RGBA{191, 0, 0, 255},
		114: color.RGBA{36, 36, 36, 255},
		115: color.RGBA{127, 42, 0, 255},
		116: color.RGBA{127, 85, 0, 255},
		117: color.RGBA{127, 127, 0, 255},
		118: color.RGBA{84, 127, 0, 255},
		119: color.RGBA{42, 127, 0, 255},
		120: color.RGBA{0, 127, 0, 255},
		121: color.RGBA{0, 127, 42, 255},
		122: color.RGBA{0, 127, 85, 255},
		123: color.RGBA{0, 127, 127, 255},
		124: color.RGBA{0, 84, 127, 255},
		125: color.RGBA{0, 42, 127, 255},
		126: color.RGBA{0, 0, 127, 255},
		127: color.RGBA{42, 0, 127, 255},
		128: color.RGBA{85, 0, 127, 255},
		129: color.RGBA{127, 0, 127, 255},
		130: color.RGBA{127, 0, 84, 255},
		131: color.RGBA{127, 0, 42, 255},
		132: color.RGBA{127, 0, 0, 255},
	}
)

// Outfit renders a tibia outfit with the given looks
func Outfit(path string, looktype, lookhead, lookbody, looklegs, lookfeet, lookaddons int) ([]byte, error) {
	dst, err := os.Open(fmt.Sprintf("%v/outfits/%v.png", path, looktype))
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	dstImg, err := png.Decode(dst)
	if err != nil {
		return nil, err
	}
	tpl, err := os.Open(fmt.Sprintf("%v/outfits/%v_template.png", path, looktype))
	if err != nil {
		return nil, err
	}
	defer tpl.Close()
	tplImg, err := png.Decode(tpl)
	if err != nil {
		return nil, err
	}
	tplImg = colorize(tplImg, lookhead, lookbody, looklegs, lookfeet)
	out := drawOutfitBase(dstImg, tplImg)
	output := bytes.Buffer{}
	err = png.Encode(&output, out)
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func colorize(img image.Image, lookhead, lookbody, looklegs, lookfeet int) image.Image {
	b := image.NewRGBA(img.Bounds())
	draw.Draw(b, b.Bounds(), img, image.ZP, draw.Src)
	for x := 0; x < b.Bounds().Dx(); x++ {
		for y := 0; y < b.Bounds().Dy(); y++ {
			colorRGBA := color.RGBAModel.Convert(b.At(x, y)).(color.RGBA)
			if colorRGBA == red {
				b.Set(x, y, colors[lookbody])
			} else if colorRGBA == green {
				b.Set(x, y, colors[looklegs])
			} else if colorRGBA == blue {
				b.Set(x, y, colors[lookfeet])
			} else if colorRGBA == yellow {
				b.Set(x, y, colors[lookhead])
			}
		}
	}
	return b
}

func drawOutfitBase(dst, tpl image.Image) *image.RGBA {
	outputRGBA := image.NewRGBA(dst.Bounds())
	draw.Draw(outputRGBA, dst.Bounds(), dst, image.ZP, draw.Src)
	mask := image.NewUniform(color.Alpha{140})
	draw.DrawMask(outputRGBA, tpl.Bounds(), tpl, image.ZP, mask, image.ZP, draw.Over)
	return outputRGBA
}
