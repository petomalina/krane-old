package internal

import (
	"github.com/petomalina/krane/operator/pkg"
	"golang.org/x/net/context"
)

type GatewayServer struct {
}

func (GatewayServer) Create(context.Context, *operator.Canary) (*operator.Canary, error) {
	panic("implement me")
}

func (GatewayServer) List(context.Context, *operator.CanaryGatewayListQuery) (*operator.CanaryList, error) {
	panic("implement me")
}

func (GatewayServer) Describe(context.Context, *operator.CanaryGatewayDescribeQuery) (*operator.Canary, error) {
	panic("implement me")
}
