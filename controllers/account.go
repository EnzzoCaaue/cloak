package controllers

import (
	"fmt"
	"net/http"

	"crypto/sha1"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/dgryski/dgoogauth"
	"github.com/julienschmidt/httprouter"
	"net/url"
	"time"
)

type AccountController struct {
	*BaseController
}

type deletionForm struct {
	Password string
	Captcha  string `validate:"validCaptcha" alias:"Captcha check"`
}

type creationForm struct {
	Name     string `validate:"regexp=^[A-Z a-z]+$, max=14" alias:"Character name"`
	Town     string
	Gender   string `validate:"validGender" alias:"Character Gender"`
	Vocation string `validate:"validVocation" alias:"Character Vocation"`
	Captcha  string `validate:"validCaptcha" alias:"Captcha check"`
}

// AccountManage shows the account manage page
func (base *BaseController) AccountManage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	characters, err := base.Account.GetCharacters()
	if err != nil {
		base.Error = "Error while getting your character list"
		return
	}
	base.Data["Success"] = base.Session.GetFlashes("success")
	base.Data["Characters"] = characters
	base.Data["Account"] = base.Account
	base.Template = "manage.html"
}

// AccountLogout logs the user out
func (base *BaseController) AccountLogout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Session.Delete("key")
	base.Redirect = "/account/login"
}

// AccountSetRecovery sets an account recovery key
func (base *BaseController) AccountSetRecovery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if base.Account.RecoveryKey != "" {
		base.Redirect = "/account/manage"
		return
	}
	key := uniuri.New()
	err := base.Account.UpdateRecoveryKey(key)
	if err != nil {
		base.Error = "Error while updating your recovery key"
		return
	}
	base.Session.AddFlash("Your recovery key is <b>"+key+"</b>. Please write it down!", "success")
	base.Redirect = "/account/manage"
}

//AccountTwoFactor the form to set-up a two factor auth
func (base *BaseController) AccountTwoFactor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if base.Account.TwoFactor > 0 {
		base.Redirect = "/account/manage"
	}
	secretKey := uniuri.NewLenChars(16, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"))
	codeURL := fmt.Sprintf("otpauth://totp/%v:%v?secret=%v&issuer=%v", "MyServer", base.Account.Account.Name, secretKey, "MyServer")
	base.Data["QR"] = "http://chart.apis.google.com/chart?chs=500x500&cht=qr&choe=UTF-8&chl=" + codeURL
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Session.Set("secret", secretKey)
	base.Template = "account_twofactor.html"
}

// AccountSetTwoFactor Checks and sets a two factor key
func (base *BaseController) AccountSetTwoFactor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if base.Account.TwoFactor > 0 {
		base.Redirect = "/account/manage"
		return
	}
	otpConfig := &dgoogauth.OTPConfig{
		Secret:      base.Session.GetString("secret"),
		WindowSize:  3,
		HotpCounter: 0,
	}
	success, _ := otpConfig.Authenticate(req.FormValue("password"))
	if !success {
		base.Session.AddFlash("Wrong authenticator code", "errors")
		base.Redirect = "/account/manage/twofactor"
		return
	}
	err := base.Account.EnableTwoFactor(base.Session.GetString("secret"))
	if err != nil {
		base.Error = "Error while activating two factor on your account"
		return
	}
	// TODO: REMOVE SECRET FROM SESSION
	base.Session.AddFlash("Two Factor authenticator activated. Enjoy your new security level", "success")
	base.Redirect = "/account/manage"
}

// AccountDeleteCharacter shows the form to delete a character
func (base *BaseController) AccountDeleteCharacter(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	characterName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		base.Error = "Cannot escape character name"
		return
	}
	player := base.Account.GetCharacter(characterName)
	if player == nil {
		base.Redirect = "/account/manage"
		return
	}
	if player.Cloaka.Deleted == 1 {
		base.Redirect = "/account/manage"
	}
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Name"] = player.Name
	base.Template = "delete_character.html"
}

// DeleteCharacter deletes an account character
func (base *BaseController) DeleteCharacter(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	characterName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		http.Error(w, "Oops! Something while reading character name!", http.StatusBadRequest)
		return
	}
	player := base.Account.GetCharacter(characterName)
	if player == nil {
		base.Redirect = "/account/manage"
		return
	}
	if player.Cloaka.Deleted == 1 {
		base.Redirect = "/account/manage"
		return
	}
	form := &deletionForm{
		req.FormValue("password"),
		req.FormValue("g-recaptcha-response"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "errors")
		}
		base.Redirect = "/account/manage/delete/" + ps.ByName("name")
		return
	}
	password := fmt.Sprintf("%x", sha1.Sum([]byte(form.Password)))
	if base.Account.Account.Password != password {
		base.Session.AddFlash("Wrong password", "errors")
		base.Redirect = "/account/manage/delete/" + ps.ByName("name")
		return
	}
	deletion := time.Now().Unix() + ((3600 * 24) * 7)
	err = player.Delete(deletion)
	if err != nil {
		base.Error = "Error setting character deletion time"
		return
	}
	deletionDate := time.Unix(deletion, 0)
	base.Session.AddFlash(fmt.Sprintf("Character set for deletion at: <b>%v-%v-%v</b>", deletionDate.Month().String()[:3], deletionDate.Day(), deletionDate.Year()), "success")
	base.Redirect = "/account/manage"
}

// AccountCreateCharacter shows the form to create an account character
func (base *BaseController) AccountCreateCharacter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	towns, err := models.GetTowns()
	if err != nil {
		base.Error = "Error while getting town list"
		return
	}
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Towns"] = towns
	base.Template = "create_character.html"
}

// CreateCharacter creates an account character
func (base *BaseController) CreateCharacter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	form := &creationForm{
		req.FormValue("name"),
		req.FormValue("town"),
		req.FormValue("sex"),
		req.FormValue("vocation"),
		req.FormValue("g-recaptcha-response"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "errors")
		}
		base.Redirect = "/account/manage/create"
		return
	}
	town := models.GetTownByName(form.Town)
	if town == nil {
		base.Session.AddFlash("Unknown town name", "errors")
		base.Redirect = "/account/manage/create"
		return
	}
	player := models.NewPlayer()
	player.Name = form.Name
	player.AccountID = base.Account.Account.ID
	if player.Exists() {
		base.Session.AddFlash("Character name is already in use", "errors")
		base.Redirect = "/account/manage/create"
		return
	}
	player.Level = util.Parser.Register.Level
	player.Health = util.Parser.Register.Health
	player.HealthMax = util.Parser.Register.Healthmax
	player.Mana = util.Parser.Register.Mana
	player.ManaMax = util.Parser.Register.Manamax
	player.Vocation = util.Vocation(form.Vocation)
	player.Gender = util.Gender(form.Gender)
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
	player.Town = town
	player.Stamina = util.Parser.Register.Stamina
	player.SkillAxe = util.Parser.Register.Skills.Axe
	player.SkillSword = util.Parser.Register.Skills.Sword
	player.SkillClub = util.Parser.Register.Skills.Club
	player.SkillDist = util.Parser.Register.Skills.Dist
	player.SkillFish = util.Parser.Register.Skills.Fish
	player.SkillFist = util.Parser.Register.Skills.Fist
	player.SkillShield = util.Parser.Register.Skills.Shield
	player.Experience = util.Parser.Register.Experience
	err := player.Save()
	if err != nil {
		base.Error = "Error while creating your character"
		return
	}
	base.Session.AddFlash("Character <b>"+player.Name+"</b> created successfully", "success")
	base.Redirect = "/account/manage"
}
