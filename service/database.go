package service

import (
	"context"

	"github.com/Isnor/ghexample/config"
	"github.com/Isnor/ghexample/model"

	"github.com/Isnor/generichandler"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// An interface to describe reading and writing from/to an external data store
// It's important to make these functions general, so the interface is based on the models that
// the service uses rather than a datastore-specific detail that may give us more control.
// The implementations are responsible for that kind of detail, as well as the conversion
// between the models our service defines and the types defined in the SDK of the
// data store.
type DataService interface {
	PutData(context.Context, *model.ExampleData) (*model.ExampleData, error)
	GetData(context context.Context, email *model.ExampleData) (*model.ExampleData, error)
}

// An example of a postgres implementation of the above DataService
type PostgresService struct {
	config.PostgresConfig
	DB *pgx.Conn
}

// PutData adds a model.ExampleData to the configured postgres database
func (d *PostgresService) PutData(ctx context.Context, data *model.ExampleData) (*model.ExampleData, error) {
	// of course, the table should not be hardcoded - this is just an example
	err := d.DB.QueryRow(ctx, "INSERT INTO customers(id,email) VALUES($1, $2)", data.ID, data.Email).Scan()
	if err == nil || err == pgx.ErrNoRows {
		return data, nil
	}
	return nil, err
}

// GetData reads a model.ExampleData from the configured postgres database
func (d *PostgresService) GetData(ctx context.Context, data *model.ExampleData) (*model.ExampleData, error) {
	res := &model.ExampleData{}
	// of course, the table should not be hardcoded - this is just an example
	err := d.DB.QueryRow(ctx, "SELECT id,email FROM customers WHERE email=$1", data.Email).Scan(&res.ID, &res.Email)

	return res, err
}

// An example of an implementation of DataService that reads and writes from/to a dictionary in memory.
// The purpose of this is to showcase how simple it is to implement a separate data store and replace
// a previous implementation.
type InMemoryDataService struct {
	data map[string]*model.ExampleData
}

func (d *InMemoryDataService) PutData(_ context.Context, data *model.ExampleData) (*model.ExampleData, error) {
	if len(data.Email) == 0 {
		return nil, errors.WithMessage(generichandler.ErrorNotFound, "no ID provided")
	}
	d.data[data.Email] = data
	return data, nil
}

func (d *InMemoryDataService) GetData(context context.Context, email *model.ExampleData) (*model.ExampleData, error) {
	return d.data[email.Email], nil
}
