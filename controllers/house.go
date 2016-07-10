package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
)

type HouseController struct {
	*pigo.Controller
}

// List shows the list of server houses
func (base *HouseController) List(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

}
