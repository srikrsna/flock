package flock

import (
	"context"
	"database/sql"
)

//ConnectDB -  Return database connection interface when passed the connection string and database
func ConnectDB(u string, databaseSrv string) (*sql.DB, error) {

	db, err := sql.Open(databaseSrv, u)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

//GetData  - Returns select query results
func GetData(ctx context.Context, db *sql.DB, query string, args []interface{}) ([]map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	res := make([]map[string]interface{}, 0, 10)

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colTyp := range cols {
			val := columnPointers[i].(*interface{})
			m[colTyp.Name()] = *val
		}

		res = append(res, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// GetSchema - Return the column names of every table
func GetSchema(ctx context.Context, db *sql.DB) (map[string][]string, error) {

	// TODO : Generalize table query

	tableQuery := `SELECT CONCAT(schemaname, '.', tablename) as TABLE_NAME FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';`
	t, err := db.QueryContext(ctx, tableQuery)
	if err != nil {
		return nil, err
	}
	defer t.Close()

	tables := make([]string, 0)

	for t.Next() {
		var table string

		if err := t.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	res := make(map[string][]string)

	for _, name := range tables {

		// TODO : Generalize column query

		columnQuery := `SELECT * FROM ` + name + ` WHERE 1 = 0;`

		cols, err := db.QueryContext(ctx, columnQuery)
		if err != nil {
			continue
		}

		columns, err := cols.Columns()

		if er := cols.Close(); er != nil {
			return nil, er
		}

		if err != nil {
			return nil, err
		}

		res[name] = columns
	}

	return res, nil
}
