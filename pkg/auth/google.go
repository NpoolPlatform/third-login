package auth

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type GoogleAuth struct {
	GoogleAuthorizeURL string
	GoogleTokenURL     string
	GoogleUserInfoURL  string
}

type GoogleTokenRes struct {
	AccessToken string `json:"access_token"`
}

type GoogleUserInfoRes struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type GoogleErrRes struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (a *GoogleAuth) GetRedirectURL(config *Config) (string, error) {
	url := NewURLBuilder(a.GoogleAuthorizeURL).
		AddParam("client_id", config.ClientID).
		AddParam("redirect_uri", config.RedirectURL).
		AddParam("response_type", "code").
		AddParam("scope", "https://www.googleapis.com/auth/userinfo.email").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}

func (a *GoogleAuth) GetAccessToken(ctx context.Context, code string, config *Config) (string, error) {
	url := NewURLBuilder(a.GoogleTokenURL).
		AddParam("client_id", config.ClientID).
		AddParam("client_secret", config.ClientSecret).
		AddParam("grant_type", "authorization_code").
		AddParam("redirect_uri", config.RedirectURL).
		Build()
	// google redirect code is url encode,addParam will cause duplication url encode
	url = url + "&code=" + code
	client := resty.New()
	if os.Getenv("ENV_CURRENCY_REQUEST_PROXY") != "" {
		client.SetProxy(os.Getenv("ENV_CURRENCY_REQUEST_PROXY"))
	}
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Post(url)
	if err != nil {
		return "", err
	}
	successCode := 200
	if resp.StatusCode() != successCode {
		googleRes := GoogleErrRes{}
		err = json.Unmarshal(resp.Body(), &googleRes)
		if err != nil {
			return "", err
		}
		return "", errors.New(googleRes.ErrorDescription)
	}
	googleRes := GoogleTokenRes{}
	err = json.Unmarshal(resp.Body(), &googleRes)
	if err != nil {
		return "", err
	}
	return googleRes.AccessToken, nil
}

func (a *GoogleAuth) GetUserInfo(ctx context.Context, code string, config *Config) (*appusermgrpb.AppUserThirdParty, error) {
	token, err := a.GetAccessToken(ctx, code, config)
	if err != nil {
		return &appusermgrpb.AppUserThirdParty{}, err
	}
	url := a.GoogleUserInfoURL
	client := resty.New()
	if os.Getenv("ENV_CURRENCY_REQUEST_PROXY") != "" {
		client.SetProxy(os.Getenv("ENV_CURRENCY_REQUEST_PROXY"))
	}
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		Get(url)
	if err != nil {
		return &appusermgrpb.AppUserThirdParty{}, err
	}
	successCode := 200
	if resp.StatusCode() != successCode {
		googleRes := GoogleErrRes{}
		err = json.Unmarshal(resp.Body(), &googleRes)
		if err != nil {
			return &appusermgrpb.AppUserThirdParty{}, err
		}
		return &appusermgrpb.AppUserThirdParty{}, errors.New(googleRes.ErrorDescription)
	}
	googleRes := GoogleUserInfoRes{}
	err = json.Unmarshal(resp.Body(), &googleRes)
	if err != nil {
		return nil, err
	}
	return &appusermgrpb.AppUserThirdParty{
		ThirdPartyUserID:     googleRes.ID,
		ThirdPartyUsername:   googleRes.Email,
		ThirdPartyUserAvatar: googleRes.Picture,
		ThirdPartyID:         config.ClientID,
	}, nil
}
