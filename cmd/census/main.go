package main

import (
	"github.com/GeneralKenobi/census/internal/api/httpgin"
	"github.com/GeneralKenobi/census/internal/config"
	"github.com/GeneralKenobi/census/internal/db/postgres"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/shutdown"
	"github.com/GeneralKenobi/census/pkg/util"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configure()
	parentCtx := shutdown.NewParentContext(time.Duration(config.Get().Global.ShutdownTimeoutSeconds) * time.Second)
	bootstrap(parentCtx)
	shutdownAfterStopSignal(parentCtx)
}

func configure() {
	argsCfg := commandLineArgsConfig()

	err := mdctx.SetLogLevelFromString(argsCfg.logLevel)
	if err != nil {
		mdctx.Fatalf(nil, "Error setting log level: %v", err)
	}

	err = config.Load(argsCfg.configFiles)
	if err != nil {
		mdctx.Fatalf(nil, "Error loading configuration: %v", err)
	}

	seed, err := util.RngSeed()
	if err != nil {
		mdctx.Fatalf(nil, "Error seeding random number generator: %v", err)
	}
	rand.Seed(seed)
}

func bootstrap(parentCtx shutdown.ParentContext) {
	// DB
	dbCtx, err := postgres.NewContext(parentCtx.NewContext("postgres"))
	if err != nil {
		mdctx.Fatalf(nil, "Error connecting to DB: %v", err)
	}

	// HTTP server
	httpServer := httpgin.NewServer(dbCtx)
	go httpServer.Run(parentCtx.NewContext("http server"))
}

func shutdownAfterStopSignal(parentCtx shutdown.ParentContext) {
	stopSignalChannel := make(chan os.Signal)
	// SIGINT for ctrl+c, SIGTERM for k8s stopping the container.
	signal.Notify(stopSignalChannel, syscall.SIGINT, syscall.SIGTERM)

	caughtSignal := <-stopSignalChannel
	mdctx.Infof(nil, "Caught signal %v, shutting down", caughtSignal)

	parentCtx.Cancel()
	mdctx.Infof(nil, "Shutdown completed, exiting")
}
