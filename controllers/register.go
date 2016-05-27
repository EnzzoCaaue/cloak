package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/raggaer/pigo"
)

type RegisterController struct {
	*pigo.Controller
}

// RegisterForm saves the register form
type RegisterForm struct {
	AccountName       string `validate:"regexp=^[A-Za-z0-9]+$, min=5, max=20" alias:"Account name"`
	Email             string `alias:"Account email"`
	Password          string `validate:"min=8, max=30" alias:"Account password"`
	CharacterName     string `validate:"regexp=^[A-Z a-z]+$, max=14" alias:"Character name"`
	CharacterSex      string `validate:"validGender" alias:"Character gender"`
	CharacterVocation string `validate:"validVocation" alias:"Character vocation"`
	CharacterTown     string
	Captcha           string `validate:"validCaptcha" alias:"Captcha check"`
}

// Register shows the register.html page
func (base *RegisterController) Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	towns, err := models.GetTowns()
	if err != nil {
		base.Error = "Error fetching town list"
		return
	}
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Towns"] = towns
	base.Template = "register.html"
}

// CreateAccount process register.html and creates an account
func (base *RegisterController) CreateAccount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	form := &RegisterForm{
		req.FormValue("accountname"),
		req.FormValue("email"),
		req.FormValue("password1"),
		req.FormValue("name"),
		req.FormValue("sex"),
		req.FormValue("vocation"),
		req.FormValue("town"),
		req.FormValue("g-recaptcha-response"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for _, v := range errs {
			base.Session.AddFlash(v.Error(), "errors")
		}
		base.Redirect = "/account/create"
		return
	}
	town := models.NewTown(form.CharacterTown)
	if !town.Exists() {
		base.Session.AddFlash("Invalid character town", "errors")
		base.Redirect = "/account/create"
		return
	}
	account := models.NewAccount()
	account.Account.Name = form.AccountName
	if account.NameExists() {
		base.Session.AddFlash("Account name is already in use", "errors")
		base.Redirect = "/account/create"
		return
	}
	if account.EmailExists() {
		base.Session.AddFlash("Account email is already in use", "errors")
		base.Redirect = "/account/create"
		return
	}
	account.Account.Password = fmt.Sprintf("%x", sha1.Sum([]byte(form.Password)))
	account.Account.Premdays = pigo.Config.Key("Register").Int("Premdays")
	account.Account.Email = form.Email
	err := account.Account.Save()
	if err != nil {
		base.Error = "Error while saving your account"
		return
	}
	err = account.Save()
	if err != nil {
		base.Error = "Error while saving your account"
		return
	}
	player := models.NewPlayer()
	player.AccountID = account.Account.ID
	player.Name = form.CharacterName
	if player.Exists() {
		base.Session.AddFlash("Character name is already in use", "errors")
		base.Redirect = "/account/create"
		return
	}
	player.Level = pigo.Config.Key("register").Int("level")
	player.Health = pigo.Config.Key("register").Int("health")
	player.HealthMax = pigo.Config.Key("register").Int("healthmax")
	player.Mana = pigo.Config.Key("register").Int("mana")
	player.ManaMax = pigo.Config.Key("register").Int("manamax")
	player.Vocation = util.Vocation(form.CharacterVocation)
	player.Gender = util.Gender(form.CharacterSex)
	if player.Gender == 0 { // female
		player.LookBody = pigo.Config.Key("register").Key("female").Int("lookbody")
		player.LookFeet = pigo.Config.Key("register").Key("female").Int("lookfeet")
		player.LookHead = pigo.Config.Key("register").Key("female").Int("lookhead")
		player.LookType = pigo.Config.Key("register").Key("female").Int("looktype")
		player.LookAddons = pigo.Config.Key("register").Key("female").Int("lookaddons")
	} else {
		player.LookBody = pigo.Config.Key("register").Key("male").Int("lookbody")
		player.LookFeet = pigo.Config.Key("register").Key("male").Int("lookfeet")
		player.LookHead = pigo.Config.Key("register").Key("male").Int("lookhead")
		player.LookType = pigo.Config.Key("register").Key("male").Int("looktype")
		player.LookAddons = pigo.Config.Key("register").Key("male").Int("lookaddons")
	}
	player.Town = town.Get()
	player.Stamina = pigo.Config.Key("register").Int("stamina")
	player.SkillAxe = pigo.Config.Key("register").Key("skills").Int("axe")
	player.SkillSword = pigo.Config.Key("register").Key("skills").Int("sword")
	player.SkillClub = pigo.Config.Key("register").Key("skills").Int("club")
	player.SkillDist = pigo.Config.Key("register").Key("skills").Int("dist")
	player.SkillFish = pigo.Config.Key("register").Key("skills").Int("fish")
	player.SkillFist = pigo.Config.Key("register").Key("skills").Int("fist")
	player.SkillShield = pigo.Config.Key("register").Key("skills").Int("shield")
	player.Experience = pigo.Config.Key("register").Int("experience")
	err = player.Save()
	if err != nil {
		base.Error = "Error while saving player"
		return
	}
	key := uniuri.New()
	for models.TokenExists(key) {
		key = uniuri.New()
	}
	err = account.UpdateToken(key)
	if err != nil {
		base.Error = "Error while updating your acocunt token"
		return
	}
	base.Session.Set("key", key)
	base.Session.Set("logged", 1)
	base.Redirect = "/account/manage"
}
