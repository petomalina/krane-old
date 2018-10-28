package infra

import (
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
)

func ReconcileCluster() {
	hc, err := google.DefaultClient(context.Background(), container.CloudPlatformScope)
	if err != nil {
		panic(err)
	}

	c, err := container.New(hc)
	if err != nil {
		panic(err)
	}

	req := &container.CreateClusterRequest{}
	op, err := c.Projects.Zones.Clusters.Create("", "", req).Do()
	if err != nil {
		panic(err)
	}

	c.Projects.Zones.Operations.Get("", "", op.Name)

}
