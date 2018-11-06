package main

import (
	"context"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"
	"github.com/petomalina/krane/operator/pkg"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"os"
	"sort"
)

var (
	client operator.CanaryGatewayClient
	bucket *blob.Bucket
)

func main() {
	app := cli.NewApp()
	app.Name = "Krane"
	app.Usage = "manipulates the krane-krane-operator-old"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host, h",
			Value: "127.0.0.1:22334",
			Usage: "Host to connect to",
		},
		cli.StringFlag{
			Name:  "bucket, b",
			Value: "krane-artifacts",
			Usage: "Name of the bucket to store artifacts in",
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		client, err = initClient(c.String("host"))
		if err != nil {
			return err
		}

		bucket, err = setupGCP(context.Background(), c.String("bucket"))

		return err
	}

	app.Commands = []cli.Command{
		{
			// @Deprecated
			Name:    "create",
			Aliases: []string{"c"},
			Action: func(c *cli.Context) error {
				client.Create(context.Background(), nil)

				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Action: func(c *cli.Context) error {
				client.List(context.Background(), nil)

				return nil
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) error {
				client.Describe(context.Background(), nil)

				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func initClient(host string) (operator.CanaryGatewayClient, error) {
	gwConn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return operator.NewCanaryGatewayClient(gwConn), nil
}

func setupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}
	c, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}

	return gcsblob.OpenBucket(ctx, bucket, c)
}
