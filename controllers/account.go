package controllers

import (
	"fmt"
	"net/http"

	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"crypto/sha1"
	"time"
	"github.com/dchest/uniuri"
	"github.com/dgryski/dgoogauth"
	"github.com/julienschmidt/httprouter"
	"net/url"
)

type manage struct {
	Success    []string
	Characters []*models.Player
	Account    *models.CloakaAccount
	Token      string
}

type twof struct {
	QR string
    Errors []string
}

type deletion struct {
	Token string
	Errors []string
	Name string
}

type deletionForm struct {
	Password string
	Captcha string `validate:"validCaptcha" alias:"Captcha check"`
}

type creation struct {
	Token string
	Errors []string
	Towns []*models.Town
}

type creationForm struct {
	Name string `validate:"regexp=^[A-Z a-z]+$, max=14" alias:"Character name"`
	Town string
	Gender string `validate:"validGender" alias:"Character Gender"`
	Vocation string `validate:"validVocation" alias:"Character Vocation"`
	Captcha string `validate:"validCaptcha" alias:"Captcha check"`
}

// AccountManage shows the account manage page
func (base *BaseController) AccountManage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	characters, err := account.GetCharacters()
	if err != nil {
		http.Error(w, "Oops! Something wrong happened while getting your character list", http.StatusBadRequest)
		return
	}
	response := &manage{
		base.Session.GetFlashes("success"),
		characters,
		account,
		uniuri.New(),
	}
	err = base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving current session", http.StatusBadRequest)
		return
	}
	template.Renderer.ExecuteTemplate(w, "manage.html", response)
}

// AccountLogout logs the user out
func (base *BaseController) AccountLogout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	base.Session.Delete("key")
	err := base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving current session", http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, "/account/login", http.StatusMovedPermanently)
}

// AccountSetRecovery sets an account recovery key
func (base *BaseController) AccountSetRecovery(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	if account.RecoveryKey != "" {
		http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
		return
	}
	key := uniuri.New()
	err := account.UpdateRecoveryKey(key)
	if err != nil {
		util.HandleError("Error while updating account recovery key", err)
		http.Error(w, "Oops! Something wrong happened while updating your recovery key!", http.StatusBadRequest)
		return
	}
	base.Session.AddFlash("Your recovery key is <b>"+key+"</b>. Please write it down!", "success")
	err = base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving current session", http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
}

//AccountTwoFactor the form to set-up a two factor auth
func (base *BaseController) AccountTwoFactor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	if account.TwoFactor > 0 {
		http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
		return
	}
	secretKey := uniuri.NewLenChars(16, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"))
	codeURL := fmt.Sprintf("otpauth://totp/%v:%v?secret=%v&issuer=%v", "MyServer", account.Account.Name, secretKey, "MyServer")
	response := &twof{
		"http://chart.apis.google.com/chart?chs=500x500&cht=qr&choe=UTF-8&chl=" + codeURL,
        base.Session.GetFlashes("errors"),
	}
	base.Session.Set("secret", secretKey)
	err := base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving current session", http.StatusBadRequest)
		return
	}
	template.Renderer.ExecuteTemplate(w, "account_twofactor.html", response)
}

// AccountSetTwoFactor Checks and sets a two factor key
func (base *BaseController) AccountSetTwoFactor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	if account.TwoFactor > 0 {
		http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
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
        base.Session.Save(req, w)
        http.Redirect(w, req, "/account/manage/twofactor", http.StatusMovedPermanently)
        return
    }
    err := account.EnableTwoFactor(base.Session.GetString("secret"))
    if err != nil {
        util.HandleError("Error while updating accounts secret row", err)
        http.Error(w, "Cannot activate two-factor on your account", 500)
        return
    }
    // TODO: REMOVE SECRET FROM SESSION
    base.Session.AddFlash("Two Factor authenticator activated. Enjoy your new security level", "success")
    base.Session.Save(req, w)
    http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
}

// AccountDeleteCharacter shows the form to delete a character
func (base *BaseController) AccountDeleteCharacter(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	characterName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		http.Error(w, "Oops! Something while reading character name!", http.StatusBadRequest)
		return
	}
	player := account.GetCharacter(characterName)
	if player.Cloaka.Deleted == 1 {
		http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
	}
	token := uniuri.New()
	response := &deletion{
		token,
		base.Session.GetFlashes("errors"),
		player.Name,
	}
	base.Session.Save(req, w)
	template.Renderer.ExecuteTemplate(w, "delete_character.html", response)
}

// DeleteCharacter deletes an account character
func (base *BaseController) DeleteCharacter(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	characterName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		http.Error(w, "Oops! Something while reading character name!", http.StatusBadRequest)
		return
	}
	player := account.GetCharacter(characterName)
	if player.Cloaka.Deleted == 1 {
		http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
	}
	form := &deletionForm{
		req.FormValue("password"),
		req.FormValue("g-recaptcha-response"),	
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for i := range errs {
			base.Session.AddFlash(errs[i].Error(), "errors")
		}
		base.Session.Save(req, w)
		http.Redirect(w, req, "/account/manage/delete/" + ps.ByName("name"), http.StatusMovedPermanently)
		return
	}
	password := fmt.Sprintf("%x", sha1.Sum([]byte(form.Password)))
	if account.Account.Password != password {
		base.Session.AddFlash("Wrong password", "errors")
		base.Session.Save(req, w)
		http.Redirect(w, req, "/account/manage/delete/" + ps.ByName("name"), http.StatusMovedPermanently)
		return
	}
	deletion := time.Now().Unix() + ((3600 * 24) * 7)
	err = player.Delete(deletion)
	if err != nil {
		util.HandleError("Error setting a deletion date", err)
		http.Error(w, "Error setting your character deletion date", 500)
		return
	}
	deletionDate := time.Unix(deletion, 0)
	base.Session.AddFlash(fmt.Sprintf("Character set for deletion at: <b>%v-%v-%v</b>", deletionDate.Month().String()[:3], deletionDate.Day(), deletionDate.Year()), "success")
	base.Session.Save(req, w)
	http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
}

// AccountCreateCharacter shows the form to create an account character
func (base *BaseController) AccountCreateCharacter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
	token := uniuri.New()
	towns, err := models.GetTowns()
	if err != nil {
		util.HandleError("Error getting town list", err)
		http.Error(w, "Cannot get town list", 500)
		return
	}
	response := &creation{
		token,
		base.Session.GetFlashes("errors"),
		towns,
	}
	base.Session.Save(req, w)
	template.Renderer.ExecuteTemplate(w, "create_character.html", response)
}

// CreateCharacter creates an account character
func (base *BaseController) CreateCharacter(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	account := models.GetAccountByToken(base.Session.GetString("key"))
	if account == nil {
		http.Error(w, "Oops! Something wrong happened while getting your account!", http.StatusBadRequest)
		return
	}
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
		base.Session.Save(req, w)
		http.Redirect(w, req, "/account/manage/create", http.StatusMovedPermanently)
		return
	}
	town := models.GetTownByName(form.Town)
	if town == nil {
		base.Session.AddFlash("Unknown town name", "errors")
		base.Session.Save(req, w)
		http.Redirect(w, req, "/account/manage/create", http.StatusMovedPermanently)
		return
	}
	player := models.NewPlayer()
	player.Name = form.Name
	player.AccountID = account.Account.ID
	if player.Exists() {
		base.Session.AddFlash("Character name is already in use", "errors")
		base.Session.Save(req, w)
		http.Redirect(w, req, "/account/manage/create", http.StatusMovedPermanently)
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
		util.HandleError("Error creating player", err)
		http.Error(w, "Oops! Something wrong happened while creating player", http.StatusBadRequest)
		return
	}
	base.Session.AddFlash("Character <b>" + player.Name + "</b> created successfully", "success")
	base.Session.Save(req, w)
	http.Redirect(w, req, "/account/manage", http.StatusMovedPermanently)
}