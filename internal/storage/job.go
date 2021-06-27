package storage

import (
	"context"
	"database/sql"
	"fmt"
	core "github.com/calebikhuohon/hasty-test"
	"github.com/google/uuid"
	"time"
)

const sqlAddJob = `
	INSERT into job (job_id, object_id, status, sleep_time_used, max_retries, created_at, updated_at)
	VALUES (?,?,?,?,?,?,?);
`

func (s Storage) CreateJob(ctx context.Context, job core.Job) error {
	fmt.Println(job)
	_, err := s.db.ExecContext(
		ctx,
		sqlAddJob,
		job.JobID[:],
		job.ObjectID,
		job.Status,
		job.SleepTimeUsed,
		job.MaxRetries,
		job.CreatedAt,
		job.UpdatedAt,
	)

	switch {
	case isDuplicateErr(err):
		return core.ErrAlreadyExists
	default:
		return err
	}
}

const sqlGetJobById = `
	SELECT *
	FROM job
	WHERE job_id = ? 
`

func (s Storage) GetJobById(ctx context.Context, id uuid.UUID) (core.Job, error) {
	row := s.db.QueryRowContext(ctx, sqlGetJobById, id[:])
	if row == nil {
		return core.Job{}, fmt.Errorf("job not found")
	}

	return s.jobRow(row)
}

const sqlGetJobsByObjectId = `
	SELECT *
	FROM job
	WHERE object_id = ? 
`

func (s Storage) GetJobsByObjectId(ctx context.Context, objectId string) ([]core.Job, error) {
	rows, err := s.db.QueryContext(ctx, sqlGetJobsByObjectId, objectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.listOfJobRows(rows)
}

const sqlUpdateJob = `
	UPDATE job
	SET
		status = ?,
		updated_at = ?
	WHERE job_id = ?; 
`

func (s Storage) UpdateJobStatus(ctx context.Context, job core.Job, lastUpdatedAt time.Time) error {
	res, err := s.db.ExecContext(
		ctx,
		sqlUpdateJob,
		job.Status,
		time.Now().UTC(),
		job.JobID[:],
	)

	if err != nil {
		return fmt.Errorf("could not update job: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	switch {
	case err != nil:
		return fmt.Errorf("could not verify job update: could not get rowsAffected: %w", err)
	case rowsAffected == 0:
		return fmt.Errorf("job update failed: no rows were affected")
	default:
		return nil
	}
}

func (s *Storage) jobRow(row scanner) (core.Job, error) {
	var job core.Job

	err := row.Scan(
		&job.JobID,
		&job.ObjectID,
		&job.Status,
		&job.SleepTimeUsed,
		&job.MaxRetries,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	return job, err
}

func (s *Storage) listOfJobRows(rows *sql.Rows) ([]core.Job, error) {
	var jobs []core.Job

	for rows.Next() {
		job, err := s.jobRow(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}
