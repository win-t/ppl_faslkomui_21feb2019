package main

import (
	"encoding/base64"
	"net/http"
	"strings"

	middleware "github.com/payfazz/go-middleware"
	"github.com/payfazz/go-router/defhandler"
)

func basicAuth(username, password string) middleware.Func {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if ok := func() bool {
				authPart := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
				if len(authPart) < 2 {
					return false
				}
				if strings.ToLower(authPart[0]) != "basic" {
					return false
				}
				rawPart, err := base64.StdEncoding.DecodeString(authPart[1])
				if err != nil {
					return false
				}
				basicPart := strings.SplitN(string(rawPart), ":", 2)
				if len(basicPart) < 2 {
					return false
				}
				return basicPart[0] == username && basicPart[1] == password
			}(); !ok {
				defhandler.StatusUnauthorized(w, r)
				return
			}
			next(w, r)
		}
	}
}
