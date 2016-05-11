package controllers

import (
    "net/http"
    "crypto/sha1"
	"github.com/julienschmidt/httprouter"
    "github.com/Cloakaac/cloak/template"
	"github.com/dchest/uniuri"
	"github.com/Cloakaac/cloak/util"
    "github.com/Cloakaac/cloak/models"
	"fmt"
)

type login struct {
    Errors []string
    Token string
}

// LoginForm saves the login form
type LoginForm struct {
    AccountName string `validate:"regexp=^[A-Za-z0-9]+$, min=5, max=20" alias:"Account name"`
    Password string `validate:"min=8, max=30" alias:"Account password"`
    Captcha string `validate:"validCaptcha" alias:"Captcha check"`
}

// Login shows the login form
func (base *BaseController) Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    csrfToken := uniuri.New()
	//base.Session.Set("token", csrfToken)
    response := &login{
        base.Session.GetFlashes("errors"),
        csrfToken,
        
    }
    err := base.Session.Save(req, w)
	if err != nil {
		util.HandleError("Error saving the current session", err)
		http.Error(w, "Oops! Something wrong happened while saving your request session", http.StatusBadRequest)
		return

	}
    template.Renderer.ExecuteTemplate(w, "login.html", response)
}

// SignIn process the login form
func (base *BaseController) SignIn(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    form := &LoginForm{
        req.FormValue("loginname"),
        req.FormValue("loginpassword"),
        req.FormValue("g-recaptcha-response"),
    }
    if errs := util.Validate(form); len(errs) > 0 {
        for _, v := range errs {
		    base.Session.AddFlash(v.Error(), "errors")
            base.Session.Save(req, w)
        }
        http.Redirect(w, req, "/account/login", 301)
        return
    }
    account := models.NewAccount()
    account.Account.Name = req.FormValue("loginname")
    hash := sha1.Sum([]byte(req.FormValue("loginpassword")))
    account.Account.Password = fmt.Sprintf("%x", hash)
    if !account.SignIn() {
        base.Session.AddFlash("Wrong account or password", "errors")
        base.Session.Save(req, w)
        http.Redirect(w, req, "/account/login", 301)
        return
    }
    key := uniuri.New()
    for models.TokenExists(key) {
        key = uniuri.New()
    }
    err := account.UpdateToken(key)
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