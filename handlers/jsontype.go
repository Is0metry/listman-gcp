package handlers

import (
	"github.com/Is0metry/listman-gcp/lists"
)

type OperationRequest struct {
	ReqType string
	Text    string
}

type ListResponse struct {
	GetSuccessful bool
	ErrorText     string
	Lst           *lists.List
}
