package configuration

type Config struct {
	RealmConfigURL   string   `env:"RealmConfigURL"`
	Secret           Secret   `env:"Secret"`
	GcpProjectID     string   `env:"GcpProjectID"`
	AuthExclusionUrl []string `env:"AuthExclusionUrl"  envSeparator:","`
	FirestorePrefix  string   `env:"FirestorePrefix"`
}

type Secret string
