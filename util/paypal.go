package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	cloakaPaypalProcessURL  = "/buypoints/paypal/process"
	cloakaPaypalURL         = "/buypoints/paypal"
	paypalTokenURL          = "/v1/oauth2/token"
	paypalCreatePaymentURL  = "/v1/payments/payment"
	paypalProcessPaymentURL = "/v1/payments/payment/%v/execute/"
)

// PaypalToken is the main oauth struct for paypal
type PaypalToken struct {
	Scope     string `json:"scope"`
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

// PaypalPayment struct for a created payment
type PaypalPayment struct {
	ID           string              `json:"id"`
	CreateTime   string              `json:"create_time"`
	UpdateTime   string              `json:"update_time"`
	State        string              `json:"state"`
	Intent       string              `json:"intent"`
	Payer        PaypalPayer         `json:"payer"`
	Transactions []PaypalTransaction `json:"transactions"`
	Links        []PaypalLink        `json:"links"`
}

// PaypalPaymentCreation struct for creating paypal payments
type PaypalPaymentCreation struct {
	Intent       string              `json:"intent"`
	RedirectURL  PaypalRedirectURL   `json:"redirect_urls"`
	Payer        PaypalPayer         `json:"payer"`
	Transactions []PaypalTransaction `json:"transactions"`
}

// PaypalPaymentInformation contains information about an executed payment
type PaypalPaymentInformation struct {
	ID           string
	Intent       string                 `json:"intent"`
	State        string                 `json:"state"`
	Payer        PaypalPayerInformation `json:"payer"`
	Transactions []PaypalTransaction    `json:"transactions"`
}

// PaypalRedirectURL redirect url payment
type PaypalRedirectURL struct {
	ReturnURL string `json:"return_url"`
	CancelURL string `json:"cancel_url"`
}

// PaypalPayer payment method
type PaypalPayer struct {
	PaymentMethod string `json:"payment_method"`
}

// PaypalPayerInformation holds information about the payer
type PaypalPayerInformation struct {
	PaymentMethod string          `json:"payment_method"`
	Status        string          `json:"status"`
	PayerInfo     PaypalPayerInfo `json:"payer_info"`
}

// PaypalPayerInfo holds information about the payer
type PaypalPayerInfo struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PayerID     string `json:"payer_id"`
	CountryCode string `json:"country_code"`
}

// PaypalTransaction array of paypal payment transactions
type PaypalTransaction struct {
	Amount      PaypalAmount `json:"amount"`
	Description string       `json:"description"`
}

// PaypalAmount paypal payment amount
type PaypalAmount struct {
	Total    string `json:"total"`
	Currency string `json:"currency"`
	//Details PaypalAmountDetails `json:"details"`
}

// PaypalAmountDetails detials of paypal payment amount
type PaypalAmountDetails struct {
	Subtotal float64 `json:"subtotal"`
}

// PaypalLink links of created payment
type PaypalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// PaypalPaymentProcess used to execute payments
type PaypalPaymentProcess struct {
	PayerID string `json:"payer_id"`
}

// IsEmpty checks if a payment is processed
func (p *PaypalPaymentInformation) IsEmpty() bool {
	if p.ID == "" {
		return true
	}
	return false
}

// ProcessPaypalPayment executes a paypal payment
func ProcessPaypalPayment(baseURL, payerID, paymentID, paypalToken string) (*PaypalPaymentInformation, error) {
	payment := &PaypalPaymentProcess{
		PayerID: payerID,
	}
	resp, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(resp)
	req, err := http.NewRequest(http.MethodPost, baseURL+fmt.Sprintf(paypalProcessPaymentURL, paymentID), buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+paypalToken)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	paymentResponse := &PaypalPaymentInformation{}
	err = json.Unmarshal(body, paymentResponse)
	if err != nil {
		return nil, err
	}
	return paymentResponse, nil
}

// CreatePaypalPayment creates a paypal payment and returns the response
func CreatePaypalPayment(hostURL, baseURL, paypalToken, amount, description, currency string) (*PaypalPayment, error) {
	payment := &PaypalPaymentCreation{
		Intent: "sale",
		RedirectURL: PaypalRedirectURL{
			ReturnURL: hostURL + cloakaPaypalProcessURL,
			CancelURL: hostURL + cloakaPaypalURL,
		},
		Payer: PaypalPayer{
			PaymentMethod: "paypal",
		},
		Transactions: []PaypalTransaction{
			{
				Amount: PaypalAmount{
					Total:    amount,
					Currency: currency,
				},
				Description: description,
			},
		},
	}
	resp, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(resp)
	req, err := http.NewRequest(http.MethodPost, baseURL+paypalCreatePaymentURL, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+paypalToken)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	paymentResponse := &PaypalPayment{}
	err = json.Unmarshal(body, paymentResponse)
	if err != nil {
		return nil, err
	}
	return paymentResponse, nil
}

// GetPaypalToken returns a new paypal oauth token
func GetPaypalToken(baseURL, public, private string) (*PaypalToken, error) {
	buffer := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest(http.MethodPost, baseURL+paypalTokenURL, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")
	req.SetBasicAuth(public, private)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	tokenResponse := &PaypalToken{}
	err = json.Unmarshal(body, tokenResponse)
	if err != nil {
		return nil, err
	}
	return tokenResponse, nil
}
