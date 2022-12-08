package db

import (
	"database/sql"
	"encoding/json"
)

func (me DbConnection) GetQuery(query string, params ...interface{}) []byte {
	rows, err := me.sqlDb.Query(query, params...)

	if err != nil {
		// L.Error(err, "DBPgSQL.Connect|Failed to get query %s (q: %s, p: %s)", err, query, params)
		return []byte("[]")
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return nil
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {
		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT64":
				scanArgs[i] = new(sql.NullInt64)
				break
			case "INT32":
				scanArgs[i] = new(sql.NullInt32)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil
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
	return jsonData
}

func (me DbConnection) WriteQuery(query string, params ...interface{}) int64 {
	res, err := me.sqlDb.Exec(query, params...)

	if err != nil {
		// L.Error(err, err.Error())
		return 0
	}

	//utils.Ldbg("Query Exec %s, res %s", query, res)
	affectedRows, _ := res.RowsAffected()

	return affectedRows
}

func (me DbConnection) GetValue(query string, params ...interface{}) interface{} {
	rows, err := me.sqlDb.Query(query, params...)

	if err != nil {
		// L.Error(err, "DBPgSQL.Connect|Failed to get query %s (q: %s, p: %s)", err, query, params)
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return ""
	}

	count := len(columnTypes)
	//finalRows := []interface{}{};

	for rows.Next() {
		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT64":
				scanArgs[i] = new(sql.NullInt64)
				break
			case "INT32":
				scanArgs[i] = new(sql.NullInt32)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil
		}

		for i, _ := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				return z.Bool
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				return z.String
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				return z.Int64
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				return z.Float64
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				return z.Int32
			}

			return scanArgs[i]
		}
	}

	return nil
}
