package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	//DistributeTaskSendVerifyEmail to Distribute task to redis when new user was created
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

// NewRedisTaskDistributor create redis task distributor with given opt
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
