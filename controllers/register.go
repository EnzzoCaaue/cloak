package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/dchest/uniuri"
	"crypto/sha1"
	"fmt"
)

type register struct {
	Errors []string
	Token  string
	Towns  []*models.Town
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
    Captcha string `validate:"validCaptcha" alias:"Captcha check"`
}

// Register shows the register.html page
func (base *BaseController) Register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	towns, err := models.GetTowns()
	if err != nil {
		util.HandleError("Error on models.GetTowns", err)
		http.Error(w, "Oops! Something wrong happened while getting town list", http.StatusBadRequest)
		return
	}
	csrfToken := uniuri.New()
	//base.Session.Set("token", csrfToken)
	response := &register{
		base.Session.GetFlashes("errors"),
		csrfToken,
		towns,
	}
	err = base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving your request session", http.StatusBadRequest)
		return

	}
	template.Renderer.ExecuteTemplate(w, "register.html", response)
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
            base.Session.Save(req, w)
        }
        http.Redirect(w, req, "/account/create", 301)
        return
	}
	town := models.NewTown(form.CharacterTown)
	if !town.Exists() {
		base.Session.AddFlash("Invalid character town", "errors")
        base.Session.Save(req, w)
		http.Redirect(w, req, "/account/create", 301)
        return
	}
    account := models.NewAccount()
	account.Account.Name = form.AccountName
	if account.NameExists() {
		base.Session.AddFlash("Account name is already in use", "errors")
        base.Session.Save(req, w)
		http.Redirect(w, req, "/account/create", 301)
        return
	}
	if account.EmailExists() {
		base.Session.AddFlash("Account email is already in use", "errors")
        base.Session.Save(req, w)
		http.Redirect(w, req, "/account/create", 301)
        return
	}
	account.Account.Password = fmt.Sprintf("%x", sha1.Sum([]byte(form.Password)))
	account.Account.Premdays = util.Parser.Register.Premdays
	account.Account.Email = form.Email
	err := account.Account.Save()
	if err != nil {
		util.HandleError("Error creating account", err)
		http.Error(w, "Oops! Something wrong happened while creating your account", http.StatusBadRequest)
		return
	}
	err = account.Save()
	if err != nil {
		util.HandleError("Error creating account", err)
		http.Error(w, "Oops! Something wrong happened while creating your cloaka account", http.StatusBadRequest)
		return
	}
	player := models.NewPlayer()
	player.AccountID = account.Account.ID
	player.Name = form.CharacterName
	if player.Exists() {
		base.Session.AddFlash("Character name is already in use", "errors")
        base.Session.Save(req, w)
		http.Redirect(w, req, "/account/create", 301)
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
		player.LookBody = util.Parser.Femalelooktype.Lookbody
		player.LookFeet = util.Parser.Femalelooktype.Lookfeet
		player.LookHead = util.Parser.Femalelooktype.Lookhead
		player.LookType = util.Parser.Femalelooktype.Looktype
		player.LookAddons = util.Parser.Femalelooktype.Lookaddons
	} else {
		player.LookBody = util.Parser.Malelooktype.Lookbody
		player.LookFeet = util.Parser.Malelooktype.Lookfeet
		player.LookHead = util.Parser.Malelooktype.Lookhead
		player.LookType = util.Parser.Malelooktype.Looktype
		player.LookAddons = util.Parser.Malelooktype.Lookaddons
	}
	player.Town = town.Get()
	player.Stamina = util.Parser.Register.Stamina
	player.SkillAxe = util.Parser.Skills.Axe
	player.SkillSword = util.Parser.Skills.Sword
	player.SkillClub = util.Parser.Skills.Club
	player.SkillDist = util.Parser.Skills.Dist
	player.SkillFish = util.Parser.Skills.Fish
	player.SkillFist = util.Parser.Skills.Fist
	player.SkillShield = util.Parser.Skills.Shield
	player.Experience = util.Parser.Register.Experience
	err = player.Save()
	if err != nil {
		util.HandleError("Error creating player", err)
		http.Error(w, "Oops! Something wrong happened while creating player", http.StatusBadRequest)
		return
	}
	key := uniuri.New()
	for models.TokenExists(key) {
		key = uniuri.New()
	}
	err = account.UpdateToken(key)
	if err != nil {
		util.HandleError("Error updating your account token", err)
		http.Error(w, "Oops! Something wrong happened while setting your login token", http.StatusBadRequest)
		return
	}
	base.Session.Set("key", key)
	base.Session.Set("logged", 1)	
    err = base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while getting town list", http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, "/account/manage", 301)
}