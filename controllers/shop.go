package controllers

import (
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"net/http"
	"strconv"
	"time"
	"fmt"
	"log"
	"github.com/Cloakaac/cloak/models"
)

const (
	sandbox = "sandbox"
	live    = "live"
)

var (
	baseURL     string
	paypalToken *util.PaypalToken
)

type ShopController struct {
	*pigo.Controller
}

func init() {
	if pigo.Config.Key("paypal").String("mode") == sandbox {
		baseURL = "https://api.sandbox.paypal.com"
	} else {
		baseURL = "https://api.paypal.com"
	}
}

// Paypal shows the paypal buypoints page
func (base *ShopController) Paypal(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Success"] = base.Session.GetFlashes("success")
	base.Data["Points"] = pigo.Config.Key("paypal").Key("payment").Float("points")
	base.Data["Min"] = pigo.Config.Key("paypal").Key("payment").Float("min")
	base.Data["Max"] = pigo.Config.Key("paypal").Key("payment").Float("max")
	base.Data["Promo"] = pigo.Config.Key("paypal").Float("promo")
	base.Template = "paypal.html"
}

// PaypalPay process a paypal buypoints request
func (base *ShopController) PaypalPay(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	payAmount, err := strconv.ParseFloat(req.FormValue("pay"), 64)
	if err != nil {
		base.Session.AddFlash("Payment amount needs to be a number", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payAmount > pigo.Config.Key("paypal").Key("payment").Float("max") {
		base.Session.AddFlash("Payment amount is too high", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payAmount < pigo.Config.Key("paypal").Key("payment").Float("min") {
		base.Session.AddFlash("Payment amount is too low", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	timeNow := time.Now().Unix()
	if paypalToken == nil || (timeNow+paypalToken.ExpiresIn) < timeNow {
		token, err := util.GetPaypalToken(baseURL, pigo.Config.Key("paypal").String("public"), pigo.Config.Key("paypal").String("private"))
		if err != nil {
			base.Error = err.Error()
			return
		}
		paypalToken = token
	}
	hostURL := ""
	if pigo.Config.Key("https").Bool("enabled") {
		hostURL = fmt.Sprintf("%v://%v", "https", req.Host)
	} else {
		hostURL = fmt.Sprintf("%v://%v", "https", req.Host)
	}
	log.Println(hostURL)
	payment, err := util.CreatePaypalPayment(
		hostURL,
		baseURL, 
		paypalToken.Token, 
		req.FormValue("pay"),
		pigo.Config.Key("paypal").String("description"),
		pigo.Config.Key("paypal").String("currency"),
	)
	if err != nil {
		base.Session.AddFlash("Something went wrong while creating your payment", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payment.State != "created" {
		base.Session.AddFlash("Your payment cannot be created. Please try again later", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	for i := range payment.Links {
		if payment.Links[i].Rel == "approval_url" {
			base.Redirect = payment.Links[i].Href
			return
		}
	}
	base.Session.AddFlash("Error while trying to get your payment approval URL. Please try again later", "errors")
	base.Redirect = "/buypoints/paypal"
}

// PaypalProcess process a paypal payment
func (base *ShopController) PaypalProcess(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	paymentID := req.URL.Query().Get("paymentId")
	payerID := req.URL.Query().Get("PayerID")
	if paymentID == "" {
		base.Session.AddFlash("Missing payment ID. Cannot process your payment", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payerID == "" {
		base.Session.AddFlash("Missing payer ID. Cannot process your payment", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	payment, err := util.ProcessPaypalPayment(baseURL, payerID, paymentID, paypalToken.Token)
	if err != nil {
		base.Session.AddFlash("Your payment could notbe processed. No money its been taking from your account", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payment.IsEmpty() {
		base.Session.AddFlash("Invalid payment ID", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	paid, err := strconv.ParseFloat(payment.Transactions[0].Amount.Total, 64)
	if err != nil {
		base.Session.AddFlash("Your payment amount could not be processed. Please contant an administrator", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	totalCoins := 0
	if pigo.Config.Key("paypal").Float("promo") > 0 {
		totalCoinsNoPromo := paid * pigo.Config.Key("paypal").Key("payment").Float("points")
		totalCoins = int(((pigo.Config.Key("paypal").Float("promo") * totalCoinsNoPromo) / 100) + totalCoinsNoPromo);
	} else {
		totalCoins = int(paid * pigo.Config.Key("paypal").Key("payment").Float("points"))
	}
	err = base.Hook["account"].(*models.CloakaAccount).UpdatePoints(totalCoins)
	if err != nil {
		base.Error = "Unable to update account points"
		return
	}
	base.Session.AddFlash("Payment completed. We added "+strconv.Itoa(totalCoins)+" coins to your account. Enjoy!", "success")
	base.Redirect = "/buypoints/paypal"
}