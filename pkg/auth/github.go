package auth

import (
	"context"
	"encoding/json"
	"errors"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	"github.com/go-resty/resty/v2"
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
	url := NewURLBuilder(a.authorizeURL).
		AddParam("response_type", "code").
		AddParam("client_id", a.config.ClientID).
		AddParam("redirect_uri", a.config.RedirectURL).
		AddParam("scope", "snsapi_login").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}

func (a *GitHubAuth) GetAccessToken(ctx context.Context, code string) (string, error) {
	url := NewURLBuilder(a.TokenURL).
		AddParam("client_id", a.config.ClientID).
		AddParam("client_secret", a.config.ClientSecret).
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

func (a *GitHubAuth) GetUserInfo(ctx context.Context, code string) (*appusermgrpb.AppUserThird, error) {
	token, err := a.GetAccessToken(ctx, code)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}
	url := a.userInfoURL

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
	m := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &m)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}
	if _, ok := m["error"]; ok {
		return &appusermgrpb.AppUserThird{}, errors.New(m["error_description"].(string))
	}
	return &appusermgrpb.AppUserThird{
		ThirdUserId:      m["id"].(string),
		ThirdUserName:    m["login"].(string),
		ThirdUserPicture: m["avatar_url"].(string),
		ThirdExtra:       string(resp.Body()),
	}, nil
}
