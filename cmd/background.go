package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mirusky-dev/challenge-18/core/background"
	"github.com/mirusky-dev/challenge-18/core/background/events"
	"github.com/mirusky-dev/challenge-18/core/background/handlers"
	"github.com/mirusky-dev/challenge-18/core/env"
	"github.com/mirusky-dev/challenge-18/core/mailer"
	"github.com/mirusky-dev/challenge-18/repositories"
)

func newCmdBackground() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "background",
		Short: "Runs Background Jobs",
		Long:  "TODO", // TODO: Define a long run description
		RunE:  backgroundJobs,
	}

	return cmd
}

func setupBackground(config env.Config, backgroundClient *asynq.Client) (*asynq.Server, *asynq.ServeMux, error) {
	svr, mux, err := background.NewServerMux(config)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	mailer, _ := mailer.NewNoopMailer(config)

	userRepository := repositories.NewUserRepository(db)
	taskRepository := repositories.NewTaskRepository(db)

	ctrl := handlers.Controller{
		Config:           config,
		Mailer:           mailer,
		BackgroundClient: backgroundClient,

		UserRepository: userRepository,
		TaskRepository: taskRepository,
	}

	mux.HandleFunc(events.TypeTaskCompleted, ctrl.HandleTaskCompleted)

	return svr, mux, err
}

func backgroundJobs(cmd *cobra.Command, args []string) error {

	config, err := env.Load()
	if err != nil {
		log.Fatal("Failed to load env vars", err)
		return err
	}

	backgroundJobClient, err := background.NewClient(config)
	if err != nil {
		return err
	}

	svr, mux, err := setupBackground(config, backgroundJobClient)
	if err != nil {
		return err
	}

	// Listen from a different goroutine
	go func() {
		if err := svr.Start(mux); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	log.Println("[background] Gracefully shutting down...")
	svr.Shutdown()

	log.Println("[background] Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()

	backgroundJobClient.Close()

	log.Println("[background] Server was successful shutdown.")

	return nil
}
