package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-redis/pkg/client"
	cfg "github.com/conductorone/baton-redis/pkg/config"
	connectorSchema "github.com/conductorone/baton-redis/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-redis",
		getConnector,
		cfg.Config,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, rc *cfg.Redis) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	if err := cfg.ValidateConfig(rc); err != nil {
		return nil, err
	}

	redisClient := client.NewClient(rc.Username, rc.Password, rc.ClusterHost, rc.ApiPort)

	connectorBuilder, err := connectorSchema.New(ctx, redisClient)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	opts := make([]connectorbuilder.Opt, 0)

	connector, err := connectorbuilder.NewConnector(ctx, connectorBuilder, opts...)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	return connector, nil
}
