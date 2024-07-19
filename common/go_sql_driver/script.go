package go_sql_driver

import (
	"fmt"
	"sync"
)

const (
	TableType_Primary   = 0
	TableType_Duplicate = 1
	TableType_Aggregate = 2
)

var (
	columnsSQL                 = `show columns from %s.%s;`
	createTableSQLForPrimary   = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time datetime DEFAULT CURRENT_TIMESTAMP) PRIMARY KEY (%s,_create_time) PARTITION BY date_trunc('month', _create_time) DISTRIBUTED BY HASH (%s) ORDER BY (%s) PROPERTIES ("enable_persistent_index" = "true");`
	createTableSQLForDuplicate = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time datetime DEFAULT CURRENT_TIMESTAMP) PARTITION BY date_trunc('month', _create_time) ORDER BY (%s) PROPERTIES ("enable_persistent_index" = "true");`
	createTableSQLForAggregate = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time datetime DEFAULT CURRENT_TIMESTAMP) PARTITION BY date_trunc('month', _create_time) AGGREGATE KEY(%s) DISTRIBUTED BY HASH (%s) ORDER BY (%s) PROPERTIES ("enable_persistent_index" = "true");`
	dropTableSQL               = `DROP TABLE IF EXISTS %s.%s `
	addColumnSQL               = `ALTER TABLE %s.%s ADD COLUMN %s %s`
)

func getTableColumns(database, tablename string) string {
	var query = fmt.Sprintf(columnsSQL, database, tablename)
	return query
}

func createTable(engine int, database, tablename, primaryKey, distributedKey, orderByKey string, columnsMap *sync.Map) string {
	var columns = ""
	columnsMap.Range(func(key, value any) bool {
		if key.(string) == "_create_time" {
			return true
		}
		cwt := value.(*ColumnWithType)
		if !cwt.Type.Nullable {
			columns += fmt.Sprintf("%s %s NOT NULL, ", cwt.Name, cwt.Type.ToString())
		} else {
			columns += fmt.Sprintf("%s %s, ", cwt.Name, cwt.Type.ToString())
		}
		return true
	})
	query := ""
	if engine == TableType_Primary {
		query = fmt.Sprintf(createTableSQLForPrimary, database, tablename, columns, primaryKey, distributedKey, orderByKey)
	} else if engine == TableType_Aggregate {
		query = fmt.Sprintf(createTableSQLForAggregate, database, tablename, columns, primaryKey, distributedKey, orderByKey)
	} else {
		query = fmt.Sprintf(createTableSQLForDuplicate, database, tablename, columns, orderByKey)
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
