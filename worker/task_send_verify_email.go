package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

// PayloadSendVerifyEmail struct of data needed for worker to do tas send verify email
type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

// DistributeTaskSendVerifyEmail to Distribute new task to redis when a user was created
func (redisDistributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	//extract payload data from json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	// create new task structswith (task name, task payload, task option)
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload)
	// signed task to redis queque

	info, err := redisDistributor.client.EnqueueContext(ctx, task)
	log.Info().Err(err)
	if err != nil { //somethings happend when signed a task
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	//done signed task to redis and ready to consume by processor task
	return nil
}

// ProcessTaskSendVerifyEmail - function to process "task:send_verify_email" do sending verify email when new user are created
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	//get payload data from redis task
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		//returning error and give worker a signal to skip retring this task since next try always fail
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows { //returning error and give worker a signal to skip retring this task since next try always fail
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	//TO-DO : do send email to user
	//done processing a task
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}
