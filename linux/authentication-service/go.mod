module authentication

go 1.23.3

require (
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx/v5 v5.7.1
	github.com/tsizism/semplicita/linux/shared v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.30.0
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/tsizism/semplicita/linux/shared => ./../shared
