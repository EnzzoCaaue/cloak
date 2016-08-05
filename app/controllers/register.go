package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"github.com/spf13/viper"
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
	base.Data("Errors", base.Session.GetFlashes("errors"))
	base.Data("Towns", util.Towns.GetList())
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
	if !util.Towns.Exists(form.CharacterTown) {
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
	account.Account.Premdays = viper.GetInt("Register.Premdays")
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
	player.Level = viper.GetInt("register.level")
	player.Health = viper.GetInt("register.health")
	player.HealthMax = viper.GetInt("register.healthmax")
	player.Mana = viper.GetInt("register.mana")
	player.ManaMax = viper.GetInt("register.manamax")
	player.Vocation = util.Vocation(form.CharacterVocation)
	player.Gender = util.Gender(form.CharacterSex)
	if player.Gender == 0 {
		player.LookBody = viper.GetInt("register.female.lookbody")
		player.LookFeet = viper.GetInt("register.female.lookfeet")
		player.LookHead = viper.GetInt("register.female.lookhead")
		player.LookType = viper.GetInt("register.female.looktype")
		player.LookAddons = viper.GetInt("register.female.lookaddons")
	} else {
		player.LookBody = viper.GetInt("register.male.lookbody")
		player.LookFeet = viper.GetInt("register.male.lookfeet")
		player.LookHead = viper.GetInt("register.male.lookhead")
		player.LookType = viper.GetInt("register.male.looktype")
		player.LookAddons = viper.GetInt("register.male.lookaddons")
	}
	player.Town = util.Towns.Get(form.CharacterTown)
	player.Stamina = viper.GetInt("register.stamina")
	player.SkillAxe = viper.GetInt("register.skills.axe")
	player.SkillSword = viper.GetInt("register.skills.sword")
	player.SkillClub = viper.GetInt("register.skills.club")
	player.SkillDist = viper.GetInt("register.skills.dist")
	player.SkillFish = viper.GetInt("register.skills.fish")
	player.SkillFist = viper.GetInt("register.skills.fist")
	player.SkillShield = viper.GetInt("register.skills.shield")
	player.Experience = viper.GetInt("register.experience")
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
