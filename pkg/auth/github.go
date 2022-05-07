package auth

import (
	"github.com/NpoolPlatform/third-login-gateway/pkg/utils"
	"github.com/google/uuid"
)

type GitHubAuth struct {
	BaseRequest
}

func NewGitHubAuth(conf *Config) *GitHubAuth {
	authRequest := &GitHubAuth{}
	authRequest.Set(conf)

	authRequest.authorizeURL = "https://github.com/login/oauth/authorize"
	authRequest.TokenURL = "https://github.com/login/oauth/access_token"
	authRequest.userInfoURL = "https://api.github.com/user"

	return authRequest
}

func (a *GitHubAuth) GetRedirectURL() (string, error) {
	url := utils.NewURLBuilder(a.authorizeURL).
		AddParam("response_type", "code").
		AddParam("client_id", a.config.ClientID).
		AddParam("redirect_uri", a.config.RedirectURL).
		AddParam("scope", "snsapi_login").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}
