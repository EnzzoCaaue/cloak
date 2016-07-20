package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Cloakaac/cloak/models"
	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/pigo"
	"github.com/spf13/viper"
)

const (
	sandbox = "sandbox"
	live    = "live"
)

var (
	baseURL     string
	paypalToken = &Token{
		nil,
		&sync.RWMutex{},
	}
)

type paypalForm struct {
	Captcha string `validate:"validCaptcha" alias:"Captcha check"`
}

// Token holds a paypal token with a mutex
type Token struct {
	PaypalToken *util.PaypalToken
	rw          *sync.RWMutex
}

type ShopController struct {
	*pigo.Controller
}

func init() {
	if viper.GetString("paypal.mode") == sandbox {
		baseURL = "https://api.sandbox.paypal.com"
	} else {
		baseURL = "https://api.paypal.com"
	}
}

// Paypal shows the paypal buypoints page
func (base *ShopController) Paypal(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	base.Data["Success"] = base.Session.GetFlashes("success")
	base.Data["Points"] = viper.GetInt("paypal.payment.points")
	base.Data["Min"] = viper.GetInt("paypal.payment.min")
	base.Data["Max"] = viper.GetInt("paypal.payment.max")
	base.Data["Promo"] = viper.GetInt("paypal.promo")
	base.Template = "paypal.html"
}

// PaypalPay process a paypal buypoints request
func (base *ShopController) PaypalPay(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	form := &paypalForm{
		req.FormValue("g-recaptcha-response"),
	}
	if errs := util.Validate(form); len(errs) > 0 {
		for _, v := range errs {
			base.Session.AddFlash(v.Error(), "errors")
		}
		base.Redirect = "/buypoints/paypal"
		return
	}
	payAmount, err := strconv.Atoi(req.FormValue("pay"))
	if err != nil {
		base.Session.AddFlash("Payment amount needs to be a number", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payAmount > viper.GetInt("paypal.payment.max") {
		base.Session.AddFlash("Payment amount is too high", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	if payAmount < viper.GetInt("paypal.payment.min") {
		base.Session.AddFlash("Payment amount is too low", "errors")
		base.Redirect = "/buypoints/paypal"
		return
	}
	timeNow := time.Now().Unix()
	if paypalToken.PaypalToken == nil || (timeNow+paypalToken.PaypalToken.ExpiresIn) < timeNow {
		paypalToken.rw.Lock()
		token, err := util.GetPaypalToken(baseURL, viper.GetString("paypal.public"), viper.GetString("paypal.private"))
		if err != nil {
			base.Error = err.Error()
			return
		}
		paypalToken.PaypalToken = token
		paypalToken.rw.Unlock()
	}
	hostURL := ""
	if viper.GetBool("https.enabled") {
		hostURL = fmt.Sprintf("%v://%v", "https", req.Host)
	} else {
		hostURL = fmt.Sprintf("%v://%v", "https", req.Host)
	}
	paypalToken.rw.RLock()
	defer paypalToken.rw.RUnlock()
	payment, err := util.CreatePaypalPayment(
		hostURL,
		baseURL,
		paypalToken.PaypalToken.Token,
		req.FormValue("pay"),
		viper.GetString("paypal.description"),
		viper.GetString("paypal.currency"),
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
	paypalToken.rw.RLock()
	defer paypalToken.rw.RUnlock()
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
	payment, err := util.ProcessPaypalPayment(baseURL, payerID, paymentID, paypalToken.PaypalToken.Token)
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
	if viper.GetInt("paypal.promo") > 0 {
		totalCoinsNoPromo := paid * viper.GetFloat64("paypal.payment.points")
		totalCoins = int(((viper.GetFloat64("paypal.promo") * totalCoinsNoPromo) / 100) + totalCoinsNoPromo)
	} else {
		totalCoins = int(paid * viper.GetFloat64("paypal.payment.points"))
	}
	err = base.Hook["account"].(*models.CloakaAccount).UpdatePoints(totalCoins)
	if err != nil {
		base.Error = "Unable to update account points"
		return
	}
	base.Session.AddFlash("Payment completed. We added "+strconv.Itoa(totalCoins)+" coins to your account. Enjoy!", "success")
	base.Redirect = "/buypoints/paypal"
}

// ShopView loads and shows the "donation" shop
func (base *ShopController) ShopView(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	base.Data["Errors"] = base.Session.GetFlashes("errors")
	categories, err := models.GetCategories()
	if err != nil {
		base.Error = err.Error()
		return
	}
	randomCategory := rand.Intn(len(categories))
	categories[randomCategory].Active = true
	offers, err := categories[randomCategory].GetOffers()
	if err != nil {
		base.Error = err.Error()
		return
	}
	base.Data["Offers"] = offers
	base.Data["Categories"] = categories
	base.Template = "shop.html"
}
