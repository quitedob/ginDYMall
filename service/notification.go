package service

import (
	"context"
	"encoding/json"
	"fmt" // Added for fmt.Errorf
	"time" // For timeout on BRPop

	"github.com/redis/go-redis/v9"
	"douyin/pkg/utils/email" // Your email client package
	"douyin/mylog"          // Your logger
)

const emailQueueKey = "email_queue" // Redis key for the email queue

// EmailJob defines the structure of an email task in the queue.
type EmailJob struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

// NotificationService handles asynchronous notifications.
type NotificationService struct {
	RedisClient *redis.Client
	EmailClient *email.Client
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(rdb *redis.Client, mailer *email.Client) *NotificationService {
	return &NotificationService{
		RedisClient: rdb,
		EmailClient: mailer,
	}
}

// EnqueueEmail adds an email sending task to the Redis queue.
func (ns *NotificationService) EnqueueEmail(ctx context.Context, job EmailJob) error {
	if ns.RedisClient == nil {
		mylog.Error("Redis client is nil in NotificationService. Cannot enqueue email.")
		// Optionally, try to send synchronously or return critical error
		return fmt.Errorf("redis client not initialized")
	}

	payloadBytes, err := json.Marshal(job)
	if err != nil {
		mylog.Errorf("Failed to marshal email job for %v: %v", job.To, err)
		return fmt.Errorf("failed to marshal email job: %w", err)
	}

	if err := ns.RedisClient.LPush(ctx, emailQueueKey, string(payloadBytes)).Err(); err != nil {
		mylog.Errorf("Failed to LPUSH email job to Redis for %v: %v", job.To, err)
		return fmt.Errorf("failed to enqueue email job: %w", err)
	}
	mylog.Infof("Email job enqueued for %v, subject: %s", job.To, job.Subject)
	return nil
}

// ListenAndSend continuously listens to the email queue and sends emails.
// This is intended to be run as a background goroutine.
func (ns *NotificationService) ListenAndSend(ctx context.Context) {
	if ns.RedisClient == nil || ns.EmailClient == nil {
		mylog.Error("NotificationService not properly initialized (Redis or Email client is nil). Worker stopping.")
		return
	}
	mylog.Info("Starting email queue listener...")
	for {
		select {
		case <-ctx.Done(): // Context cancelled, stop worker
			mylog.Info("Email queue listener shutting down...")
			return
		default:
			// BRPop blocks until an item is available or timeout occurs.
			// Using a timeout allows checking ctx.Done() periodically.
			result, err := ns.RedisClient.BRPop(ctx, 5*time.Second, emailQueueKey).Result()
			if err == redis.Nil {
				continue // Timeout, no item, loop again
			}
			if err != nil {
				mylog.Errorf("Error during BRPop from email queue: %v. Retrying in 5s.", err)
				time.Sleep(5 * time.Second) // Wait before retrying on error
				continue
			}

			// result is []string{queueName, value}
			if len(result) < 2 {
				mylog.Warnf("BRPop returned unexpected result: %v", result)
				continue
			}

			var job EmailJob
			if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
				mylog.Errorf("Failed to unmarshal email job from queue (%s): %v", result[1], err)
				// Consider moving to a dead-letter queue or logging extensively
				continue
			}

			mylog.Infof("Processing email job for %v, subject: %s", job.To, job.Subject)
			if err := ns.EmailClient.Send(job.To, job.Subject, job.Body); err != nil {
				// Email sending failed, log error. Consider retry mechanisms or dead-letter queue.
				mylog.Errorf("Failed to send email for job (%v, %s): %v", job.To, job.Subject, err)
			} else {
				mylog.Infof("Successfully sent email for job (%v, %s)", job.To, job.Subject)
			}
		}
	}
}
