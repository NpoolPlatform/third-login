package auth

import (
	"fmt"
	"net/url"
)

type URLBuilder struct {
	baseURL string
	params  url.Values
}

func NewURLBuilder(baseURL string) *URLBuilder {
	uv, err := url.ParseRequestURI(baseURL)
	builder := &URLBuilder{}
	if err != nil {
		return builder
	}
	builder.baseURL = uv.Scheme + "://" + uv.Host + uv.Path
	builder.params = uv.Query()
	return builder
}

func (c *URLBuilder) AddParam(key string, value interface{}) *URLBuilder {
	if key == "" {
		return c
	}
	c.params.Add(key, fmt.Sprint(value))
	return c
}

func (c *URLBuilder) Build() string {
	if c.baseURL == "" {
		return ""
	}
	if len(c.params) == 0 {
		return c.baseURL
	}
	return c.baseURL + "?" + c.params.Encode()
}
