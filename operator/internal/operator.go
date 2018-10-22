package internal

import (
	"github.com/petomalina/krane/operator/pkg"
	"golang.org/x/net/context"
)

type OperatorServer struct {
}

func (OperatorServer) Initialize(context.Context, *operator.Job) (*operator.Job, error) {
	panic("implement me")
}

func (OperatorServer) Finish(context.Context, *operator.Job) (*operator.Job, error) {
	panic("implement me")
}
