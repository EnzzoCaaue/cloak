package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/dgryski/dgoogauth"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"log"
	"net/http"
	"net/url"
	"time"
)

type AccountController struct {
	*pigo.Controller
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

type passwordLostForm struct {
	Password string `validate:"min=8, max=30" alias:"Account password"`
}

// AccountLost shows the recover account form
func (base *AccountController) AccountLost(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Data["ErrorsPassword"] = base.Session.GetFlashes("errorsPassword")
	base.Data["SuccessPassword"] = base.Session.GetFlashes("successPassword")
	base.Data["ErrorsName"] = base.Session.GetFlashes("errorsName")
	base.Data["SuccessName"] = base.Session.GetFlashes("successName")
	base.Template = "account_lost.html"
}

// AccountLostName recovers an account name using the recovery key and the password
func (base *AccountController) AccountLostName(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	passwordSha1 := fmt.Sprintf("%x", sha1.Sum([]byte(req.FormValue("password"))))
	log.Println(req.FormValue("key"))
	name := models.RecoverAccountName(req.FormValue("key"), passwordSha1)
	if name == "" {
		base.Session.AddFlash("Wrong account password or recovery key", "errorsName")
		base.Redirect = "/account/lost"
		return
	}
	base.Session.AddFlash("Your account name is <b>"+name+"</b>", "successName")
	base.Redirect = "/account/lost"
}

// AccountLostPassword recovers an account using the recovery key and the name
func (base *AccountController) AccountLostPassword(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if !models.RecoverAccountPassword(req.FormValue("key"), req.FormValue("name")) {
		base.Session.AddFlash("Wrong account name or recovery key", "errorsPassword")
		base.Redirect = "/account/lost"
		return
	}
	form := &passwordLostForm{
		req.FormValue("password"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "errorsPassword")
		}
		base.Redirect = "/account/lost"
		return
	}
	passwordSha1 := sha1.Sum([]byte(req.FormValue("password")))
	err := models.SetNewPassword(req.FormValue("name"), fmt.Sprintf("%x", passwordSha1))
	if err != nil {
		base.Error = "Error while updating your acocunt password"
		return
	}
	base.Session.AddFlash("Password changed successfully", "successPassword")
	base.Redirect = "/account/lost"
}

// AccountManage shows the account manage page
func (base *AccountController) AccountManage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	characters, err := base.Hook["account"].(*models.CloakaAccount).GetCharacters()
	if err != nil {
		base.Error = "Error while getting your character list"
		return
	}
	base.Data["Success"] = base.Session.GetFlashes("success")
	base.Data["Characters"] = characters
	base.Data["Account"] = base.Hook["account"].(*models.CloakaAccount)
	base.Template = "manage.html"
}

// AccountLogout logs the user out
func (base *AccountController) AccountLogout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Session.Delete("key")
	base.Redirect = "/account/login"
}

// AccountSetRecovery sets an account recovery key
func (base *AccountController) AccountSetRecovery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if base.Hook["account"].(*models.CloakaAccount).RecoveryKey != "" {
		base.Redirect = "/account/manage"
		return
	}
	key := uniuri.New()
	err := base.Hook["account"].(*models.CloakaAccount).UpdateRecoveryKey(key)
	if err != nil {
		base.Error = "Error while updating your recovery key"
		return
	}
	base.Session.AddFlash("Your recovery key is <b>"+key+"</b>. Please write it down!", "success")
	base.Redirect = "/account/manage"
}

//AccountTwoFactor the form to set-up a two factor auth
func (base *AccountController) AccountTwoFactor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if base.Hook["account"].(*models.CloakaAccount).TwoFactor > 0 {
		base.Redirect = "/account/manage"
	}
	secretKey := uniuri.NewLenChars(16, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"))
	codeURL := fmt.Sprintf("otpauth://totp/%v:%v?secret=%v&issuer=%v", "MyServer", base.Hook["account"].(*models.CloakaAccount).Account.Name, secretKey, "MyServer")
	base.Data["QR"] = "http://chart.apis.google.com/chart?chs=500x500&cht=qr&choe=UTF-8&chl=" + codeURL
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Session.Set("secret", secretKey)
	base.Template = "account_twofactor.html"
}

// AccountSetTwoFactor Checks and sets a two factor key
func (base *AccountController) AccountSetTwoFactor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if base.Hook["account"].(*models.CloakaAccount).TwoFactor > 0 {
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
	err := base.Hook["account"].(*models.CloakaAccount).EnableTwoFactor(base.Session.GetString("secret"))
	if err != nil {
		base.Error = "Error while activating two factor on your account"
		return
	}
	// TODO: REMOVE SECRET FROM SESSION
	base.Session.AddFlash("Two Factor authenticator activated. Enjoy your new security level", "success")
	base.Redirect = "/account/manage"
}

// AccountDeleteCharacter shows the form to delete a character
func (base *AccountController) AccountDeleteCharacter(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	characterName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		base.Error = "Cannot escape character name"
		return
	}
	player := base.Hook["account"].(*models.CloakaAccount).GetCharacter(characterName)
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
func (base *AccountController) DeleteCharacter(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	characterName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		http.Error(w, "Oops! Something while reading character name!", http.StatusBadRequest)
		return
	}
	player := base.Hook["account"].(*models.CloakaAccount).GetCharacter(characterName)
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
	if base.Hook["account"].(*models.CloakaAccount).Account.Password != password {
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
func (base *AccountController) AccountCreateCharacter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
func (base *AccountController) CreateCharacter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
	player.AccountID = base.Hook["account"].(*models.CloakaAccount).Account.ID
	if player.Exists() {
		base.Session.AddFlash("Character name is already in use", "errors")
		base.Redirect = "/account/manage/create"
		return
	}
	player.Level = pigo.Config.Key("register").Int("level")
	player.Health = pigo.Config.Key("register").Int("health")
	player.HealthMax = pigo.Config.Key("register").Int("healthmax")
	player.Mana = pigo.Config.Key("register").Int("mana")
	player.ManaMax = pigo.Config.Key("register").Int("manamax")
	player.Vocation = util.Vocation(form.Vocation)
	player.Gender = util.Gender(form.Gender)
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
	player.Town = town
	player.Stamina = pigo.Config.Key("register").Int("stamina")
	player.SkillAxe = pigo.Config.Key("register").Key("skills").Int("axe")
	player.SkillSword = pigo.Config.Key("register").Key("skills").Int("sword")
	player.SkillClub = pigo.Config.Key("register").Key("skills").Int("club")
	player.SkillDist = pigo.Config.Key("register").Key("skills").Int("dist")
	player.SkillFish = pigo.Config.Key("register").Key("skills").Int("fish")
	player.SkillFist = pigo.Config.Key("register").Key("skills").Int("fist")
	player.SkillShield = pigo.Config.Key("register").Key("skills").Int("shield")
	player.Experience = pigo.Config.Key("register").Int("experience")
	err := player.Save()
	if err != nil {
		base.Error = "Error while creating your character"
		return
	}
	base.Session.AddFlash("Character <b>"+player.Name+"</b> created successfully", "success")
	base.Redirect = "/account/manage"
}
