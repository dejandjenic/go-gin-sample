package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func getConfigurationBody(u string) string {
	b, _ := os.ReadFile("config.json")
	return strings.ReplaceAll(string(b), "{server}", u)
}

func setupMockServer() (*httptest.Server, string) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.URL.Path == "/realms/DEMOREALM/.well-known/openid-configuration" {
			b := getConfigurationBody(r.Host)
			w.Write([]byte(b))
		} else {
			b, _ := os.ReadFile("keys.json")
			w.Write(b)
		}
	}))

	return server, server.URL + "/realms/DEMOREALM"
}
