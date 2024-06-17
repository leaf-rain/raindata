package clickhouse_sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/RoaringBitmap/roaring"
	"github.com/leaf-rain/raindata/common/ecode"
	"go.uber.org/zap"
	"log"
	"time"
)

// ClickhouseConfig configuration parameters
type ClickhouseConfig struct {
	Cluster  string     `json:"cluster,omitempty" yaml:"cluster"`
	DB       string     `json:"DB,omitempty" yaml:"DB"`
	Hosts    [][]string `json:"hosts,omitempty" yaml:"hosts"`
	Username string     `json:"username,omitempty" yaml:"username"`
	Password string     `json:"password,omitempty" yaml:"password"`
	Protocol string     `json:"protocol,omitempty" yaml:"protocol"` //native, http

	// Whether enable TLS encryption with clickhouse-server
	Secure bool `json:"secure,omitempty" yaml:"secure"`
	// Whether skip verify clickhouse-server cert
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify"`

	RetryTimes   int `json:"retryTimes,omitempty" yaml:"retryTimes"` // <=0 means retry infinitely
	MaxOpenConns int `json:"maxOpenConns,omitempty" yaml:"maxOpenConns"`
}

func NewClickhouse(cfg *ClickhouseConfig) (*Conn, error) {
	var err error
	if len(cfg.Hosts) == 0 {
		err = errors.New("invalid configuration, Clickhouse section is missing")
		return nil, err
	}
	if cfg.RetryTimes < 0 {
		cfg.RetryTimes = 0
	}
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = defaultMaxOpenConns
	}

	if cfg.Protocol == "" {
		cfg.Protocol = clickhouse.Native.String()
	}
	if cfg.Cluster == "" {
		var numHosts int
		for _, shard := range cfg.Hosts {
			numHosts += len(shard)
		}
		if numHosts > 1 {
			err = ecode.Newf("Need to set cluster name when DynamicSchema is enabled and number of hosts is more than one")
			return nil, err
		}
	}
	proto := clickhouse.Native
	if cfg.Protocol == clickhouse.HTTP.String() {
		proto = clickhouse.HTTP
	}
	conn := Conn{
		protocol: proto,
		ctx:      context.Background(),
	}
	var opts = new(clickhouse.Options)
	for i := range cfg.Hosts {
		opts.Addr = append(opts.Addr, cfg.Hosts[i]...)
	}
	if cfg.Username != "" && cfg.Password != "" && cfg.DB != "" {
		opts.Auth = clickhouse.Auth{
			Database: cfg.DB,
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}
	if proto == clickhouse.HTTP {
		// refers to https://github.com/ClickHouse/clickhouse-go/issues/1150
		// An obscure error in the HTTP protocol when using compression
		// disable compression in the HTTP protocol
		conn.db = clickhouse.OpenDB(opts)
		conn.db.SetMaxOpenConns(cfg.MaxOpenConns)
		conn.db.SetMaxIdleConns(cfg.MaxOpenConns)
		conn.db.SetConnMaxLifetime(time.Minute * 10)
	} else {
		opts.Compression = &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		}
		conn.c, err = clickhouse.Open(opts)
		if err != nil {
			return nil, err
		}
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

type Conn struct {
	protocol clickhouse.Protocol
	c        driver.Conn
	db       *sql.DB
	ctx      context.Context
}

func (c *Conn) Query(query string, args ...any) (*Rows, error) {
	var rs Rows
	rs.protocol = c.protocol
	if c.protocol == clickhouse.HTTP {
		rows, err := c.db.Query(query, args...)
		if err != nil {
			return &rs, err
		} else {
			rs.rs1 = rows
		}
	} else {
		rows, err := c.c.Query(c.ctx, query, args...)
		if err != nil {
			return &rs, err
		} else {
			rs.rs2 = rows
		}
	}
	return &rs, nil
}

func (c *Conn) QueryRow(query string, args ...any) *Row {
	var row Row
	row.proto = c.protocol
	if c.protocol == clickhouse.HTTP {
		row.r1 = c.db.QueryRow(query, args...)
	} else {
		row.r2 = c.c.QueryRow(c.ctx, query, args...)
	}
	return &row
}

func (c *Conn) Exec(query string, args ...any) error {
	if c.protocol == clickhouse.HTTP {
		_, err := c.db.Exec(query, args...)
		return err
	} else {
		return c.c.Exec(c.ctx, query, args...)
	}
}

func (c *Conn) Ping() error {
	if c.protocol == clickhouse.HTTP {
		return c.db.Ping()
	} else {
		return c.c.Ping(c.ctx)
	}
}

func (c *Conn) write_v1(prepareSQL string, rows [][]interface{}, idxBegin, idxEnd int) (numBad int, err error) {
	var errExec error

	var stmt *sql.Stmt
	var tx *sql.Tx
	tx, err = c.db.Begin()
	if err != nil {
		err = ecode.Wrapf(err, "pool.Conn.Begin")
		return
	}

	if stmt, err = tx.Prepare(prepareSQL); err != nil {
		err = ecode.Wrapf(err, "tx.Prepare %s", prepareSQL)
		return
	}
	defer stmt.Close()

	var bmBad *roaring.Bitmap
	for i, row := range rows {
		if _, err = stmt.Exec((row)[idxBegin:idxEnd]...); err != nil {
			if bmBad == nil {
				errExec = ecode.Wrapf(err, "driver.Batch.Append")
				bmBad = roaring.NewBitmap()
			}
			bmBad.AddInt(i)
		}

	}
	if errExec != nil {
		_ = tx.Rollback()
		numBad = int(bmBad.GetCardinality())
		// write rows again, skip bad ones
		if stmt, err = tx.Prepare(prepareSQL); err != nil {
			err = ecode.Wrapf(err, "tx.Prepare %s", prepareSQL)
			return
		}
		for i, row := range rows {
			if !bmBad.ContainsInt(i) {
				if _, err = stmt.Exec((row)[idxBegin:idxEnd]...); err != nil {
					break
				}
			}
		}
		if err = tx.Commit(); err != nil {
			err = ecode.Wrapf(err, "tx.Commit")
			_ = tx.Rollback()
			return
		}
		return
	}
	if err = tx.Commit(); err != nil {
		err = ecode.Wrapf(err, "tx.Commit")
		_ = tx.Rollback()
		return
	}
	return
}

func (c *Conn) write_v2(prepareSQL string, rows [][]interface{}, idxBegin, idxEnd int) (numBad int, err error) {
	var errExec error
	var batch driver.Batch
	if batch, err = c.c.PrepareBatch(c.ctx, prepareSQL); err != nil {
		return
	}
	var bmBad *roaring.Bitmap
	for i, row := range rows {
		if err = batch.Append((row)[idxBegin:idxEnd]...); err != nil {
			if bmBad == nil {
				errExec = errors.New(err.Error() + " driver.Batch.Append")
				bmBad = roaring.NewBitmap()
			}
			bmBad.AddInt(i)
		}
	}
	if errExec != nil {
		_ = batch.Abort()
		numBad = int(bmBad.GetCardinality())
		log.Printf(fmt.Sprintf("writeRows skipped %d rows of %d due to invalid content", numBad, len(rows)), zap.Error(errExec))
		// write rows again, skip bad ones
		if batch, err = c.c.PrepareBatch(c.ctx, prepareSQL); err != nil {
			err = Wrapf(err, "pool.Conn.PrepareBatch %s", prepareSQL)
			return
		}
		for i, row := range rows {
			if !bmBad.ContainsInt(i) {
				if err = batch.Append((row)[idxBegin:idxEnd]...); err != nil {
					break
				}
			}
		}
		if err = batch.Send(); err != nil {
			err = ecode.Wrapf(err, "driver.Batch.Send")
			_ = batch.Abort()
			return
		}
		return
	}
	if err = batch.Send(); err != nil {
		err = ecode.Wrapf(err, "driver.Batch.Send")
		_ = batch.Abort()
		return
	}
	return
}

func (c *Conn) Write(prepareSQL string, rows [][]interface{}, idxBegin, idxEnd int) (numBad int, err error) {
	log.Printf("start write to ck, begin:%d, idEnd:%d", idxBegin, idxEnd)
	if c.protocol == clickhouse.HTTP {
		numBad, err = c.write_v1(prepareSQL, rows, idxBegin, idxEnd)
	} else {
		numBad, err = c.write_v2(prepareSQL, rows, idxBegin, idxEnd)
	}
	log.Printf("loop write completed, numbad:%d", numBad)
	return numBad, err
}

func (c *Conn) Close() error {
	if c.protocol == clickhouse.HTTP {
		return c.db.Close()
	} else {
		return c.c.Close()
	}
}
