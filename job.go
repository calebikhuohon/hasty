package hasty_test

import (
	"github.com/google/uuid"
	"time"
)

type JobStatus string

const (
	JobCompleted    JobStatus = "completed"
	JobNotCompleted           = "not_completed"
	JobFailed                 = "failed"
)

type Job struct {
	JobID         uuid.UUID     `json:"job_id"`
	ObjectID      string        `json:"object_id"`
	Status        JobStatus     `json:"status"`
	SleepTimeUsed time.Duration `json:"sleep_time_used"`
	MaxRetries    int           `json:"max_retries"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
