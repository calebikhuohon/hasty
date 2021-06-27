package service

import (
	"context"
	"fmt"
	core "github.com/calebikhuohon/hasty-test"
	"github.com/calebikhuohon/hasty-test/internal/pkg/cron"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"time"
)

type JobService struct {
	rdb *redis.Client
	storage    JobStorage
	MaxRetries int
}

type JobStorage interface {
	CreateJob(ctx context.Context, job core.Job) error
	GetJobById(ctx context.Context, id uuid.UUID) (core.Job, error)
	GetJobsByObjectId(ctx context.Context, objectId string) ([]core.Job, error)
	UpdateJobStatus(ctx context.Context, job core.Job, lastUpdatedAt time.Time) error
}

func NewJobService(storage JobStorage, maxRetries int, rdb *redis.Client) *JobService {
	return &JobService{
		storage:    storage,
		MaxRetries: maxRetries,
		rdb: rdb,
	}
}

func (j JobService) RunJob(ctx context.Context, job core.Job) (uuid.UUID, error) {
	timeout := time.Unix(time.Now().Add(time.Minute * 5).Unix(), 0)

	val:= j.rdb.Exists(ctx, job.ObjectID).Val()
	if val != 0 {
		return uuid.Nil, fmt.Errorf("duplicate task")
	}

	job.JobID = uuid.New()
	job.Status = core.JobNotCompleted
	job.SleepTimeUsed = 34 * time.Second
	job.CreatedAt = time.Now().UTC()
	job.UpdatedAt = job.CreatedAt
	job.MaxRetries = j.MaxRetries

	//store job in redis cache
	err := j.rdb.Set(ctx, job.ObjectID, job.JobID, timeout.Sub(time.Now())).Err()
	if err != nil {
		return uuid.Nil, err
	}

	if err := j.storage.CreateJob(ctx, job); err != nil {
		log.Println("create job err: ", err)
		return uuid.Nil, err
	}

	jobChan := make(chan core.Job )

	go worker(ctx, jobChan, j.storage, j.rdb)
	jobChan <- job

	//re-run failing jobs
	j.RerunFailingJobs(ctx, job, timeout)

	return job.JobID, nil
}

func (j JobService) GetJobDetails(ctx context.Context, jobId string) (core.Job, error) {
	parsedJobId, err := uuid.Parse(jobId)
	if err != nil {
		return core.Job{}, err
	}

	return j.storage.GetJobById(ctx, parsedJobId)
}

func (j JobService) RerunFailingJobs(ctx context.Context, job core.Job, timeout time.Time)  {
	jt := cron.NewJobTicker()

	for {
		<-jt.T.C
		jobs, err := j.storage.GetJobsByObjectId(ctx, job.ObjectID)
		if err != nil {
			return
		}
		for _, jb := range jobs {
			log.Println("running a failed job")
			if jb.Status != core.JobCompleted {
				err := j.rdb.Set(ctx, job.ObjectID, job.JobID, timeout.Sub(time.Now())).Err()
				if err != nil {
					return
				}

				jobChan := make(chan core.Job )

				go worker(ctx, jobChan, j.storage, j.rdb)
				jobChan <- job
			}
		}
		jt.UpdateJobTicker()
	}
}

func worker(ctx context.Context, jobChan <-chan core.Job, storage JobStorage, rdb *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-jobChan:
			time.Sleep(job.SleepTimeUsed)
			updatedJob := job
			updatedJob.UpdatedAt = time.Now().UTC()
			updatedJob.Status = core.JobCompleted

			_, err := rdb.Del(ctx, job.ObjectID).Result()
			if err != nil {
				log.Print(err)
			}

			if err := storage.UpdateJobStatus(ctx, updatedJob, job.UpdatedAt); err != nil {
				log.Print(err)
			}
		}

	}
}
