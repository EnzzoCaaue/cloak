package util

import (
    "bytes"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

// PaypalToken is the main oauth struct for paypal
type PaypalToken struct {
    Scope string `json:"scope"`
    Token string `json:"access_token"`
    Type string `json:"token_type"`
    ExpiresIn int64 `json:"expires_in"`
}

// GetPaypalToken returns a new paypal oauth token
func GetPaypalToken(baseURL, public, private string) (*PaypalToken, error) {
    buffer := bytes.NewBuffer([]byte("grant_type=client_credentials"))
    req, err := http.NewRequest(http.MethodPost, baseURL + "/v1/oauth2/token", buffer)
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