package auth

import (
	"context"
	"encoding/json"
	"fmt"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var ThirdMap = make(map[string]ThirdMethod)

func init() {
	ThirdMap[appusermgrconst.ThirdGithub] = &GitHubAuth{}
	ThirdMap[appusermgrconst.ThirdGoogle] = &GoogleAuth{}
}

type ThirdMethod interface {
	GetUserInfo(ctx context.Context, code string, config *Config) (*appusermgrpb.AppUserThird, error)
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
