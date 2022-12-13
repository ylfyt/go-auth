package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

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

func getData[T any](onlyOneRow bool, conn DbConnection, query string, params ...interface{}) ([]T, error) {
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
	if columnNum == 0 {
		return nil, fmt.Errorf("there is no field in query")
	}

	var fieldNames = make(map[string]string)
	for _, v := range columnTypes {
		camelCase := snakeCaseToCamelCase(v.Name())
		fieldNames[v.Name()] = camelCase
	}

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

	var finalValues []T
	for rows.Next() {
		err := rows.Scan(scannedData...)
		if err != nil {
			return nil, err
		}

		var tempData T
		for i, v := range columnTypes {
			var val interface{}
			if z, ok := (scannedData[i]).(*sql.NullBool); ok {
				val = z.Bool
			} else if z, ok := (scannedData[i]).(*sql.NullString); ok {
				val = z.String
			} else if z, ok := (scannedData[i]).(*sql.NullInt64); ok {
				val = z.Int64
			} else if z, ok := (scannedData[i]).(*sql.NullFloat64); ok {
				val = z.Float64
			} else if z, ok := (scannedData[i]).(*sql.NullInt32); ok {
				val = z.Int32
			} else {
				val = scannedData[i]
			}

			field := reflect.ValueOf(&tempData).Elem().FieldByName(fieldNames[v.Name()])
			if (field == reflect.Value{}) {
				return nil, fmt.Errorf("cannot find %s or %s field in %T", fieldNames[v.Name()], v.Name(), tempData)
			}
			if field.Type().String() != "uuid.UUID" {
				field.Set(reflect.ValueOf(val))
				continue
			}

			val, err := uuid.Parse(val.(string))
			if err != nil {
				return nil, err
			}
			field.Set(reflect.ValueOf(val))
		}

		if onlyOneRow {
			return []T{tempData}, nil
		}
		finalValues = append(finalValues, tempData)
	}
	return finalValues, nil
}

func Get[T any](conn DbConnection, query string, params ...interface{}) ([]T, error) {
	return getData[T](false, conn, query, params...)
}

func GetFirst[T any](conn DbConnection, query string, params ...interface{}) (*T, error) {
	result, err := getData[T](true, conn, query, params...)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func Write(conn DbConnection, query string, params ...interface{}) (int64, error) {
	res, err := conn.sqlDb.Exec(query, params...)
	if err != nil {
		return 0, err
	}

	affectedRows, err := res.RowsAffected()
	return affectedRows, err
}

func GetRowCount(conn DbConnection, query string, params ...interface{}) (int, error) {
	rows, err := conn.sqlDb.Query(query, params...)

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

func GetFieldFirst[T any](conn DbConnection, query string, params ...interface{}) (*T, error) {
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
	if columnNum == 0 {
		return nil, fmt.Errorf("there is no field in query")
	}

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

	if rows.Next() {
		err := rows.Scan(scannedData...)
		if err != nil {
			return nil, err
		}

		var val interface{}
		if z, ok := (scannedData[0]).(*sql.NullBool); ok {
			val = z.Bool
		} else if z, ok := (scannedData[0]).(*sql.NullString); ok {
			val = z.String
		} else if z, ok := (scannedData[0]).(*sql.NullInt64); ok {
			val = z.Int64
		} else if z, ok := (scannedData[0]).(*sql.NullFloat64); ok {
			val = z.Float64
		} else if z, ok := (scannedData[0]).(*sql.NullInt32); ok {
			val = z.Int32
		} else {
			val = scannedData[0]
		}

		var data T
		typeName := reflect.TypeOf(data).String()
		if typeName == "uuid.UUID" {
			newUuid, err := uuid.Parse(val.(string))
			if err != nil {
				return nil, err
			}
			var iUuid interface{} = newUuid
			data = iUuid.(T)

			return &data, nil
		}

		if typeName == "time.Time" {
			if val.(string) == "" {
				return nil, nil
			}
			newTime, err := time.Parse(time.RFC3339Nano, val.(string))
			if err != nil {
				return nil, err
			}
			var iTime interface{} = newTime
			data = iTime.(T)

			return &data, nil
		}

		data = val.(T)
		return &data, nil
	}

	return nil, nil
}
