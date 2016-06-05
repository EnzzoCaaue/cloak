package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"log"
	"time"

	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
)

type CommunityController struct {
	*pigo.Controller
}

// CharacterView shows a character
func (base *CommunityController) CharacterView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name, err := url.QueryUnescape(p.ByName("name"))
	if err != nil {
		base.Error = "Invalid character name"
		return
	}
	player := models.GetPlayerByName(name)
	if player == nil {
		base.Redirect = "/"
		return
	}
	player.GetGuild()
	deaths, err := player.GetDeaths()
	if err != nil {
		base.Error = "Error while getting character deaths"
		return
	}
	base.Data["Info"] = player
	base.Data["Deaths"] = deaths
	base.Template = "character_view.html"
}

// SignatureView shows a signature
func (base *CommunityController) SignatureView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name, err := url.QueryUnescape(p.ByName("name"))
	if err != nil {
		base.Error = "Invalid character name"
		return
	}
	player := models.GetPlayerByName(name)
	if player == nil {
		base.Error = "Unknown character name"
		return
	}
	signatureFile, err := os.Open(pigo.Config.String("template") + "/public/signatures/" + player.Name + ".png")
	if err != nil { // No signature
		signature, err := util.CreateSignature(player.Name, player.Gender, player.Vocation, player.Level, player.LastLogin)
		if err != nil {
			base.Error = "Error while creating signature"
			return
		}
		w.Header().Set("Content-type", "image/png")
		w.Write(signature)
		return
	}
	defer signatureFile.Close()
	signatureFileStats, err := signatureFile.Stat()
	if err != nil {
		base.Error = "Error while reading signature stats"
		return
	}
	if signatureFileStats.ModTime().Unix()+(1*60) > time.Now().Unix() {
		buffer, err := ioutil.ReadAll(signatureFile)
		if err != nil {
			base.Error = "Error while reading signature bytes"
			return
		}
		w.Header().Set("Content-type", "image/png")
		w.Write(buffer)
		return
	}
	signature, err := util.CreateSignature(player.Name, player.Gender, player.Vocation, player.Level, player.LastLogin)
	if err != nil {
		base.Error = "Error while creating signature"
		return
	}
	w.Header().Set("Content-type", "image/png")
	w.Write(signature)
}

// SearchCharacter searchs for names LIKE
func (base *CommunityController) SearchCharacter(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	players, err := models.SearchPlayers(req.FormValue("name"))
	if err != nil {
		base.Error = "Error while searching for players"
		return
	}
	base.Data["Current"] = req.FormValue("name")
	base.Data["Characters"] = players
	base.Template = "character_search.html"
}

// OutfitView shows a player outfit
func (base *CommunityController) OutfitView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	playerName := p.ByName("name")
	player := models.GetPlayerByName(playerName)
	if player == nil {
		return
	}
	log.Println(player.LookBody)
	outfit, err := util.Outfit(pigo.Config.String("template"), player.LookType, player.LookHead, player.LookBody, player.LookLegs, player.LookFeet, player.LookAddons)
	if err != nil {
		base.Error = err.Error()
		log.Println(err)
		return
	}
	w.Write(outfit)
	w.Header().Set("Content-type", "image/png")
}