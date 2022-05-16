package auth

import (
	"fmt"
	"net/url"
	"strings"
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
	var n = 2
	urls := strings.SplitN(uv.String(), "?", n)
	builder.baseURL = urls[0]
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
