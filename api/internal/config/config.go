package config

import (
	"os"
	"strconv"
)

var Config *config

const SessionCookieName string = "PSID"

const OAuth2RedirectURI string = ""
const OAuth2ClientID string = ""
const OAUth2ClientSecret string = ""

type session struct {
	CookieDomain string
	CookieSecure bool
}

type config struct {
	ApiAddress      string
	ServerHost      string
	ServerPortRange string
	Runtime         string
	RuntimeUri      string

	Session *session
}

func New() {
	apiAddress := os.Getenv("PANEL_ADDRESS")
	if apiAddress == "" {
		apiAddress = "0.0.0.0:8000"
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

	Config = &config{
		ApiAddress:      apiAddress,
		ServerHost:      serverHost,
		ServerPortRange: serverPortRange,
		Runtime:         runtime,
		RuntimeUri:      runtimeUri,
		Session:         session,
	}
}
