package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgPath = flag.String("f", "configs/config.yaml", "path to the go's config")
	useGrpc = flag.Bool("g", true, "use grpc server?")
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	queries, conn, err := repository.ConfigurePostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := httpserver.New(cfg, conn, queries)

	var grpcServer *grpc.Server
	if *useGrpc {
		grpcServer, err := pvzv1.New(cfg.GRPCServerCfg, &app.Service.PvzService)
		if err != nil {
			log.Fatalf("failed to create grpc server: %v", err)
		}
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}

	go func() {
		<-ctx.Done()

		app.Stop(ctx)
		if grpcServer != nil {
			grpcServer.Stop()
		}
	}()

	metrics.StartMetricsServer()

	if err := app.Start(ctx); err != nil {
		log.Printf("http server returned error: %v", err)
	}
}
