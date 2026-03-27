package zinc_sdk

import (
	"net/http"
	"strings"
)

func buildHostUrl(host string, authToken string) string {
	if authToken != "" {
		if !strings.HasPrefix(host, "https://") {
			host = "https://" + host
		}
	} else {
		if !strings.HasPrefix(host, "http://") {
			host = "http://" + host
		}
	}
	return host
}

func setAuthHeader(req *http.Request, authToken string) {
	if authToken != "" {
		req.Header.Set("auth-token", authToken)
	}
}
