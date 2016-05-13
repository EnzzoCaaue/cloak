package controllers

import (
	"github.com/Cloakaac/cloak/models"
	//"github.com/Cloakaac/cloak/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
)

// CharacterView shows a character
func (base *BaseController) CharacterView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

}

// SignatureView shows a signature
func (base *BaseController) SignatureView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name, err := url.QueryUnescape(p.ByName("name"))
	if err != nil {
		http.Error(w, "Oops! Invalid character name", 500)
		return
	}
	player := models.GetPlayerByName(name)
    if player == nil {
        http.Error(w, "Oops! Unknown character name", 500)
	    return
    }
	signatureFile, err := os.Open(util.Parser.Template + "/public/signatures/" + player.Name + ".png")
	if err != nil { // No signature
		signature, err := util.CreateSignature(player.Name, player.Gender, player.Vocation, player.Level, player.LastLogin)
		if err != nil {
			http.Error(w, "Oops! Cannot create signature", 500)
			return
		}
		w.Header().Set("Content-type", "image/png")
		w.Write(signature)
		return
	}
	defer signatureFile.Close()
	signatureFileStats, err := signatureFile.Stat()
	if err != nil {
		util.HandleError("Cannot get signature file stats", err)
		http.Error(w, "Oops! Cannot read signature stats", 500)
		return
	}
	if signatureFileStats.ModTime().Unix()+(1*60) > time.Now().Unix() {
		buffer, err := ioutil.ReadAll(signatureFile)
		if err != nil {
			util.HandleError("Cannot get signature file bytes", err)
			http.Error(w, "Oops! Cannot read signature file", 500)
			return
		}
		w.Header().Set("Content-type", "image/png")
		w.Write(buffer)
		return
	}
	signature, err := util.CreateSignature(player.Name, player.Gender, player.Vocation, player.Level, player.LastLogin)
	if err != nil {
		http.Error(w, "Oops! Cannot create signature", 500)
		return
	}
	w.Header().Set("Content-type", "image/png")
	w.Write(signature)
	return
}
