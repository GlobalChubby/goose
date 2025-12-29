package dialects

import (
	"errors"

	"github.com/pressly/goose/v3/database/dialect"
)

// Dialect is the type of database dialect.
type Dialect string

var ErrUnknownDialect = errors.New("unknown dialect")

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

// aliases is the set of supported external names for a Dialect.
type aliases map[string]struct{}

// spec specifies a supported dialect's external aliases as well
// as its default [dialect.Querier] implementation.
type spec struct {
	// aliases is the set of supported external names for a Dialect.
	aliases aliases
	// querier is a function that returns the default [dialect.Querier] implementation for a Dialect.
	querier func() dialect.Querier
}

func alias(values ...string) aliases {
	a := make(map[string]struct{}, len(values))
	for _, v := range values {
		a[v] = struct{}{}
	}
	return a
}

// registry maps a canonical [Dialect] value to its corresponding specification
// and is the source of truth enumerating all supported Dialects.
var registry = map[Dialect]spec{
	Postgres: {
		alias("postgres", "pgx"),
		// A few querier values need to be wrapped in a type-appropriate function, because the
		// backing constructor returns a database.QuerierExtender, and its function type won't
		// match the querier type, even if database.QuerierExtender embeds database.Querier.
		func() dialect.Querier { return NewPostgres() },
	},
	MySQL: {
		alias("mysql"),
		func() dialect.Querier { return NewMysql() },
	},
	SQLite3: {
		alias("sqlite3", "sqlite"),
		NewSqlite3,
	},
	Spanner: {
		alias("spanner"),
		NewSpanner,
	},
	MSSQL: {
		alias("mssql", "azuresql", "sqlserver"),
		NewSqlserver,
	},
	Redshift: {
		alias("redshift"),
		NewRedshift,
	},
	TiDB: {
		alias("tidb"),
		NewTidb,
	},
	ClickHouse: {
		alias("clickhouse"),
		NewClickhouse,
	},
	Vertica: {
		alias("vertica"),
		NewVertica,
	},
	YdB: {
		alias("ydb"),
		NewYDB,
	},
	Turso: {
		alias("turso"),
		NewTurso,
	},
	Starrocks: {
		alias("starrocks"),
		NewStarrocks,
	},
	AuroraDSQL: {
		alias("dsql"),
		func() dialect.Querier { return NewAuroraDSQL() },
	},
}

// Parse returns the corresponding [Dialect], if any, for a given [Dialect] alias.
func Parse(s string) (Dialect, error) {
	for d, spec := range registry {
		_, ok := spec.aliases[s]
		if ok {
			return d, nil
		}
	}
	return "", ErrUnknownDialect
}

// Querier returns the default [dialect.Querier] implementation for the dialect.
func Querier(d Dialect) (dialect.Querier, error) {
	s, ok := registry[d]
	if ok {
		return s.querier(), nil
	}
	return nil, ErrUnknownDialect
}
