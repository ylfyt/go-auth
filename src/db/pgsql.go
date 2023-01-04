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

func getFieldNames(dataRef reflect.Type, columnTypes []*sql.ColumnType) map[string]string {
	var fieldNames = make(map[string]string)
	for i := 0; i < dataRef.NumField(); i++ {
		field := dataRef.Field(i)
		columnName := strings.TrimSpace(field.Tag.Get("col"))
		if columnName == "" {
			continue
		}
		fieldNames[columnName] = field.Name
	}

	for _, v := range columnTypes {
		if fieldNames[v.Name()] != "" {
			continue
		}
		camelCase := snakeCaseToCamelCase(v.Name())
		isExist := false
		for _, val := range fieldNames {
			if val == camelCase {
				isExist = true
				break
			}
		}
		if !isExist {
			fieldNames[v.Name()] = camelCase
		}
	}
	return fieldNames
}

func getInterfaceValue(data interface{}) interface{} {
	if z, ok := (data).(*sql.NullBool); ok {
		return z.Bool
	}
	if z, ok := (data).(*sql.NullString); ok {
		return z.String
	}
	if z, ok := (data).(*sql.NullInt64); ok {
		return z.Int64
	}
	if z, ok := (data).(*sql.NullFloat64); ok {
		return z.Float64
	}
	if z, ok := (data).(*sql.NullInt32); ok {
		return z.Int32
	}

	return data
}

func getDataContainer(columnTypes []*sql.ColumnType) []interface{} {
	if len(columnTypes) == 0 {
		return nil
	}
	container := make([]interface{}, len(columnTypes))
	for i, v := range columnTypes {
		switch v.DatabaseTypeName() {
		case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
			container[i] = new(sql.NullString)
		case "BOOL":
			container[i] = new(sql.NullBool)
		case "INT8":
			container[i] = new(sql.NullInt64)
		case "INT4":
			container[i] = new(sql.NullInt32)
		default:
			container[i] = new(sql.NullString)
		}
	}
	return container
}

func getData[T any](onlyOneRow bool, conn *sql.DB, query string, params ...interface{}) ([]T, error) {
	rows, err := conn.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	scannedData := getDataContainer(columnTypes)
	if scannedData == nil {
		return nil, fmt.Errorf("there is no field in query")
	}

	var tempData T
	fieldNames := getFieldNames(reflect.TypeOf(tempData), columnTypes)

	var finalValues []T = make([]T, 0)
	for rows.Next() {
		err := rows.Scan(scannedData...)
		if err != nil {
			return nil, err
		}

		for i, v := range columnTypes {
			fieldName := fieldNames[v.Name()]
			if fieldName == "" {
				continue
			}
			field := reflect.ValueOf(&tempData).Elem().FieldByName(fieldName)
			if (field == reflect.Value{}) {
				continue
			}

			var val interface{} = getInterfaceValue(scannedData[i])
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

func Get[T any](conn *sql.DB, query string, params ...interface{}) ([]T, error) {
	return getData[T](false, conn, query, params...)
}

func GetOne[T any](conn *sql.DB, query string, params ...interface{}) (*T, error) {
	result, err := getData[T](true, conn, query, params...)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func Write(conn *sql.DB, query string, params ...interface{}) (int64, error) {
	res, err := conn.Exec(query, params...)
	if err != nil {
		return 0, err
	}

	affectedRows, err := res.RowsAffected()
	return affectedRows, err
}

func GetRowCount(conn *sql.DB, query string, params ...interface{}) (int, error) {
	rows, err := conn.Query(query, params...)

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

func GetFieldOne[T any](conn *sql.DB, query string, params ...interface{}) (T, error) {
	rows, err := conn.Query(query, params...)
	var result T
	if err != nil {
		return result, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return result, err
	}

	scannedData := getDataContainer(columnTypes)
	if scannedData == nil {
		return result, fmt.Errorf("there is no field in query")
	}

	if !rows.Next() {
		return result, nil
	}

	err = rows.Scan(scannedData...)
	if err != nil {
		return result, err
	}

	var val interface{} = getInterfaceValue(scannedData[0])

	typeName := reflect.TypeOf(result).String()
	if typeName == "uuid.UUID" {
		newUuid, err := uuid.Parse(val.(string))
		if err != nil {
			return result, err
		}
		var iUuid interface{} = newUuid
		return iUuid.(T), nil
	}

	if typeName == "time.Time" {
		if val.(string) == "" {
			return result, nil
		}
		newTime, err := time.Parse(time.RFC3339Nano, val.(string))
		if err != nil {
			return result, err
		}
		var iTime interface{} = newTime
		return iTime.(T), nil
	}

	if z, ok := val.(T); ok {
		return z, nil
	}

	return result, fmt.Errorf("return type is not valid")
}
