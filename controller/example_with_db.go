package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Isnor/ghexample/model"
	"github.com/Isnor/ghexample/service"

	"github.com/rs/zerolog/log"
)

// This is a simple HTTP controller with two endpoints to read from and write to an external data store.
// I like to organize controllers per-path; e.g. /users and /orgs would be two files, each with their own
// struct that defines all of the functionality for those routes - a UserController and an OrgController.
type ExampleControllerWithDatabase struct {
	DataService service.DataService
}

// This is an example of how typical functions are currently written http.Handlers - we need to use the
// http package to read and write, somewhat awkwardly - and we always need to write the boilerplate to
// decode and encode JSON because we want to use the typed object and Golang's Marshal/Unmarshal functionality.
func (m *ExampleControllerWithDatabase) ExampleGetDataHandler(w http.ResponseWriter, r *http.Request) {

	req := &model.ExampleData{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Error().Msg("no data provided")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ExampleErrorResponse{Error: "no email provided"})
		return
	}

	data, err := m.DataService.GetData(r.Context(), req)
	if err != nil {
		log.Err(err).Msg("Failed to get data from database")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ExampleErrorResponse{Error: "failed getting from database"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// This is an example of what the generichandler library calls an "APIEndpoint". The signature uses a context and some
// JSON object, and returns a JSON object and an error. It looks like a regular function that we'd expect to handle
// as a Golang programmer, it's extremely easy to test and reason about, and it doesn't have any HTTP details or boilerplate.
// All of that is handled by the generichandler library
func (m *ExampleControllerWithDatabase) ExamplePutDataEndpoint(ctx context.Context, data *model.ExampleData) (*model.ExampleData, error) {
	d, err := m.DataService.PutData(ctx, data)
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}
	return d, nil
}
