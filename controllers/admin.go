package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
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
	onlineRecords, err := models.GetOnlineRecords(10)
	if err != nil {
		base.Error = "Error while trying to get online records"
		return
	}
	base.Data["Records"] = onlineRecords
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

// Server shows the TFS server manager
func (base *AdminController) Server(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	base.Template = "admin_server.html"
}
