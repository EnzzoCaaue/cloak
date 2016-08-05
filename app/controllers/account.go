package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Cloakaac/cloak/app/models"
	"github.com/Cloakaac/cloak/app/util"
	"github.com/dchest/uniuri"
	"github.com/dgryski/dgoogauth"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"github.com/yaimko/yaimko"
)

type AccountController struct {
	*yaimko.Controller
}

// AccountLost shows the recover account form
func (base AccountController) AccountLost() *yaimko.Result {
	base.Data["ErrorsPassword"] = session.GetFlash("errorsPassword")
	base.Data["SuccessPassword"] = session.GetFlash("successPassword")
	base.Data["ErrorsName"] = session.GetFlash("errorsName")
	base.Data["SuccessName"] = session.GetFlash("successName")
	return base.Render("account_lost.html")
}

// AccountLostName recovers an account name using the recovery key and the password
func (base AccountController) AccountLostName(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	passwordSha1 := fmt.Sprintf("%x", sha1.Sum([]byte(req.FormValue("password"))))
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
	base.Data("Success", base.Session.GetFlashes("success"))
	base.Data("Characters", characters)
	base.Data("Account", base.Hook["account"].(*models.CloakaAccount))
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
	base.Data("QR", "http://chart.apis.google.com/chart?chs=500x500&cht=qr&choe=UTF-8&chl="+codeURL)
	base.Data("Errors", base.Session.GetFlashes("errors"))
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
	base.Data("Errors", base.Session.GetFlashes("errors"))
	base.Data("Name", player.Name)
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
	base.Data("Errors", base.Session.GetFlashes("errors"))
	base.Data("Towns", util.Towns.GetList())
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
	if !util.Towns.Exists(form.Town) {
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
	player.Level = viper.GetInt("register.level")
	player.Health = viper.GetInt("register.health")
	player.HealthMax = viper.GetInt("register.healthmax")
	player.Mana = viper.GetInt("register.mana")
	player.ManaMax = viper.GetInt("register.manamax")
	player.Vocation = util.Vocation(form.Vocation)
	player.Gender = util.Gender(form.Gender)
	if player.Gender == 0 { // female
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
	player.Town = util.Towns.Get(form.Town)
	player.Stamina = viper.GetInt("register.stamina")
	player.SkillAxe = viper.GetInt("register.skills.axe")
	player.SkillSword = viper.GetInt("register.skills.sword")
	player.SkillClub = viper.GetInt("register.skills.club")
	player.SkillDist = viper.GetInt("register.skills.dist")
	player.SkillFish = viper.GetInt("register.skills.fish")
	player.SkillFist = viper.GetInt("register.skills.fist")
	player.SkillShield = viper.GetInt("register.skills.shield")
	player.Experience = viper.GetInt("register.experience")
	err := player.Save()
	if err != nil {
		base.Error = "Error while creating your character"
		return
	}
	base.Session.AddFlash("Character <b>"+player.Name+"</b> created successfully", "success")
	base.Redirect = "/account/manage"
}
