package application

import (
	"context"
	"os"
	"reflect"

	"github.com/dejandjenic/go-gin-sample/application/configuration"
	"github.com/dejandjenic/go-gin-sample/application/repository"
	secrets "github.com/dejandjenic/go-gin-sample/application/secret-manager"

	"cloud.google.com/go/firestore"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Application struct {
	Configuration  configuration.Config
	TodoRepository repository.ITodoRepository
}

func NewApplication(ctx context.Context, oidcUrl string) Application {
	cfg := loadConfiguration()
	if oidcUrl != "" {
		cfg.RealmConfigURL = oidcUrl
	}
	app := Application{
		Configuration: cfg,
		TodoRepository: repository.TodoRepository{
			Db:     createClient(ctx, cfg.GcpProjectID),
			Config: cfg,
		},
	}
	return app
}

func createClient(ctx context.Context, projectID string) *firestore.Client {

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create client: %v")
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func detectEnvironment() string {
	env := os.Getenv("GO_ENV")
	if "" == env {
		env = "development"
	}
	return env
}

func loadConfiguration() configuration.Config {
	e := detectEnvironment()
	c := configuration.Config{}

	godotenv.Load(".env." + e + ".local")
	if "test" != e {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + e)
	godotenv.Load()
	if err := env.ParseWithOptions(&c,
		env.Options{
			FuncMap: map[reflect.Type]env.ParserFunc{
				reflect.TypeFor[configuration.Secret](): parseSecret,
			},
		},
	); err != nil {
		log.Fatal().Msg("configuration could not be loaded")
	}

	log.Info().Any("value", c).Msg("configuration value")
	return c
}

func parseSecret(v string) (interface{}, error) {
	s := secrets.GetSecret(v)
	return configuration.Secret(s), nil
}
