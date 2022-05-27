package auth

import (
	"context"
	"encoding/json"
	"fmt"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	thirdlogingawconst "github.com/NpoolPlatform/third-login-gateway/pkg/const"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var ThirdMap = make(map[string]ThirdMethod)

func init() {
	ThirdMap[thirdlogingawconst.ThirdPartyGithub] = &GitHubAuth{
		GithubAuthorizeURL: "https://github.com/login/oauth/authorize",
		GithubTokenURL:     "https://github.com/login/oauth/access_token",
		GithubUserInfoURL:  "https://api.github.com/user",
	}
	ThirdMap[thirdlogingawconst.ThirdPartyGoogle] = &GoogleAuth{
		GoogleAuthorizeURL: "https://accounts.google.com/o/oauth2/v2/auth",
		GoogleTokenURL:     "https://oauth2.googleapis.com/token",
		GoogleUserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
	}
	ThirdMap[thirdlogingawconst.ThirdPartyFaceBook] = &FaceBookAuth{
		FaceBookAuthorizeURL: "https://www.facebook.com/v13.0/dialog/oauth",
		FaceBookTokenURL:     "https://graph.facebook.com/v13.0/oauth/access_token",
		FaceBookUserInfoURL:  "https://graph.facebook.com/me",
	}
}

type ThirdMethod interface {
	GetUserInfo(ctx context.Context, code string, config *Config) (*appusermgrpb.AppUserThirdParty, error)
	GetRedirectURL(config *Config) (string, error)
}

type Context struct {
	ThirdMethod
}

func NewContext(thirdMethod ThirdMethod) *Context {
	return &Context{
		thirdMethod,
	}
}

func JSONToMSS(s string) (map[string]string, error) {
	if s == "" {
		return nil, fmt.Errorf("map empty")
	}
	msi := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &msi)
	if err != nil {
		return nil, err
	}
	mss := make(map[string]string)
	for k, v := range msi {
		mss[k] = fmt.Sprintf("%v", v)
	}
	return mss, nil
}
