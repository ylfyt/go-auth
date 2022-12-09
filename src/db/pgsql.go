package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func snakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func Get[T any](conn DbConnection, query string, params ...interface{}) ([]T, error) {
	rows, err := conn.sqlDb.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columnNum := len(columnTypes)

	var finalValues []T
	for rows.Next() {
		scannedData := make([]interface{}, columnNum)

		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scannedData[i] = new(sql.NullString)
			case "BOOL":
				scannedData[i] = new(sql.NullBool)
			case "INT64":
				scannedData[i] = new(sql.NullInt64)
			case "INT32":
				scannedData[i] = new(sql.NullInt32)
			default:
				scannedData[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scannedData...)
		if err != nil {
			return nil, err
		}

		rowData := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scannedData[i]).(*sql.NullBool); ok {
				rowData[v.Name()] = z.Bool
				continue
			}
			if z, ok := (scannedData[i]).(*sql.NullString); ok {
				rowData[v.Name()] = z.String
				continue
			}
			if z, ok := (scannedData[i]).(*sql.NullInt64); ok {
				rowData[v.Name()] = z.Int64
				continue
			}
			if z, ok := (scannedData[i]).(*sql.NullFloat64); ok {
				rowData[v.Name()] = z.Float64
				continue
			}
			if z, ok := (scannedData[i]).(*sql.NullInt32); ok {
				rowData[v.Name()] = z.Int32
				continue
			}

			rowData[v.Name()] = scannedData[i]
		}

		var tempData T
		for key := range rowData {
			fieldName := snakeCaseToCamelCase(key)
			field := reflect.ValueOf(&tempData).Elem().FieldByName(fieldName)
			if (field == reflect.Value{}) {
				return nil, fmt.Errorf("cannot find %s or %s field in %T", fieldName, key, tempData)
			}
			if field.Type().String() == "uuid.UUID" {
				val, _ := uuid.Parse(rowData[key].(string))
				field.Set(reflect.ValueOf(val))
				continue
			}
			field.Set(reflect.ValueOf(rowData[key]))
		}

		finalValues = append(finalValues, tempData)
	}
	return finalValues, nil
}

func (me DbConnection) GetQuery(query string, params ...interface{}) ([]byte, error) {
	rows, err := me.sqlDb.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {
		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
			case "INT64":
				scanArgs[i] = new(sql.NullInt64)
			case "INT32":
				scanArgs[i] = new(sql.NullInt32)
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		masterData := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	jsonData, err := json.Marshal(finalRows)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (me DbConnection) WriteQuery(query string, params ...interface{}) (int64, error) {
	res, err := me.sqlDb.Exec(query, params...)
	if err != nil {
		return 0, err
	}

	affectedRows, err := res.RowsAffected()
	return affectedRows, err
}

func (me DbConnection) GetRowCount(query string, params ...interface{}) (int, error) {
	rows, err := me.sqlDb.Query(query, params...)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (me DbConnection) GetValue(query string, params ...interface{}) (interface{}, error) {
	rows, err := me.sqlDb.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	count := len(columnTypes)
	//finalRows := []interface{}{};

	for rows.Next() {
		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
			case "INT64":
				scanArgs[i] = new(sql.NullInt64)
			case "INT32":
				scanArgs[i] = new(sql.NullInt32)
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		for i := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				return z.Bool, nil
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				return z.String, nil
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				return z.Int64, nil
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				return z.Float64, nil
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				return z.Int32, nil
			}

			return scanArgs[i], nil
		}
	}

	return nil, nil
}
