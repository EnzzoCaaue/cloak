package controllers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// DisplayError displays the error page
func (base *BaseController) DisplayError(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http.Error(w, "Oops! Something bad happened", 500)
}
