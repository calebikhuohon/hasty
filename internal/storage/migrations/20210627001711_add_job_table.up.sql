CREATE TABLE job
(
	job_id          BINARY(16)  NOT NULL PRIMARY KEY,
	object_id       VARCHAR(32) NOT NULL,
	status          VARCHAR(32) NOT NULL,
	sleep_time_used VARCHAR(32)  NOT NULL,
	max_retries     INT         NOT NULL,
	created_at      DATETIME    NOT NULL,
	updated_at      DATETIME    NOT NULL,

	UNIQUE KEY (job_id)
);