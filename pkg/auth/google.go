package auth

import (
	"github.com/NpoolPlatform/third-login-gateway/pkg/utils"
	"github.com/google/uuid"
)

type GoogleAuth struct {
	BaseRequest
}

func NewGoogleAuth(conf *Config) *GoogleAuth {
	authRequest := &GoogleAuth{}
	authRequest.Set(conf)

	authRequest.authorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	authRequest.TokenURL = "https://oauth2.googleapis.com/token"
	authRequest.userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

	return authRequest
}

func (a *GoogleAuth) GetRedirectURL() (string, error) {
	url := utils.NewURLBuilder(a.authorizeURL).
		AddParam("client_id", a.config.ClientID).
		AddParam("redirect_uri", a.config.RedirectURL).
		AddParam("response_type", "code").
		AddParam("scope", "https://www.googleapis.com/auth/userinfo.email").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}
