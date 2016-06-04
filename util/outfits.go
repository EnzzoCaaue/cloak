package util

import (
    "log"
    "bytes"
    "os"
    "image"
    "image/draw"
    "image/color"
    "image/png"
)

// Outfit renders a tibia outfit with the given looks
func Outfit() ([]byte, error) {
    path := "C:/Users/ragga/Desktop/outfits"
    dst, err := os.Open(path+"/1_1_1_3.png")
    if err != nil {
        return nil, err
    }
    defer dst.Close()
    dstImg, err := png.Decode(dst)
    if err != nil {
        return nil, err
    }
    tpl, err := os.Open(path+"/1_1_1_3_template.png")
    if err != nil {
        return nil, err
    }
    defer tpl.Close()
    tplImg, err := png.Decode(tpl)
    if err != nil {
        return nil, err
    }
    colorize(tplImg)
    out := drawOutfitBase(dstImg, tplImg)
    output := bytes.Buffer{}
    err = png.Encode(&output, out)
    if err != nil {
        return nil, err
    }
    return output.Bytes(), nil
}

func colorize(img image.Image) {
    b := image.NewRGBA(img.Bounds())
    draw.Draw(b, b.Bounds(), img, image.ZP, draw.Src)
    for x := 0; x < b.Bounds().Dx(); x++ {
        for y := 0; y < b.Bounds().Dy(); y++ {
            log.Println(b.At(x, y).RGBA())
        }
    }
}

func drawOutfitBase(dst, tpl image.Image) *image.RGBA {
    outputRGBA := image.NewRGBA(dst.Bounds())
    draw.Draw(outputRGBA, dst.Bounds(), dst, image.ZP, draw.Src)
    mask := image.NewUniform(color.Alpha{110})
    draw.DrawMask(outputRGBA, tpl.Bounds(), tpl, image.ZP, mask, image.ZP, draw.Over)
    return outputRGBA
}