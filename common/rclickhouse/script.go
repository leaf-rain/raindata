package rclickhouse

import (
	"fmt"
	"sync"
)

var (
	columnsSQL               = `select name, type, default_kind from system.columns where database = '%s' and table = '%s'`
	createTableSQLForReplace = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time DateTime DEFAULT NOW()) ENGINE=ReplacingMergeTree(%s) PARTITION BY toYYYYMMDD(_create_time) ORDER BY (%s) SETTINGS index_granularity = 8192;`
	createTableSQL           = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time DateTime DEFAULT NOW()) ENGINE=MergeTree() PARTITION BY toYYYYMMDD(_create_time) ORDER BY (%s) SETTINGS index_granularity = 8192;`
	dropTableSQL             = `DROP TABLE IF EXISTS %s.%s `
	addColumnSQL             = `ALTER TABLE %s.%s ADD COLUMN IF NOT EXISTS %s %s`
)

func getTableColumns(database, tablename string) string {
	var query = fmt.Sprintf(columnsSQL, database, tablename)
	return query
}

func createTable(engine int, database, tablename, replayKey, orderByKey string, columnsMap *sync.Map) string {
	var columns = ""
	columnsMap.Range(func(key, value any) bool {
		if key.(string) == "_create_time" {
			return true
		}
		cwt := value.(*ColumnWithType)
		columns += fmt.Sprintf("%s %s, ", cwt.Name, cwt.Type.ToString())
		return true
	})
	query := ""
	if engine == EngineMergeTree {
		query = fmt.Sprintf(createTableSQL, database, tablename, columns, orderByKey)
	} else {
		query = fmt.Sprintf(createTableSQLForReplace, database, tablename, columns, replayKey, orderByKey)
	}
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
