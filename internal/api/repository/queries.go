package repository

const createTaskQuery = `INSERT INTO sh_task.tasks (id, type, payload, status, retries, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)`

const getByIdQuery = `SELECT id, type, payload, status, retries, created_at, updated_at
		FROM sh_task.tasks
		WHERE id = $1`

const createOutBoxQuery = `INSERT INTO sh_task.outbox (id, topic, key, payload, created_at) VALUES ($1,$2,$3,$4,$5)`

const fetchOutboxBatch = `SELECT id, topic, key, payload, created_at
	FROM sh_task.outbox
	WHERE processed = false
	ORDER BY created_at
	LIMIT $1
	FOR UPDATE SKIP LOCKED`

const updateOutBoxQuery = `UPDATE sh_task.outbox
	SET processed = true
	WHERE id = $1`
