package authorization

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	oidc "github.com/coreos/go-oidc"
)

type Claims struct {
	JTI   string `json:"jti,omitempty"`
	Scope string `json:"scope,omitempty"`
}

func authorisationFailed(message string) {
	fmt.Println("error %v", message)
}

func Authorize(rawAccessToken string, realmConfigurationUrl string) (Claims, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Duration(6000) * time.Second,
		Transport: tr,
	}
	ctx := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ctx, realmConfigurationUrl)
	if err != nil {
		authorisationFailed("authorisation failed while getting the provider: " + err.Error())
		return Claims{}, err
	}

	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}
	verifier := provider.Verifier(oidcConfig)
	idToken, err := verifier.Verify(ctx, rawAccessToken)
	if err != nil {
		authorisationFailed("authorisation failed while verifying the token: " + err.Error())
		return Claims{}, err
	}

	var IDTokenClaims Claims // ID Token payload is just JSON.
	if err := idToken.Claims(&IDTokenClaims); err != nil {
		authorisationFailed("claims : " + err.Error())
		return Claims{}, err
	}
	return IDTokenClaims, nil
}
