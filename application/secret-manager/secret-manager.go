package secrets

import (
	"context"

	"github.com/rs/zerolog/log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func GetSecret(name string) string {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)

	log.Info().Str("name", name).Msg("getting the secret")

	if err != nil {
		log.Fatal().Msgf("failed to setup client: %v", err)
		return "ERROR"
	}

	defer client.Close()

	res, err := client.AccessSecretVersion(context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})

	if err != nil {
		log.Fatal().Msgf("failed to access secret: %v", err)
		return "ERROR"
	}

	return string(res.Payload.Data)
}
