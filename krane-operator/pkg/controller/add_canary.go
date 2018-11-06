package controller

import (
	"github.com/petomalina/krane/operator/pkg/controller/canary"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, canary.Add)
}
