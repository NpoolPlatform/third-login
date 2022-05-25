package auth

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type FaceBookAuth struct {
	FaceBookAuthorizeURL string
	FaceBookTokenURL     string
	FaceBookUserInfoURL  string
}

type FaceBookUserInfoRes struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture struct {
		Date struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type FaceBookTokenRes struct {
	AccessToken string `json:"access_token"`
}

type FaceBookErrRes struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (a *FaceBookAuth) GetRedirectURL(config *Config) (string, error) {
	url := NewURLBuilder(a.FaceBookAuthorizeURL).
		AddParam("client_id", config.ClientID).
		AddParam("redirect_uri", config.RedirectURL).
		AddParam("state", uuid.New().String()).
		Build()
	return url, nil
}

func (a *FaceBookAuth) GetAccessToken(ctx context.Context, code string, config *Config) (string, error) {
	url := NewURLBuilder(a.FaceBookTokenURL).
		AddParam("client_id", config.ClientID).
		AddParam("client_secret", config.ClientSecret).
		AddParam("code", code).
		Build()
	url = url + "&redirect_uri=" + config.RedirectURL
	client := resty.New()
	if os.Getenv("ENV_CURRENCY_REQUEST_PROXY") != "" {
		client.SetProxy(os.Getenv("ENV_CURRENCY_REQUEST_PROXY"))
	}
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Post(url)
	if err != nil {
		return "", err
	}
	successCode := 200
	if resp.StatusCode() != successCode {
		facebookRes := FaceBookErrRes{}
		err = json.Unmarshal(resp.Body(), &facebookRes)
		if err != nil {
			return "", err
		}
		return "", errors.New(facebookRes.Error.Message)
	}
	faceBookRes := FaceBookTokenRes{}
	err = json.Unmarshal(resp.Body(), &faceBookRes)
	if err != nil {
		return "", err
	}

	return faceBookRes.AccessToken, err
}

func (a *FaceBookAuth) GetUserInfo(ctx context.Context, code string, config *Config) (*appusermgrpb.AppUserThirdParty, error) {
	token, err := a.GetAccessToken(ctx, code, config)
	if err != nil {
		return &appusermgrpb.AppUserThirdParty{}, err
	}
	url := NewURLBuilder(a.FaceBookUserInfoURL).
		AddParam("fields", "id,name,picture").
		AddParam("access_token", token).
		Build()
	client := resty.New()
	if os.Getenv("ENV_CURRENCY_REQUEST_PROXY") != "" {
		client.SetProxy(os.Getenv("ENV_CURRENCY_REQUEST_PROXY"))
	}
	resp, err := client.R().
		SetContext(ctx).
		// batter is use Bearer
		SetAuthToken(token).
		Get(url)
	if err != nil {
		return &appusermgrpb.AppUserThirdParty{}, err
	}
	successCode := 200
	if resp.StatusCode() != successCode {
		facebookRes := FaceBookErrRes{}
		err = json.Unmarshal(resp.Body(), &facebookRes)
		if err != nil {
			return &appusermgrpb.AppUserThirdParty{}, err
		}
		return &appusermgrpb.AppUserThirdParty{}, errors.New(facebookRes.Error.Message)
	}
	faceBookRes := FaceBookUserInfoRes{}
	err = json.Unmarshal(resp.Body(), &faceBookRes)
	if err != nil {
		return nil, err
	}
	return &appusermgrpb.AppUserThirdParty{
		ThirdPartyUserID:     faceBookRes.ID,
		ThirdPartyUserName:   faceBookRes.Name,
		ThirdPartyUserAvatar: faceBookRes.Picture.Date.URL,
		ThirdPartyID:         config.ClientID,
	}, nil
}
