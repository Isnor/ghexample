package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Isnor/ghexample/config"
	"github.com/Isnor/ghexample/controller"
	"github.com/Isnor/ghexample/model"
	"github.com/Isnor/ghexample/service"

	"github.com/Isnor/generichandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	// configure
	usConfig := config.ServiceConfig{}
	envconfig.MustProcess("usgen", &usConfig)

	log.Warn().Msgf("config: %+v", usConfig)
	var dataService service.DataService

	if len(usConfig.Datastore.PostgreSQL.Host) != 0 {
		// try to connect to postgres given the configuration
		pgConfig := usConfig.Datastore.PostgreSQL
		conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s/%s", pgConfig.Username, pgConfig.Password, pgConfig.Host, pgConfig.Database))
		if err != nil {
			log.Fatal().Err(err).Msg("Could not connect to postgres")
			return
		}
		defer conn.Close(ctx)

		dataService = &service.PostgresService{
			PostgresConfig: pgConfig,
			DB:             conn,
		}
	}
	//

	// a simple example controller that uses some data store
	mockWithDatabaseController := &controller.ExampleControllerWithDatabase{
		DataService: dataService,
	}

	// create router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// this is how we would write the GetDataHandler as a typical http.Handler - see controller.ExampleGetDataHandler
	r.Mount("/regularhandler", http.HandlerFunc(mockWithDatabaseController.ExampleGetDataHandler))

	// this is how we would use `generichandler` to write our API with Endpoint functions
	r.Mount("/put", generichandler.DefaultJSONHandlerFunc(mockWithDatabaseController.ExamplePutDataEndpoint))

	// in this example, we call a service method directly; we don't even need the `mockWithDatabaseController`
	// I think this is a bad idea because we typically separate the
	// controller and service logic and the service shouldn't be trying to tell the API how to respond. It's easier
	// to think about the service as its own module/set of functions that we can test independently of the API - this
	// way we can refactor more easily, we can cleanly add different service implementations, and it's a simpler mental
	// model. Nevertheless, it's convenient that we could potentially serve just about any SDK as a REST API without
	// writing any API code.
	r.Mount("/get", generichandler.DefaultJSONHandlerFunc(dataService.GetData))

	// this is just a dummy handler
	r.Mount("/echo", generichandler.DefaultJSONHandlerFunc(func(ctx context.Context, request *model.EchoRequest) (*model.EchoRequest, error) {
		return &model.EchoRequest{
			Message:  request.Message,
			Response: "echo!",
		}, nil
	}))

	// create HTTP server
	end := make(chan error, 1)
	go func() {
		end <- http.ListenAndServe(":8080", r)
	}()

	log.Info().Msg("ghexample started on port 8080.")

	if res := <-end; res != nil {
		log.Fatal().Err(res).Msg("there was an issue running the service")
	}
}
