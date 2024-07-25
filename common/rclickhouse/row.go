package rclickhouse

import (
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Row struct {
	proto clickhouse.Protocol
	r1    *sql.Row
	r2    driver.Row
}

func (r *Row) Scan(dest ...any) error {
	if r.proto == clickhouse.HTTP {
		return r.r1.Scan(dest...)
	} else {
		return r.r2.Scan(dest...)
	}
}

type Rows struct {
	protocol clickhouse.Protocol
	rs1      *sql.Rows
	rs2      driver.Rows
}

func (r *Rows) Close() error {
	if r.protocol == clickhouse.HTTP {
		return r.rs1.Close()
	} else {
		return r.rs2.Close()
	}
}

func (r *Rows) Columns() ([]string, error) {
	if r.protocol == clickhouse.HTTP {
		return r.rs1.Columns()
	} else {
		return r.rs2.Columns(), nil
	}
}

func (r *Rows) Next() bool {
	if r.protocol == clickhouse.HTTP {
		return r.rs1.Next()
	} else {
		return r.rs2.Next()
	}
}

func (r *Rows) Scan(dest ...any) error {
	if r.protocol == clickhouse.HTTP {
		return r.rs1.Scan(dest...)
	} else {
		return r.rs2.Scan(dest...)
	}
}
