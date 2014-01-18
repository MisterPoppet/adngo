package adngo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Our API Urls
const (
	baseURI = "https://alpha-api.app.net/"
	authURI = "https://account.app.net/oauth/"
)

// This is our scopes struct to check for that.
type Scopes []string

func (s *Scopes) Spaced() string {
	return strings.Join(s, " ")
}

func (s *Scopes) String() string {
	return strings.Join(s, ",")
}

type App struct {
	clientId     string
	clientSecret string
	accessToken  string
	RedirectURI  string
	Scopes       Scopes
}

var httpClient = &http.Client{}

func (a *App) do(method, url, bodyType string, data Values) (resp *Response, err error) {
	if data == nil {
		req := http.NewRequest(method, url, nil)
	} else {
		req := http.NewRequest(method, url, bytes.NewBufferString(data.Encode()))
	}

	if a.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+a.accessToken)
	}
	if bodyType != "" {
		req.Header.Add("Content-Type", bodyType)
	}

	return httpClient.do(req)
}

func (a *App) get(url, bodyType string) (resp *Response, err error) {
	return a.do("GET", url, bodyType, nil)
}

func (a *App) post(url string, bodyType string, data Values) (resp *Response, err error) {
	return a.do("POST", url, bodyType, data)
}

func (a *App) put(url string, bodyType string, data Values) (resp *Response, err error) {
	return a.do("PUT", url, bodyType, data)
}

func (a *App) patch(url string, bodyType string, data Values) (resp *Response, err error) {
	return a.do("PATCH", url, bodyType, data)
}

func (a *App) delete(url string) (resp *Response, err error) {
	return a.do("DELETE", url, bodyType, nil)
}

// Do we even need this??
func (a *App) VerifyToken(delegate bool) {
	if delegate {
		auth := []byte(a.clientId + ":" + a.clientSecret)
		req := http.NewRequest("GET", baseURI+"stream/0/token", nil)
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
		req.Header.Add("Identity-Delegate-Token", "True")

		resp, err := httpClient.Do(req)
	} else {
		resp, err := a.get(baseURI+"stream/0/token", "application/json")
	}
}

func (a *App) AuthURI(clientSide, appStore bool) {
	data := url.Values{}
	data.Add("client_id", a.clientId)
	data.Add("redirect_uri", a.RedirectURI)
	data.Add("scope", a.Scopes.String())

	if clientSide {
		data.Add("response_type", "token")
	} else {
		data.Add("response_type", "code")
	}
	if appStore {
		data.Add("adnview", "appstore")
	}

	return authURI + "authenticate?" + data.Encode()
}

func (a *App) GetAccessToken(code string, app bool) {
	if app {
		data := url.Values{}
		data.Add("client_id", a.clientId)
		data.Add("client_secret", a.clientSecret)
		data.Add("grant_type", "client_credentials")

		resp, err := a.post(authURI+"access_token", "", data)
	}
}

func (a *App) ProcessText(text string) {
	data := url.Values{}
	data.Add("text", text)

	resp, err := a.post(baseURI+"stream/0/text/process", "", data)
}

// Retrieves the App.Net Configuration Object
func (a *App) GetConfig() {
	resp, err := a.get(baseURI+"stream/0/config", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	var config interface{}
	err = json.Unmarshal(resp, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(config["meta"]["code"])
	return config
}

// Retrieves the App.Net Configuration Object
func (a *App) GetConfig() {
	resp, err := a.Get(baseURI+"stream/0/config", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	var config interface{}
	err = json.Unmarshal(resp, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(config['meta']['code'])
	return config
}
