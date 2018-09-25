package handlers

import (
	"github.com/Is0metry/listman-gcp/lists"
)

//OperationRequest is a struct used as the format for incoming requests. If the request is an add, Text will contain the new Item's text.
//If it is a delete, it will contain an encoded Datastore key.
type OperationRequest struct {
	Text string
}

//ListResponse is a struct used for returning a list. If the operation was a success, GetSuccessful will be true, and Lst will contain the requested list.
//Otherwise, GetSuccessful will be false and ErrorText will contain the error's text.
type ListResponse struct {
	GetSuccessful bool
	ErrorText     string
	Lst           *lists.List
}

//OperationResponse is used for returning information on an add or delete call.
type OperationResponse struct {
	Success   bool
	ErrorText string
}
