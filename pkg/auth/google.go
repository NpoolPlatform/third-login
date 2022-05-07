package auth

import (
	"github.com/NpoolPlatform/third-login-gateway/pkg/utils"
	"github.com/google/uuid"
)

type AuthGoogle struct {
	BaseRequest
}

func NewAuthGoogle(conf *AuthConfig) *AuthGoogle {
	authRequest := &AuthGoogle{}
	authRequest.Set(conf)

	authRequest.authorizeUrl = "https://accounts.google.com/o/oauth2/v2/auth "
	authRequest.TokenUrl = "https://oauth2.googleapis.com/token"
	authRequest.userInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo"

	return authRequest
}

func (a *AuthGoogle) GetRedirectUrl() (string, error) {
	url := utils.NewUrlBuilder(a.authorizeUrl).
		AddParam("client_id", a.config.ClientId).
		AddParam("redirect_uri", a.config.RedirectUrl).
		AddParam("response_type", "code").
		AddParam("scope", "snsapi_login").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}
