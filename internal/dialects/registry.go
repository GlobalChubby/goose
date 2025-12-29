package dialects

import (
	"github.com/pressly/goose/v3/database/dialect"
)

// Dialect is the type of database dialect.
type Dialect string

const (
	Custom     Dialect = ""
	ClickHouse Dialect = "clickhouse"
	AuroraDSQL Dialect = "dsql"
	MSSQL      Dialect = "mssql"
	MySQL      Dialect = "mysql"
	Postgres   Dialect = "postgres"
	Redshift   Dialect = "redshift"
	SQLite3    Dialect = "sqlite3"
	Spanner    Dialect = "spanner"
	Starrocks  Dialect = "starrocks"
	TiDB       Dialect = "tidb"
	Turso      Dialect = "turso"
	YdB        Dialect = "ydb"
	Vertica    Dialect = "vertica"
)

type Aliases map[string]struct{}

type Spec struct {
	// Aliases is the set of supported external names for a Dialect.
	Aliases Aliases
	// Querier returns the default [dialect.Querier] implementation for this dialect.                                              ..
	Querier func() dialect.Querier
}

func aliases(values ...string) Aliases {
	a := make(map[string]struct{}, len(values))
	for _, v := range values {
		a[v] = struct{}{}
	}
	return a
}

// Registry maps a canonical [Dialect] value to its corresponding specification
// and is the source of truth enumerating all supported Dialects.
var Registry = map[Dialect]Spec{
	Postgres: {
		aliases("postgres", "pgx"),
		// A few Querier values need to be wrapped in a type-appropriate function, because the
		// backing constructor returns a database.QuerierExtender, and its function type won't
		// match the querier type, even if database.QuerierExtender embeds database.Querier.
		func() dialect.Querier { return NewPostgres() },
	},
	MySQL: {
		aliases("mysql"),
		func() dialect.Querier { return NewMysql() },
	},
	SQLite3: {
		aliases("sqlite3", "sqlite"),
		NewSqlite3,
	},
	Spanner: {
		aliases("spanner"),
		NewSpanner,
	},
	MSSQL: {
		aliases("mssql", "azuresql", "sqlserver"),
		NewSqlserver,
	},
	Redshift: {
		aliases("redshift"),
		NewRedshift,
	},
	TiDB: {
		aliases("tidb"),
		NewTidb,
	},
	ClickHouse: {
		aliases("clickhouse"),
		NewClickhouse,
	},
	Vertica: {
		aliases("vertica"),
		NewVertica,
	},
	YdB: {
		aliases("ydb"),
		NewYDB,
	},
	Turso: {
		aliases("turso"),
		NewTurso,
	},
	Starrocks: {
		aliases("starrocks"),
		NewStarrocks,
	},
	AuroraDSQL: {
		aliases("dsql"),
		func() dialect.Querier { return NewAuroraDSQL() },
	},
}
