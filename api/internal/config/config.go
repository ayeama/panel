package config

import (
	"errors"
	"os"
	"strconv"
)

var Config *config

const SessionCookieName string = "PSID"

type session struct {
	CookieDomain string
	CookieSecure bool
}

type oauth2 struct {
	ClientId     string
	RedirectURI  string
	ClientSecret string
}

type config struct {
	ApiAddress      string
	Domain          string
	ServerHost      string
	ServerPortRange string
	Runtime         string
	RuntimeUri      string

	OAuth2  *oauth2
	Session *session
}

func New() {
	apiAddress := os.Getenv("PANEL_ADDRESS")
	if apiAddress == "" {
		apiAddress = "0.0.0.0:8000"
	}

	domain := os.Getenv("PANEL_DOMAIN")
	if domain == "" {
		domain = "http://localhost:5173"
	}

	serverHost := os.Getenv("PANEL_SERVER_HOST")
	if serverHost == "" {
		serverHost = "localhost"
	}

	serverPortRange := os.Getenv("PANEL_SERVER_PORT_RANGE")
	if serverPortRange == "" {
		serverPortRange = "45000-45099"
	}

	runtime := os.Getenv("PANEL_RUNTIME")
	if runtime == "" {
		runtime = "docker" // TODO default should be podman
	}

	runtimeUri := os.Getenv("PANEL_RUNTIME_URI")
	if runtimeUri == "" {
		runtimeUri = "unix:/run/user/1000/podman/podman.sock"
	}

	oauth2ClientId := os.Getenv("PANEL_OAUTH2_CLIENT_ID")
	if oauth2ClientId == "" {
		panic(errors.New("missing oauth2 client id"))
	}

	oauth2ClientSecret := os.Getenv("PANEL_OAUTH2_CLIENT_SECRET")
	if oauth2ClientSecret == "" {
		panic(errors.New("missing oauth2 client secret"))
	}

	oauth2RedirectURI := os.Getenv("PANEL_OAUTH2_REDIRECT_URI")
	if oauth2RedirectURI == "" {
		oauth2RedirectURI = "http://localhost:8000/oauth2/callback"
	}

	sessionCookieDomain := os.Getenv("PANEL_SESSION_COOKIE_DOMAIN")
	if sessionCookieDomain == "" {
		sessionCookieDomain = "localhost"
	}

	_sessionCookieSecure := os.Getenv("PANEL_SESSION_COOKIE_SECURE")
	if _sessionCookieSecure == "" {
		_sessionCookieSecure = "true"
	}
	sessionCookieSecure, err := strconv.ParseBool(_sessionCookieSecure)
	if err != nil {
		panic(err)
	}

	session := &session{
		CookieDomain: sessionCookieDomain,
		CookieSecure: sessionCookieSecure,
	}

	oauth2 := &oauth2{
		ClientId:     oauth2ClientId,
		ClientSecret: oauth2ClientSecret,
		RedirectURI:  oauth2RedirectURI,
	}

	Config = &config{
		ApiAddress:      apiAddress,
		Domain:          domain,
		ServerHost:      serverHost,
		ServerPortRange: serverPortRange,
		Runtime:         runtime,
		RuntimeUri:      runtimeUri,
		OAuth2:          oauth2,
		Session:         session,
	}
}
