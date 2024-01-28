// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/olahol/melody"

	"senior-software-engineer-takehome/internal/dispatcher"
	"senior-software-engineer-takehome/internal/postgres"
	"senior-software-engineer-takehome/internal/weather"
	"senior-software-engineer-takehome/restapi/operations"
	"senior-software-engineer-takehome/restapi/operations/service"
	"senior-software-engineer-takehome/restapi/operations/weather_assets"
)

//go:generate swagger generate server --target ../../go --name Backend --spec ../spec/openapi.yaml

var dbConfig = struct {
	Name     string `long:"db.name" default:"postgres" env:"DB_NAME" description:""`
	User     string `long:"db.user" default:"postgres" env:"DB_USER" description:""`
	Password string `long:"db.pwd" default:"example" env:"DB_PWD" description:""`
	Address  string `long:"db.address" default:"db:5432" env:"DB_ADDRESS" description:""`
}{}

func configureFlags(api *operations.BackendAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{{
		LongDescription:  "MongoDB Database options",
		ShortDescription: "MongoDB",
		Options:          &dbConfig,
	}}
}

func configureAPI(api *operations.BackendAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	db := postgres.New(api.Logger, dbConfig.Name, dbConfig.User, dbConfig.Password, dbConfig.Address)

	dispatcher := dispatcher.New(melody.New(), api.Logger)
	go dispatcher.Dispatch()

	weather := weather.Weather{
		Db:         db,
		Dispatcher: dispatcher,
		Log:        api.Logger,
	}

	api.WeatherAssetsGetDataWeatherHandler = weather_assets.GetDataWeatherHandlerFunc(weather.GetWeatherUnits)
	api.ServiceGetWsHandler = service.GetWsHandlerFunc(dispatcher.HandleWS)
	api.WeatherAssetsPostDataWeatherHandler = weather_assets.PostDataWeatherHandlerFunc(weather.AddWeatherUnit)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
