-- Добавляем колонки для retry в существующую таблицу outbox
ALTER TABLE sh_task.outbox 
ADD COLUMN retry_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN last_retry_at TIMESTAMP,
ADD COLUMN next_retry_at TIMESTAMP,
ADD COLUMN error_message TEXT;

-- Создаем индекс для эффективного поиска сообщений, готовых к retry
CREATE INDEX idx_outbox_next_retry ON sh_task.outbox(next_retry_at) 
WHERE processed = FALSE AND retry_count > 0;

-- Создаем индекс для поиска сообщений по retry count
CREATE INDEX idx_outbox_retry_count ON sh_task.outbox(retry_count) 
WHERE processed = FALSE;