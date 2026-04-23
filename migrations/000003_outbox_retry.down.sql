-- Удаляем индексы
DROP INDEX IF EXISTS sh_task.idx_outbox_next_retry;
DROP INDEX IF EXISTS sh_task.idx_outbox_retry_count;

-- Удаляем добавленные колонки
ALTER TABLE sh_task.outbox 
DROP COLUMN IF EXISTS retry_count,
DROP COLUMN IF EXISTS last_retry_at,
DROP COLUMN IF EXISTS next_retry_at,
DROP COLUMN IF EXISTS error_message;