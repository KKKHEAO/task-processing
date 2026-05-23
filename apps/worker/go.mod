module github.com/KKKHEAO/task-processing/apps/worker

go 1.25

require (
	github.com/KKKHEAO/task-processing/packages/config v0.0.0
	github.com/KKKHEAO/task-processing/packages/domain v0.0.0
	github.com/KKKHEAO/task-processing/packages/kafka v0.0.0
	github.com/KKKHEAO/task-processing/packages/logger v0.0.0
	github.com/segmentio/kafka-go v0.4.50
	go.uber.org/zap v1.28.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace (
	github.com/KKKHEAO/task-processing/packages/config => ../../packages/config
	github.com/KKKHEAO/task-processing/packages/domain => ../../packages/domain
	github.com/KKKHEAO/task-processing/packages/kafka => ../../packages/kafka
	github.com/KKKHEAO/task-processing/packages/logger => ../../packages/logger
)
