package util

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"github.com/yaimko/yaimko"
)

const (
	defaultSignature   = "signature"
	signaturePath      = "public/signatures"
	signatureExtension = "png"
	fontsPath          = "public/fonts"
	defaultFont        = "Aller_Bd"
	fontsExtension     = "ttf"
)

// CreateSignature creates a player signature image
func CreateSignature(name string, gender, vocation, level int, lastlogin int64) ([]byte, error) {
	background, err := os.Open(fmt.Sprintf("%v/%v/%v.%v", yaimko.Config.String("template.dir"), signaturePath, defaultSignature, signatureExtension))
	if err != nil {
		return nil, err
	}
	backgroundRGBA := image.NewRGBA(image.Rect(0, 0, 495, 134))
	backgroundDecoded, _, err := image.Decode(background)
	if err != nil {
		return nil, err
	}
	draw.Draw(backgroundRGBA, backgroundRGBA.Bounds(), backgroundDecoded, image.ZP, draw.Src)
	fontBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v/%v.%v", yaimko.Config.String("template.dir"), fontsPath, defaultFont, fontsExtension))
	if err != nil {
		return nil, err
	}
	signatureFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	signatureTextDrawer := freetype.NewContext()
	signatureTextDrawer.SetDPI(72)
	signatureTextDrawer.SetFont(signatureFont)
	signatureTextDrawer.SetFontSize(14)
	signatureTextDrawer.SetClip(backgroundRGBA.Bounds())
	signatureTextDrawer.SetDst(backgroundRGBA)
	signatureTextDrawer.SetSrc(image.Black)
	if _, err := signatureTextDrawer.DrawString(Config.String("serverName"), freetype.Pt(20, 30)); err != nil {
		return nil, err
	} else if _, err = signatureTextDrawer.DrawString("Name: "+name, freetype.Pt(20, 50)); err != nil {
		return nil, err
	} else if _, err = signatureTextDrawer.DrawString("Vocation: "+GetVocation(vocation), freetype.Pt(20, 70)); err != nil {
		return nil, err
	} else if _, err = signatureTextDrawer.DrawString("Level: "+strconv.Itoa(level), freetype.Pt(20, 90)); err != nil {
		return nil, err
	} else if _, err = signatureTextDrawer.DrawString("Last login: "+UnixToString(lastlogin), freetype.Pt(20, 110)); err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	err = png.Encode(buffer, backgroundRGBA)
	if err != nil {
		return nil, err
	}
	signatureFile, err := os.Create(fmt.Sprintf("%v/%v/%v.%v", yaimko.Config.String("template.dir"), signaturePath, name, signatureExtension))
	defer signatureFile.Close()
	signatureFile.Write(buffer.Bytes())
	return buffer.Bytes(), nil
}
