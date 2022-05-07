package auth

//基本配置
type AuthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type BaseRequest struct {
	authorizeUrl   string
	TokenUrl       string
	AccessTokenUrl string
	RefreshUrl     string
	userInfoUrl    string
	config         *AuthConfig
}

func (b *BaseRequest) Set(cfg *AuthConfig) {
	b.config = cfg
}

type CodeResult struct {
	Code int `json:"code"`
}

type TokenResult struct {
}

type UserResult struct {
}
