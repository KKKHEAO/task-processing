module github.com/KKKHEAO/task-processing/packages/logger

go 1.25

require (
	github.com/KKKHEAO/task-processing/packages/config v0.0.0
	go.uber.org/zap v1.28.0
)

require go.uber.org/multierr v1.11.0 // indirect

replace github.com/KKKHEAO/task-processing/packages/config => ../config
