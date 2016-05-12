package controllers

import (
	"fmt"
	"net/http"

	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/template"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/dgryski/dgoogauth"
	"github.com/julienschmidt/httprouter"
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
