package rsql

import (
	"fmt"
	"strings"
	"sync"
)

const (
	TableType_Primary   = 0
	TableType_Duplicate = 1
	TableType_Aggregate = 2
)

var (
	columnsSQL                 = `show columns from %s.%s;`
	createTableSQLForPrimary   = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP) PRIMARY KEY (%s,_create_time) PARTITION BY date_trunc('month', _create_time) DISTRIBUTED BY HASH (%s) ORDER BY (%s) PROPERTIES ("enable_persistent_index" = "true");`
	createTableSQLForDuplicate = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP) DUPLICATE KEY (%s) PARTITION BY date_trunc('month', _create_time) ORDER BY (%s) PROPERTIES ("enable_persistent_index" = "true");`
	createTableSQLForAggregate = `CREATE TABLE IF NOT EXISTS %s.%s (%s _create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP) PARTITION BY date_trunc('month', _create_time) AGGREGATE KEY(%s) DISTRIBUTED BY HASH (%s) ORDER BY (%s) PROPERTIES ("enable_persistent_index" = "true");`
	dropTableSQL               = `DROP TABLE IF EXISTS %s.%s `
	addColumnSQL               = `ALTER TABLE %s.%s ADD COLUMN %s %s`
)

func getTableColumns(database, tablename string) string {
	var query = fmt.Sprintf(columnsSQL, database, tablename)
	return query
}

func createTable(engine int, database, tablename, primaryKey, distributedKey, orderByKey string, columnsMap *sync.Map) string {
	var columns = ""
	var orderBySlice = strings.Split(orderByKey, ",")
	for i := 0; i < len(orderBySlice); i++ {
		orderBySlice[i] = strings.TrimSpace(orderBySlice[i])
		fieldInfo, ok := columnsMap.Load(orderBySlice[i])
		if !ok {
			columns += fmt.Sprintf("%s %s, ", orderBySlice[i], Bigint)
		} else {
			cwt := fieldInfo.(*ColumnWithType)
			if !cwt.Type.Nullable {
				columns += fmt.Sprintf("%s %s NOT NULL, ", cwt.Name, cwt.Type.ToString())
			} else {
				columns += fmt.Sprintf("%s %s, ", cwt.Name, cwt.Type.ToString())
			}
		}
	}
	columnsMap.Range(func(key, value any) bool {
		if key.(string) == "_create_time" || contains(orderBySlice, key.(string)) {
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
		query = fmt.Sprintf(createTableSQLForDuplicate, database, tablename, columns, orderByKey, orderByKey)
	}
	return query
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func delTable(database, tablename string) string {
	var query = fmt.Sprintf(dropTableSQL, database, tablename)
	return query
}

func addTableColumns(database, tablename string, columnName, columnType string, nullAble bool) string {
	var query = fmt.Sprintf(addColumnSQL, database, tablename, columnName, columnType)
	if !nullAble {
		switch columnType {
		case Tinyint, Smallint, Int, Bigint, Float, Double, Decimal, Boolean, Largeint:
			query = query + " NOT NULL DEFAULT 0"
		case Binary, Varchar, Char, String:
			query = query + " NOT NULL DEFAULT \"\""
		case Datetime, Date:
			query = query + " NOT NULL DEFAULT CURRENT_TIMESTAMP"
		}
	}
	return query
}
