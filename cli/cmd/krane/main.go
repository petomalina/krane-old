package main

import (
	"context"
	"github.com/petomalina/krane/operator/pkg"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"os"
	"sort"
)

func main() {
	app := cli.NewApp()
	app.Name = "Krane"
	app.Usage = "manipulates the krane-operator"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "host, h",
			Value: "127.0.0.1:22334",
			Usage: "Host to connect to",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "create",
			Aliases: []string{"c"},
			Action: func(c *cli.Context) error {
				client, err := initClient(c.String("host"))
				if err != nil {
					return err
				}

				client.Create(context.Background(), nil)

				return nil
			},
		},
		{
			Name: "list",
			Aliases: []string{"l"},
			Action: func(c *cli.Context) error {
				client, err := initClient(c.String("host"))
				if err != nil {
					return err
				}

				client.List(context.Background(), nil)

				return nil
			},
		},
		{
			Name: "describe",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) error {
				client, err := initClient(c.String("host"))
				if err != nil {
					return err
				}

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