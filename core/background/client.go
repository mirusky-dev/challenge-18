package background

import (
	"github.com/hibiken/asynq"
	"github.com/mirusky-dev/challenge-18/core/env"
)

func NewClient(config env.Config) (*asynq.Client, error) {

	redisOpt, err := asynq.ParseRedisURI(config.RedisURL + "/12")
	if err != nil {
		return nil, err
	}

	return asynq.NewClient(redisOpt), nil
}
