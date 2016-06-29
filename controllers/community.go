package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"fmt"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
)

type CommunityController struct {
	*pigo.Controller
}

// Highscores process and shows the highscores page
func (base *CommunityController) Highscores(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	highscoreType := p.ByName("type")
	page, err := strconv.Atoi(p.ByName("page"))
	if err != nil {
		base.Redirect = "/"
		return
	}
	pageIndex, query, name := util.GetHighscoreQuery(page, highscoreType, 10)
	list, err := models.GetHighscores(pageIndex, query)
	if err != nil {
		base.Error = "Error while getting highscore list"
		return
	}
	if len(list) == 0 && page > 0 {
		base.Redirect = fmt.Sprintf("/highscores/%v/%v",
			highscoreType,
			page-1,
		)
		return
	}
	base.Data["PageNext"] = page + 1
	base.Data["PageOld"] = page - 1
	base.Data["SkillName"] = name
	base.Data["Skill"] = highscoreType
	base.Data["List"] = list
	base.Data["CurrentRank"] = pageIndex
	base.Data["OldRank"] = pageIndex - 10
	base.Data["NextRank"] = pageIndex + 10
	base.Template = "highscores.html"
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
		return
	}
	w.Write(outfit)
	w.Header().Set("Content-type", "image/png")
}

// ServerOverview shows all the server information
func (base *CommunityController) ServerOverview(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	base.Template = "server_overview.html"
	base.Data["Name"] = util.Config.String("serverName")
	base.Data["WorldType"] = util.Config.String("worldType")
	base.Data["ProtectionLevel"] = util.Config.Int("protectionLevel")
	base.Data["RedSkull"] = util.Config.Int("killsToRedSkull")
	base.Data["BlackSkull"] = util.Config.Int("killsToBlackSkull")
	base.Data["InfiniteRunes"] = util.Config.Bool("removeChargesFromRunes")
	base.Data["MagicRate"] = util.Config.Int("rateMagic")
	base.Data["LootRate"] = util.Config.Int("rateLoot")
	base.Data["SkillRate"] = util.Config.Int("rateSkill")
	base.Data["SpawnRate"] = util.Config.Int("rateSpawn")
	base.Data["Motd"] = util.Config.String("motd")
	base.Data["FreePremium"] = util.Config.Bool("freePremium")
	base.Data["StagesEnabled"] = util.Stages.IsEnabled()
	base.Data["Stages"] = util.Stages.GetAll()
}
