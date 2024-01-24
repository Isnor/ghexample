package model

// Simple example data models; these don't do anything, they're just written to aid of describing
// the generichandler library

type EchoRequest struct {
	Message  string `json:"message"`
	Response string `json:"response,omitempty"`
}

type ExampleData struct {
	ID    string
	Email string
}

type ExampleErrorResponse struct {
	Error string `json:"error"`
}
