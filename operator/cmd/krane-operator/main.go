package main

import (
	"context"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/petomalina/krane/operator/internal"
)

const (
	api            = ""
	deploymentKind = ""
	namespace      = ""
	resyncPeriod   = 0
)

func main() {
	handler := &internal.Handler{}

	sdk.Watch(api, deploymentKind, namespace, resyncPeriod)
	sdk.Handle(handler)
	sdk.Run(context.Background())
}
