package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/service"
)

type OAuth2Handler struct {
	userService    *service.UserService
	sessionService *service.SessionService
}

func NewOAuth2Handler(userService *service.UserService, sessionService *service.SessionService) *OAuth2Handler {
	return &OAuth2Handler{
		userService:    userService,
		sessionService: sessionService,
	}
}

func (h *OAuth2Handler) Callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	if code == "" {
		panic(errors.New("code query parameter is required"))
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", config.OAuth2RedirectURI)

	body := strings.NewReader(form.Encode())
	request, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", body)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(config.OAuth2ClientID, config.OAUth2ClientSecret)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	type OAuth2TokenResponse struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	var data OAuth2TokenResponse
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		panic(err)
	}

	request, err = http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", "Bearer "+data.AccessToken)

	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	type OAuth2MeResponse struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
	}

	var userData OAuth2MeResponse
	if err := json.NewDecoder(response.Body).Decode(&userData); err != nil {
		panic(err)
	}

	user := h.userService.Login(userData.Email)

	session := h.sessionService.Create(user.Id)

	// TODO: tmp sessionCookie
	sessionCookie := &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    session.Id,
		Path:     "/",
		Domain:   config.Config.Session.CookieDomain,
		MaxAge:   3600,
		Secure:   config.Config.Session.CookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, sessionCookie)
	http.Redirect(w, r, "http://localhost:5173", http.StatusMovedPermanently)
}

func (h *OAuth2Handler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("GET /oauth2/callback", h.Callback)
}
