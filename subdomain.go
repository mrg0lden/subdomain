package subdomain

import (
	"net/http"
	"strings"

	"github.com/go-chi/hostrouter"
)

//Routes embeds github.com/go-chi/hostrouter's Routes
type Routes struct{ hostrouter.Routes }

//New creates an instance
func New() Routes {
	return Routes{}
}

//ServeHTTP implements the http.Handler interface
func (sr Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := requestHost(r)
	for subdomain, router := range sr.Routes {

		if strings.HasPrefix(host, subdomain) {
			router.ServeHTTP(w, r)
			return
		}
	}

	if router, ok := sr.Routes["*"]; ok {
		router.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

//Copyright (c) 2016-Present https://github.com/go-chi authors
// MIT License
func requestHost(r *http.Request) (host string) {
	// not standard, but most popular
	host = r.Header.Get("X-Forwarded-Host")
	if host != "" {
		return
	}

	// RFC 7239
	host = r.Header.Get("Forwarded")
	_, _, host = parseForwarded(host)
	if host != "" {
		return
	}

	// if all else fails fall back to request host
	host = r.Host
	return
}

//Copyright (c) 2016-Present https://github.com/go-chi authors
// MIT License
func parseForwarded(forwarded string) (addr, proto, host string) {
	if forwarded == "" {
		return
	}
	for _, forwardedPair := range strings.Split(forwarded, ";") {
		if tv := strings.SplitN(forwardedPair, "=", 2); len(tv) == 2 {
			token, value := tv[0], tv[1]
			token = strings.TrimSpace(token)
			value = strings.TrimSpace(strings.Trim(value, `"`))
			switch strings.ToLower(token) {
			case "for":
				addr = value
			case "proto":
				proto = value
			case "host":
				host = value
			}

		}
	}
	return
}
