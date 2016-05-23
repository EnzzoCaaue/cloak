package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type RegisterController struct {
	*BaseController
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
func (base *BaseController) Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
func (base *BaseController) CreateAccount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
	account.Account.Premdays = util.Parser.Register.Premdays
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
	player.Level = util.Parser.Register.Level
	player.Health = util.Parser.Register.Health
	player.HealthMax = util.Parser.Register.Healthmax
	player.Mana = util.Parser.Register.Mana
	player.ManaMax = util.Parser.Register.Manamax
	player.Vocation = util.Vocation(form.CharacterVocation)
	player.Gender = util.Gender(form.CharacterSex)
	if player.Gender == 0 { // female
		player.LookBody = util.Parser.Register.Female.Lookbody
		player.LookFeet = util.Parser.Register.Female.Lookfeet
		player.LookHead = util.Parser.Register.Female.Lookhead
		player.LookType = util.Parser.Register.Female.Looktype
		player.LookAddons = util.Parser.Register.Female.Lookaddons
	} else {
		player.LookBody = util.Parser.Register.Male.Lookbody
		player.LookFeet = util.Parser.Register.Male.Lookfeet
		player.LookHead = util.Parser.Register.Male.Lookhead
		player.LookType = util.Parser.Register.Male.Looktype
		player.LookAddons = util.Parser.Register.Male.Lookaddons
	}
	player.Town = town.Get()
	player.Stamina = util.Parser.Register.Stamina
	player.SkillAxe = util.Parser.Register.Skills.Axe
	player.SkillSword = util.Parser.Register.Skills.Sword
	player.SkillClub = util.Parser.Register.Skills.Club
	player.SkillDist = util.Parser.Register.Skills.Dist
	player.SkillFish = util.Parser.Register.Skills.Fish
	player.SkillFist = util.Parser.Register.Skills.Fist
	player.SkillShield = util.Parser.Register.Skills.Shield
	player.Experience = util.Parser.Register.Experience
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
