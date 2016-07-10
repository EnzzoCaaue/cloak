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

// List shows the list of server houses
func (base *HouseController) List(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

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
	base.Data["Owner"] = false
	base.Data["OwnerName"] = "hola"
	base.Template = "house_view.html"
}
