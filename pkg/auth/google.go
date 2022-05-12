package auth

import (
	"context"
	"encoding/json"
	"errors"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	"github.com/go-resty/resty/v2"
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
	url := NewURLBuilder(a.authorizeURL).
		AddParam("client_id", a.config.ClientID).
		AddParam("redirect_uri", a.config.RedirectURL).
		AddParam("response_type", "code").
		AddParam("scope", "https://www.googleapis.com/auth/userinfo.email").
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}

func (a *GoogleAuth) GetAccessToken(ctx context.Context, code string) (string, error) {
	url := NewURLBuilder(a.TokenURL).
		AddParam("client_id", a.config.ClientID).
		AddParam("client_secret", a.config.ClientSecret).
		AddParam("grant_type", "authorization_code").
		AddParam("redirect_uri", a.config.RedirectURL).
		Build()
	// google redirect code is url encode,addParam will cause duplication url encode
	url = url + "&code=" + code
	client := resty.New()
	client.SetProxy("http://192.168.31.135:7890") // update to ENV
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
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

func (a *GoogleAuth) GetUserInfo(ctx context.Context, code string) (*appusermgrpb.AppUserThird, error) {
	token, err := a.GetAccessToken(ctx, code)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}
	url := a.userInfoURL
	client := resty.New()
	client.SetProxy("http://192.168.31.135:7890") // update to ENV
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		Get(url)
	if err != nil {
		return &appusermgrpb.AppUserThird{}, err
	}

	m := JsonToMSS(string(resp.Body()))
	if _, ok := m["error"]; ok {
		return &appusermgrpb.AppUserThird{}, errors.New(m["error_description"])
	}
	return &appusermgrpb.AppUserThird{
		ThirdUserId:      m["id"],
		ThirdUserName:    m["email"],
		ThirdUserPicture: m["picture"],
		ThirdExtra:       string(resp.Body()),
	}, nil
}
