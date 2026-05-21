module github.com/KKKHEAO/task-processing/apps/api

go 1.25.0

require (
	github.com/google/uuid v1.6.0
	github.com/KKKHEAO/task-processing/packages/config v0.0.0
	github.com/KKKHEAO/task-processing/packages/domain v0.0.0
	github.com/KKKHEAO/task-processing/packages/postgres v0.0.0
	github.com/KKKHEAO/task-processing/packages/repository v0.0.0
	google.golang.org/grpc v1.79.1
	task-processing/proto v0.0.0
)

require (
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace (
	github.com/KKKHEAO/task-processing/packages/config => ../../packages/config
	github.com/KKKHEAO/task-processing/packages/domain => ../../packages/domain
	github.com/KKKHEAO/task-processing/packages/postgres => ../../packages/postgres
	github.com/KKKHEAO/task-processing/packages/repository => ../../packages/repository
	task-processing/proto => ../../proto
)
