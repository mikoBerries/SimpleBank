package worker

import (
	"context"

	db "github.com/MikoBerries/SimpleBank/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

// queque name
const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

// TaskProcessor write your new task processor for each task func
type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

// NewRedisTaskProcessor returning redis server processor
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	//create Task processor with cotum option
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			//queues of priority higher run first
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			//overide eror handle func with our
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			//set costume logger
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

// Start - to set some woker and start redis server
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	//list your task here (task name , task worker func() )
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
