package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mirusky-dev/challenge-18/core/background"
	"github.com/mirusky-dev/challenge-18/core/env"
	"github.com/mirusky-dev/challenge-18/router"
	"github.com/spf13/cobra"
)

var (
	apiShort = "GolangBoilerplate API server"
	apiLong  = "Run the GolangBoilerplate API server for seamless interaction with the web framework boilerplate."
)

func newCmdAPI() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "api",
		Short: apiShort,
		Long:  apiLong,
		RunE:  api,
	}

	return cmd
}

func api(cmd *cobra.Command, args []string) error {
	config, err := env.Load()
	if err != nil {
		return err
	}

	backgroundJobClient, err := background.NewClient(config)
	if err != nil {
		log.Fatal("Failed to initialize background jobs", err)
		return err
	}

	app := router.Setup(config, backgroundJobClient)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":" + config.Port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	log.Println("[api] Gracefully shutting down...")
	app.ShutdownWithTimeout(10 * time.Second)

	log.Println("[api] Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()

	log.Println("[api] Server was successful shutdown.")

	return nil
}
