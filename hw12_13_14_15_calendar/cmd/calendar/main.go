package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

const defaultConfigPath = "configs/config.toml"

func init() {
	flag.StringVar(&configFile, "config", defaultConfigPath, "Path to configuration file")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return nil
	}

	config, err := NewConfig(configFile)
	if err != nil {
		return err
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		return err
	}

	var stor storage.Storage
	switch config.Storage.Type {
	case "postgres":
		stor = sqlstorage.New()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := stor.Connect(ctx, config.Storage.ConnectionString)
		if err != nil {
			logg.Error(err)

			os.Exit(1)
		}
		defer func() {
			cancel()
			stor.Close(context.Background())
		}()
	case "inmemory":
		stor = memorystorage.New()
	default:
		stor = memorystorage.New()
	}

	calendar := app.New(logg, stor)

	addr := net.JoinHostPort("localhost", config.Server.Port)
	server := internalhttp.NewServer(logg, calendar, addr)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	return nil
}
