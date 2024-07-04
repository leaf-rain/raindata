package clickhouse_sqlx

import (
	"fmt"
	"sync"
)

var (
	columnsSQL     = `select name, type, default_kind from system.columns where database = '%s' and table = '%s'`
	createTableSQL = `CREATE TABLE IF NOT EXISTS %s.%s (%s create_time DateTime) ENGINE=ReplacingMergeTree(%s) PARTITION BY toYYYYMMDD(create_time) ORDER BY (%s) SETTINGS index_granularity = 8192;`
	dropTableSQL   = `DROP TABLE IF EXISTS %s.%s `
	addColumnSQL   = `ALTER TABLE %s.%s ADD COLUMN IF NOT EXISTS %s %s`
)

func getTableColumns(database, tablename string) string {
	var query = fmt.Sprintf(columnsSQL, database, tablename)
	return query
}

func createTable(database, tablename, replayKey, orderByKey string, columnsMap *sync.Map) string {
	var columns = ""
	columnsMap.Range(func(key, value any) bool {
		if key.(string) == "create_time" {
			return true
		}
		cwt := value.(*ColumnWithType)
		columns += fmt.Sprintf("%s %s, ", cwt.Name, cwt.Type.ToString())
		return true
	})
	var query = fmt.Sprintf(createTableSQL, database, tablename, columns, replayKey, orderByKey)
	return query
}

func delTable(database, tablename string) string {
	var query = fmt.Sprintf(dropTableSQL, database, tablename)
	return query
}

func addTableColumns(database, tablename string, columnName, columnType string) string {
	var query = fmt.Sprintf(addColumnSQL, database, tablename, columnName, columnType)
	return query
}
