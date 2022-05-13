package auth

import (
	"context"
	"encoding/json"
	"errors"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

var (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserInfoURL  = "https://api.github.com/user"
)

type GitHubAuth struct{}

func (a *GitHubAuth) GetRedirectURL(config *Config) (string, error) {
	url := NewURLBuilder(githubAuthorizeURL).
		AddParam("response_type", "code").
		AddParam("client_id", config.ClientID).
		AddParam("redirect_uri", config.RedirectURL).
		AddParam("scope", "snsapi_login").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}

func (a *GitHubAuth) GetAccessToken(ctx context.Context, code string, config *Config) (string, error) {
	url := NewURLBuilder(githubTokenURL).
		AddParam("client_id", config.ClientID).
		AddParam("client_secret", config.ClientSecret).
		AddParam("code", code).
		Build()
	client := resty.New()
	client.SetProxy("http://192.168.31.135:7890") // update to ENV
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Post(url)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &m)
	if err != nil {
		return "", err
	}
	if _, ok := m["error"]; ok {
		return "", errors.New(m["error_description"].(string))
	}
	return m["access_token"].(string), err
}

func (a *GitHubAuth) GetUserInfo(ctx context.Context, code string, config *Config) (*appusermgrpb.AppUserThird, error) {
	token, err := a.GetAccessToken(ctx, code, config)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}
	url := githubUserInfoURL

	client := resty.New()
	client.SetProxy("http://192.168.31.135:7890") // update to ENV
	resp, err := client.R().
		SetContext(ctx).
		// batter is use Bearer
		SetAuthToken(token).
		Get(url)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}
	m, err := JSONToMSS(string(resp.Body()))
	if err != nil {
		return nil, err
	}
	if _, ok := m["error"]; ok {
		return &appusermgrpb.AppUserThird{}, errors.New(m["error_description"])
	}
	return &appusermgrpb.AppUserThird{
		ThirdUserId:      m["id"],
		ThirdUserName:    m["login"],
		ThirdUserPicture: m["avatar_url"],
		ThirdExtra:       string(resp.Body()),
	}, nil
}
