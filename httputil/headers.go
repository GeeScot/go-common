package httputil

import "fmt"

type Headers map[string]string

func (h Headers) FormURLEncoded() {
	h["Content-Type"] = "application/x-www-form-urlencoded"
}

func (h Headers) JSON() {
	h["Content-Type"] = "application/json"
}

func (h Headers) Authorization(tokenType string, accessToken string) {
	h["Authorization"] = fmt.Sprintf("%s %s", tokenType, accessToken)
}
