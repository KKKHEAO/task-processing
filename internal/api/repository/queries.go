package repository

const createTaskQuery = `INSERT INTO sh_task.tasks (id, type, payload, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)`

const getByIdQuery = `SELECT id, type, payload, status, created_at, updated_at
		FROM sh_task.tasks
		WHERE id = $1`

const createOutBoxQuery = `INSERT INTO sh_task.outbox (id, topic, key, payload, created_at) VALUES ($1,$2,$3,$4,$5)`

const fetchOutboxBatch = `SELECT id, topic, key, payload, created_at, retry_count, last_retry_at, next_retry_at, error_message
	FROM sh_task.outbox
	WHERE processed = false
	AND (next_retry_at IS NULL OR next_retry_at <= NOW())
	ORDER BY created_at
	LIMIT $1
	FOR UPDATE SKIP LOCKED`

const updateOutBoxQuery = `UPDATE sh_task.outbox
	SET processed = true
	WHERE id = $1`

const updateOutboxRetry = `UPDATE sh_task.outbox
	SET retry_count = $2,
		last_retry_at = $3,
		next_retry_at = $4,
		error_message = $5
	WHERE id = $1`
