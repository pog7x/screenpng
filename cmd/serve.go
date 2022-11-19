package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pog7x/screenpng/configs"
	"github.com/pog7x/screenpng/internal/logger"
	"github.com/pog7x/screenpng/internal/server"

	"github.com/pog7x/ssfactory"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launching and serving Screenshot Factory",
	Run: func(cmd *cobra.Command, args []string) {
		serve(context.Background(), configs.Configuration)
	},
}

func serve(ctx context.Context, cfg *configs.Config) {
	log, err := logger.NewLogger(cfg.Debug)
	if err != nil {
		panic(err)
	}

	f, stopFunc, err := ssfactory.NewFactory(
		ssfactory.InitFactory{
			WebdriverPort:     cfg.WebdriverPort,
			UseBrowser:        cfg.UseBrowser,
			FirefoxBinaryPath: cfg.FirefoxBinaryPath,
			GeckodriverPath:   cfg.GeckodriverPath,
			FirefoxArgs:       cfg.FirefoxArgs,
			ChromeBinaryPath:  cfg.ChromeBinaryPath,
			ChromedriverPath:  cfg.ChromedriverPath,
			ChromeArgs:        cfg.ChromeArgs,
			WorkersCount:      8,
		},
	)
	if err != nil {
		log.Error("Factory init error", zap.Error(err))
		return
	}
	defer stopFunc()

	srv := server.NewHTTPServer(log, cfg, *f)

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			log.Error("Stopping server error", zap.Error(err))
		}
	}()
	defer func() { _ = log.Sync() }()
	errChan := make(chan error)
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	go func() {
		errChan <- srv.Start()
	}()

	select {
	case err := <-errChan:
		log.Error("Run server error", zap.Error(err))
	case sig := <-sigChan:
		log.Error("Caught os signal", zap.Any("sig", sig))
	}
}
