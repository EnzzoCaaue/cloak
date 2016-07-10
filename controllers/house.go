package controllers

import (
	"net/http"

	"net/url"

	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
)

type HouseController struct {
	*pigo.Controller
}

// List shows the list of server houses within a random town
func (base *HouseController) List(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	town := util.Towns.GetRandom()
	houses := util.Houses.GetList(town.ID)
	base.Data["Houses"] = houses
	base.Data["Town"] = town
	base.Data["Towns"] = util.Towns.GetList()
	base.Template = "houses.html"
}

// ListName shows the list of server houses by its town
func (base *HouseController) ListName(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	town := util.Towns.Get(req.FormValue("town"))
	if town == nil {
		base.Redirect = "/houses/list"
		return
	}
	houses := util.Houses.GetList(town.ID)
	base.Data["Houses"] = houses
	base.Data["Town"] = town
	base.Data["Towns"] = util.Towns.GetList()
	base.Template = "houses.html"
}

// View shows a house page
func (base *HouseController) View(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	houseName, err := url.QueryUnescape(ps.ByName("name"))
	if err != nil {
		base.Error = err.Error()
		return
	}
	house := util.Houses.GetHouseByName(houseName)
	if house == nil {
		base.Redirect = "/"
		return
	}
	base.Data["Info"] = house
	base.Data["Town"] = util.Towns.GetByID(house.TownID)
	base.Template = "house_view.html"
}
