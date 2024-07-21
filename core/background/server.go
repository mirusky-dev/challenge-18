package background

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mirusky-dev/challenge-18/core/env"
)

type wrapLogger struct {
	*log.Logger
}

func (wrap wrapLogger) Debug(args ...interface{}) {
	wrap.Logger.Print(args...)
}

func (wrap wrapLogger) Info(args ...interface{}) {
	wrap.Logger.Print(args...)
}

func (wrap wrapLogger) Warn(args ...interface{}) {
	wrap.Logger.Print(args...)
}

func (wrap wrapLogger) Error(args ...interface{}) {
	wrap.Logger.Print(args...)
}

func NewWrapLogger() asynq.Logger {
	return wrapLogger{
		log.Default(),
	}
}

func NewServerMux(config env.Config) (*asynq.Server, *asynq.ServeMux, error) {

	redisOpt, err := asynq.ParseRedisURI(config.RedisURL + "/12")
	if err != nil {
		return nil, nil, err
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Logger: NewWrapLogger(),
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"highest": 3,
				"high":    3,
				"default": 2,
				"low":     1,
				"lowest":  1,
			},
			// See the godoc for other configuration options
			// LogLevel:        asynq.FatalLevel,
			ShutdownTimeout: 5 * time.Second,
		},
	)

	mux := asynq.NewServeMux()

	return srv, mux, nil
}
