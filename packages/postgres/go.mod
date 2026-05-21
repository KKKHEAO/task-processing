module github.com/KKKHEAO/task-processing/packages/postgres

go 1.25.0

require (
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/KKKHEAO/task-processing/packages/config v0.0.0
)

require (
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace github.com/KKKHEAO/task-processing/packages/config => ../config
