package auth

import (
	"github.com/NpoolPlatform/third-login-gateway/pkg/utils"
	"github.com/google/uuid"
)

type AuthGitHub struct {
	BaseRequest
}

func NewAuthGitHub(conf *AuthConfig) *AuthGitHub {
	authRequest := &AuthGitHub{}
	authRequest.Set(conf)

	authRequest.authorizeUrl = "https://github.com/login/oauth/authorize"
	authRequest.TokenUrl = "https://github.com/login/oauth/access_token"
	authRequest.userInfoUrl = "https://api.github.com/user"

	return authRequest
}

func (a *AuthGitHub) GetRedirectUrl() (string, error) {
	url := utils.NewUrlBuilder(a.authorizeUrl).
		AddParam("response_type", "code").
		AddParam("client_id", a.config.ClientId).
		AddParam("redirect_uri", a.config.RedirectUrl).
		AddParam("scope", "snsapi_login").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}
