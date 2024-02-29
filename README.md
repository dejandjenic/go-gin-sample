# go-gin-sample

application should fulfill requirements

1. use firestore as main storage
2. authenticate with jwt obtained from oidc server
3. expose swagger client with swagger.json
4. should load configuration from file or env variable
5. should host prometheus server and add automatic performance counters
6. should be able to allow adding additional prometheus counters
7. should log using structured json logging
8. should work with otel tracing and should respond in response header with current request trace id
9. should be able to load configuration value from GCP secret 
10. have functional (e2e) tests with mocked oidc configuration
11. should expose todo CRUD operations