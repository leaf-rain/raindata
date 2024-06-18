package clickhouse_sqlx

type Fetch struct {
	TraceId string
	Data    []string
	Sql     []string
}
