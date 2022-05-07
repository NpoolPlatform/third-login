package auth

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type BaseRequest struct {
	authorizeURL string //nolint
	TokenURL     string
	userInfoURL  string //nolint
	config       *Config
}

func (b *BaseRequest) Set(cfg *Config) {
	b.config = cfg
}

type CodeResult struct {
	Code int `json:"code"`
}
