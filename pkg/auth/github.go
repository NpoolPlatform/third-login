package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	appuserconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type GitHubAuth struct {
	GithubAuthorizeURL string
	GithubTokenURL     string
	GithubUserInfoURL  string
}

type GitHubUserInfoRes struct {
	ID               int    `json:"id"`
	Login            string `json:"login"`
	AvatarURL        string `json:"avatar_url"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type GitHubTokenRes struct {
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (a *GitHubAuth) GetRedirectURL(config *Config) (string, error) {
	url := NewURLBuilder(a.GithubAuthorizeURL).
		AddParam("response_type", "code").
		AddParam("client_id", config.ClientID).
		AddParam("redirect_uri", config.RedirectURL).
		AddParam("scope", "snsapi_login").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}

func (a *GitHubAuth) GetAccessToken(ctx context.Context, code string, config *Config) (string, error) {
	url := NewURLBuilder(a.GithubTokenURL).
		AddParam("client_id", config.ClientID).
		AddParam("client_secret", config.ClientSecret).
		AddParam("code", code).
		Build()
	client := resty.New()

	client.SetProxy(os.Getenv("ENV_CURRENCY_REQUEST_PROXY"))
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Post(url)
	if err != nil {
		return "", err
	}
	gitHubRes := GitHubTokenRes{}
	err = json.Unmarshal(resp.Body(), &gitHubRes)
	if err != nil {
		return "", err
	}
	if gitHubRes.Error != "" {
		return "", errors.New(gitHubRes.ErrorDescription)
	}
	return gitHubRes.AccessToken, err
}

func (a *GitHubAuth) GetUserInfo(ctx context.Context, code string, config *Config) (*appusermgrpb.AppUserThird, error) {
	token, err := a.GetAccessToken(ctx, code, config)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}
	url := a.GithubUserInfoURL

	client := resty.New()
	client.SetProxy(os.Getenv("ENV_CURRENCY_REQUEST_PROXY"))
	resp, err := client.R().
		SetContext(ctx).
		// batter is use Bearer
		SetAuthToken(token).
		Get(url)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}

	gitHubRes := GitHubUserInfoRes{}
	err = json.Unmarshal(resp.Body(), &gitHubRes)
	if err != nil {
		return nil, err
	}
	if gitHubRes.Error != "" {
		return &appusermgrpb.AppUserThird{}, errors.New(gitHubRes.ErrorDescription)
	}
	return &appusermgrpb.AppUserThird{
		ThirdUserID:     fmt.Sprintf("%v", gitHubRes.ID),
		ThirdUserName:   gitHubRes.Login,
		ThirdUserAvatar: gitHubRes.AvatarURL,
		Third:           appuserconst.ThirdGithub,
		ThirdID:         config.ClientID,
	}, nil
}
