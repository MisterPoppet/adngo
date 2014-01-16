package adngo

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseURI = "https://alpha-api.app.net/"
	authURI = "https://account.app.net/oauth/"
)

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

func (a *App) Do(method, url, bodyType string, data Values) (resp *Response, err error) {
	req := http.NewRequest(method, url, bytes.NewBufferString(data.Encode()))

	if a.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+a.accessToken)
	}
	if bodyType != "" {
		req.Header.Add("Content-Type", bodyType)
	}

	return httpClient.Do(req)
}

func (a *App) Get(url, bodyType string) (resp *Response, err error) {
	return a.Do("GET", url, bodyType, url.Values{})
}

func (a *App) Post(url string, bodyType string, data Values) (resp *Response, err error) {
	return a.Do("POST", url, bodyType, data)
}

func (a *App) Put(url string, bodyType string, data Values) (resp *Response, err error) {
	return a.Do("PUT", url, bodyType, data)
}

func (a *App) Patch(url string, bodyType string, data Values) (resp *Response, err error) {
	return a.Do("PATCH", url, bodyType, data)
}

func (a *App) Delete(url string) (resp *Response, err error) {
	return a.Do("DELETE", url, bodyType, url.Values{})
}

func (a *App) VerifyToken(delegate bool) {
	if delegate {
		auth := []byte(a.clientId + ":" + a.clientSecret)
		req := http.NewRequest("GET", baseURI+"stream/0/token", nil)
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
		req.Header.Add("Identity-Delegate-Token", "True")

		resp, err := httpClient.Do(req)
	} else {
		resp, err := a.Get(baseURI+"stream/0/token", "application/json")
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

		resp, err := a.Post(authURI+"access_token", "", data)
	}
}

func (a *App) ProcessText(text string) {
	data := url.Values{}
	data.Add("text", text)

	resp, err := a.Post(baseURI+"stream/0/text/process", "", data)
}
