package db

import (
	"database/sql"
	"encoding/json"
)

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