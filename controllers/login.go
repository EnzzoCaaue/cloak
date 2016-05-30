package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/dchest/uniuri"
	"github.com/dgryski/dgoogauth"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/raggaer/pigo"
)

type LoginController struct {
	*pigo.Controller
}

// LoginForm saves the login form
type LoginForm struct {
	AccountName string `validate:"regexp=^[A-Za-z0-9]+$, min=5, max=20" alias:"Account name"`
	Password    string `validate:"min=8, max=30" alias:"Account password"`
	Captcha     string `validate:"validCaptcha" alias:"Captcha check"`
}

// Login shows the login form
func (base *LoginController) Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Template = "login.html"
}

// SignIn process the login form
func (base *LoginController) SignIn(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	form := &LoginForm{
		req.FormValue("loginname"),
		req.FormValue("loginpassword"),
		req.FormValue("g-recaptcha-response"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for _, v := range errs {
			base.Session.AddFlash(v.Error(), "errors")
		}
		base.Redirect = "/account/login"
		return
	}
	account := models.NewAccount()
	account.Account.Name = req.FormValue("loginname")
	hash := sha1.Sum([]byte(req.FormValue("loginpassword")))
	account.Account.Password = fmt.Sprintf("%x", hash)
	if !account.SignIn() {
		base.Session.AddFlash("Wrong account or password", "errors")
		base.Redirect = "/account/login"
		return
	}
	if account.TwoFactor > 0 && pigo.Config.String("mode") != "DEV" {
		otpConfig := &dgoogauth.OTPConfig{
			Secret:      account.Account.SecretKey,
			WindowSize:  3,
			HotpCounter: 0,
		}
		success, _ := otpConfig.Authenticate(req.FormValue("logincode"))
		if !success {
			base.Session.AddFlash("Wrong two-factor code", "errors")
			base.Redirect = "/account/login"
			return
		}
	}
	key := uniuri.New()
	for models.TokenExists(key) {
		key = uniuri.New()
	}
	err := account.UpdateToken(key)
	if err != nil {
		base.Error = "Unable to update your token key"
		return
	}
	base.Session.Set("key", key)
	base.Session.Set("logged", 1)
	base.Redirect = "/account/manage"
}
