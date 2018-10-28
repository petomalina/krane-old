package internal

import (
	"context"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/apps/v1"
)

const (
	// CanaryAnnotation results in registering Deployment into the
	// Canary Deployments database of Crane. If the annotation
	// is set to true, the deployment will automatically begin testing,
	// otherwise, it will only be registered and will be waiting until
	// approval
	CanaryAnnotation = "canary.krane.io/run"
)

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1.Deployment:
		if _, ok := o.Annotations[""]; ok {

		}
	}
	return nil
}
