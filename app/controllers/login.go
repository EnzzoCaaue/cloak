package controllers

import "github.com/yaimko/yaimko"

type LoginController struct {
	*yaimko.Controller
}

// Login shows the login form
func (base LoginController) Login() *yaimko.Result {
	return base.Render("login.html")
}
