package util

import (
    "os"
    "github.com/golang/freetype"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
    "bytes"
    "io/ioutil"
)

// CreateSignature creates a player signature image
func CreateSignature(name string, gender, vocation, level int, lastlogin int64) ([]byte, error) {
    background, err := os.Open(Parser.Style.Template + "/public/images/signature.png")
    if err != nil {
        HandleError("Cannot open signature background image", err)
        return nil, err
    }
    backgroundRGBA := image.NewRGBA(image.Rect(0, 0, 495, 134))
    backgroundDecoded, _, err := image.Decode(background)
    if err != nil {
        HandleError("Cannot decode signature background image", err)
        return nil, err
    }
    draw.Draw(backgroundRGBA, backgroundRGBA.Bounds(), backgroundDecoded, image.ZP, draw.Src)
    fontBytes, err := ioutil.ReadFile(Parser.Style.Template + "/public/fonts/Aller_Bd.ttf")
    if err != nil {
        HandleError("Cannot open signature font file", err)
        return nil, err
    }
    signatureFont, err := freetype.ParseFont(fontBytes)
    if err != nil {
        HandleError("Cannot parse free type font", err)
        return nil, err
    }
    signatureTextDrawer := freetype.NewContext()
    signatureTextDrawer.SetDPI(72)
    signatureTextDrawer.SetFont(signatureFont)
    signatureTextDrawer.SetFontSize(14)
    signatureTextDrawer.SetClip(backgroundRGBA.Bounds())
    signatureTextDrawer.SetDst(backgroundRGBA)
    signatureTextDrawer.SetSrc(image.Black)
    _, err = signatureTextDrawer.DrawString("Name: " + name, freetype.Pt(20, 30))
    if err != nil {
        HandleError("Error while drawing name text", err)
        return nil, err
    }
    _, err = signatureTextDrawer.DrawString("Gender: " + GetGender(gender), freetype.Pt(20, 50))
    if err != nil {
        HandleError("Error while drawing gender text", err)
        return nil, err
    }
    _, err = signatureTextDrawer.DrawString("Vocation: " + GetVocation(vocation), freetype.Pt(20, 70))
    if err != nil {
        HandleError("Error while drawing vocation text", err)
        return nil, err
    }
    _, err = signatureTextDrawer.DrawString("Level: 1", freetype.Pt(20, 90))
    if err != nil {
        HandleError("Error while drawing level text", err)
        return nil, err
    }
    _, err = signatureTextDrawer.DrawString("Last login: " + UnixToString(lastlogin), freetype.Pt(20, 110))
    if err != nil {
        HandleError("Error while drawing last login text", err)
        return nil, err
    }
    buffer := &bytes.Buffer{}
    err = png.Encode(buffer, backgroundRGBA)
    if err != nil {
        HandleError("Error writing to output buffer", err)
        return nil, err
    }
    signatureFile, err := os.Create(Parser.Style.Template + "/public/signatures/" + name + ".png")
    defer signatureFile.Close()
    signatureFile.Write(buffer.Bytes())
    return buffer.Bytes(), nil   
}


/*
	
	
	
	
	buffer_signature := bytes.Buffer{}
	b := bufio.NewWriter(&buffer_signature)
	err = png.Encode(b, rgba)
	if err != nil {
		http.Error(res, "Error while encoding PNG file", 500)
		return
	}
	err = b.Flush()
	if err != nil {
		http.Error(res, "Error while writing the PNG file", 500)
		return
	}
	signature_out, err := os.Create(config.Parser.Style.Template + "/public/signatures/" + character_info.Name + ".png")
	if err != nil {
		http.Error(res, "Error while saving signature", 500)
		return
	}
	defer signature_out.Close()
	signature_out.Write(buffer_signature.Bytes())
	res.Header().Set("Content-type", "image/png")
	res.Write(buffer_signature.Bytes())*/