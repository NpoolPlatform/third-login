package auth

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// build url with param
type URLBuilder struct {
	baseURL string
	params  url.Values
}

func NewURLBuilder(baseURL string) *URLBuilder {
	uv, err := url.ParseRequestURI(baseURL)
	builder := &URLBuilder{}
	if err != nil {
		log.Println(err)
		return builder
	}
	urls := strings.SplitN(uv.String(), "?", 2)
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
