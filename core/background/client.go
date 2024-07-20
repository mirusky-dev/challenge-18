package background

import (
	"github.com/gobp/gobp/core/env"
	"github.com/hibiken/asynq"
)

func NewClient(config env.Config) (*asynq.Client, error) {

	redisOpt, err := asynq.ParseRedisURI(config.RedisURL + "/12")
	if err != nil {
		return nil, err
	}

	return asynq.NewClient(redisOpt), nil
}
