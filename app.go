package adngo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
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

func (s Scopes) Spaced() string {
	return strings.Join(s, " ")
}

func (s Scopes) String() string {
	return strings.Join(s, ",")
}

// A custom type that satisfies the io.ReadCloser needed by the http Request
type dataCloser struct {
	io.Reader
}

func (dataCloser) Close() error { return nil }

// Our primary API struct. It's the source of all our awesome.
type App struct {
	clientId     string
	clientSecret string
	accessToken  string
	RedirectURI  string
	Scopes       Scopes
}

var httpClient = &http.Client{}

func (a *App) do(method, url, bodyType string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if data != nil {
		req.Body = dataCloser{bytes.NewBufferString(data.Encode())}
	}
	if a.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+a.accessToken)
	}
	if bodyType != "" {
		req.Header.Add("Content-Type", bodyType)
	}

	return httpClient.Do(req)
}

func (a *App) get(url, bodyType string) (resp *http.Response, err error) {
	return a.do("GET", url, bodyType, nil)
}

func (a *App) post(url string, bodyType string, data url.Values) (resp *http.Response, err error) {
	return a.do("POST", url, bodyType, data)
}

func (a *App) put(url string, bodyType string, data url.Values) (resp *http.Response, err error) {
	return a.do("PUT", url, bodyType, data)
}

func (a *App) patch(url string, bodyType string, data url.Values) (resp *http.Response, err error) {
	return a.do("PATCH", url, bodyType, data)
}

func (a *App) delete(url string) (resp *http.Response, err error) {
	return a.do("DELETE", url, "application/json", nil)
}

// Do we even need this??
func (a *App) VerifyToken(delegate bool) *http.Response {
	if delegate {
		auth := []byte(a.clientId + ":" + a.clientSecret)
		req, err := http.NewRequest("GET", baseURI+"stream/0/token", nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(auth))
		req.Header.Add("Identity-Delegate-Token", "True")

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		return resp
	} else {
		resp, err := a.get(baseURI+"stream/0/token", "application/json")
		if err != nil {
			log.Fatal(err)
		}

		return resp
	}
}

func (a *App) AuthURI(clientSide, appStore bool) (uri string) {
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

func (a *App) GetAccessToken(code string, app bool) *http.Response {
	if app {
		data := url.Values{}
		data.Add("client_id", a.clientId)
		data.Add("client_secret", a.clientSecret)
		data.Add("grant_type", "client_credentials")

		resp, err := a.post(authURI+"access_token", "", data)
		if err != nil {
			log.Fatal(err)
		}

		return resp
	}

	return nil
}

func (a *App) ProcessText(text string) *http.Response {
	data := url.Values{}
	data.Add("text", text)

	resp, err := a.post(baseURI+"stream/0/text/process", "", data)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

// Retrieves the App.Net Configuration Object
func (a *App) GetConfig() (config interface{}) {
	resp, err := a.get(baseURI+"stream/0/config", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	var conf interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
