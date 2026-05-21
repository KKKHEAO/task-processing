module github.com/KKKHEAO/task-processing/apps/outboxer

go 1.25.0

require (
	github.com/segmentio/kafka-go v0.4.50
	github.com/KKKHEAO/task-processing/packages/config v0.0.0
	github.com/KKKHEAO/task-processing/packages/domain v0.0.0
	github.com/KKKHEAO/task-processing/packages/postgres v0.0.0
	github.com/KKKHEAO/task-processing/packages/repository v0.0.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace (
	github.com/KKKHEAO/task-processing/packages/config => ../../packages/config
	github.com/KKKHEAO/task-processing/packages/domain => ../../packages/domain
	github.com/KKKHEAO/task-processing/packages/postgres => ../../packages/postgres
	github.com/KKKHEAO/task-processing/packages/repository => ../../packages/repository
)
