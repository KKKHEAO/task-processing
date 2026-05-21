module github.com/KKKHEAO/task-processing/apps/worker

go 1.25

require (
	github.com/KKKHEAO/task-processing/packages/domain v0.0.0
	github.com/KKKHEAO/task-processing/packages/kafka v0.0.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.50 // indirect
)

replace (
	github.com/KKKHEAO/task-processing/packages/domain => ../../packages/domain
	github.com/KKKHEAO/task-processing/packages/kafka => ../../packages/kafka
)
