package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
    "github.com/Cloakaac/cloak/models"
    "net/http"
    "runtime"
)

type AdminController struct {
	*pigo.Controller
}

// Dashboard shows the admin dashboard
func (base *AdminController) Dashboard(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
    adminInfo := models.GetAdminInformation()
    m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
    base.Data["Memstats"] = m
    base.Data["Goversion"] = runtime.Version()
    base.Data["Numcpu"] = runtime.NumCPU()
    base.Data["Numroutine"] = runtime.NumGoroutine()
    base.Data["Numcgo"] = runtime.NumCgoCall()
    base.Data["AccountTotal"] = adminInfo.Accounts
    base.Data["PlayerTotal"] = adminInfo.Players
    base.Data["MaleTotal"] = adminInfo.Males
    base.Data["FemaleTotal"] = adminInfo.Females
    base.Data["SorcererTotal"] = adminInfo.Sorcerers
    base.Data["DruidTotal"] = adminInfo.Druids
    base.Data["PaladinTotal"] = adminInfo.Paladins
    base.Data["KnightTotal"] = adminInfo.Knights
    base.Template = "admin.html"
}