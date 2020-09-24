package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func updateSQL(tableName string, params url.Values, body map[string]string) (string, bool) {
	if len(body) == 0 {
		return statusCodeBadRequest, false
	}
	SQLQuery := "update `" + tableName + "` set "

	init := false
	for key, val := range body {
		if init {
			SQLQuery += ","
		}
		SQLQuery += "`" + key + "` = '" + strings.Replace(val, "'", "\\'", -1) + "' "
		init = true
	}

	SQLQuery += " where "
	init = false
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += "`" + key + "` = '" + strings.Replace(val[0], "'", "\\'", -1) + "' "
		init = true
	}

	logger(SQLQuery)

	_, err = db.Exec(SQLQuery)
	if err != nil {
		fmt.Println("updateSQL", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			return sqlErrorCheck(driverErr.Number), false
		}
		return statusCodeServerError, false // default
	}
	return statusCodeOk, true
}

func insertSQL(tableName string, body map[string]string) (string, bool) {
	if len(body) == 0 {
		return statusCodeBadRequest, false
	}
	SQLQuery := buildInsertStatement(tableName, body)
	logger(SQLQuery)

	_, err = db.Exec(SQLQuery)
	if err != nil {
		fmt.Println("insertSQL", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			errCode := sqlErrorCheck(driverErr.Number)
			if strings.EqualFold(errCode, statusCodeDuplicateEntry) {
				if strings.Contains(err.Error(), "_u_id") {
					return statusCodeDuplicateEntry, false
				}
				return statusCodeBadRequest, false
			}
			return errCode, false
		}
		return statusCodeServerError, false // default
	}
	return statusCodeCreated, true
}

func buildInsertStatement(tableName string, body map[string]string) string {
	SQLQuery := "insert into `" + tableName + "` "
	keys := " ("
	values := " ("
	init := false
	for key, val := range body {
		if init {
			keys += ","
			values += ","
		}
		keys += "`" + key + "`"
		values += "'" + strings.Replace(val, "'", "\\'", -1) + "'"
		init = true
	}
	keys += ")"
	values += ")"
	SQLQuery += keys + " values " + values
	return SQLQuery
}

func deleteSQL(tableName string, params url.Values) (string, bool) {
	if len(params) == 0 {
		return statusCodeBadRequest, false
	}
	SQLQuery := "delete from `" + tableName + "` where "

	init := false
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += "`" + key + "` = '" + strings.Replace(val[0], "'", "\\'", -1) + "' "
		init = true
	}
	logger(SQLQuery)

	_, err = db.Exec(SQLQuery)
	if err != nil {
		fmt.Println("deleteSQL", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			return sqlErrorCheck(driverErr.Number), false
		}
		return statusCodeServerError, false // default
	}
	return statusCodeOk, true
}

func selectSQL(tableName string, params url.Values) ([]map[string]string, string, bool) {
	where := ""
	init := false
	for key, val := range params {
		if init {
			where += " and "
		}
		where += "`" + key + "` = '" + strings.Replace(val[0], "'", "\\'", -1) + "' "
		init = true
	}
	SQLQuery := "select * from `" + tableName + "`"
	if strings.Compare(where, "") != 0 {
		SQLQuery += " where " + where
	}
	return selectProcess(SQLQuery)
}

func selectProcess(SQLQuery string) ([]map[string]string, string, bool) {
	logger(SQLQuery)

	// SQLQuery = html.EscapeString(SQLQuery)
	rows, err := db.Query(SQLQuery)
	if err != nil {
		fmt.Println("selectProcess", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr != nil {
			return []map[string]string{}, sqlErrorCheck(driverErr.Number), false
		}
		return []map[string]string{}, statusCodeServerError, false // default
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("selectProcess", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			return []map[string]string{}, sqlErrorCheck(driverErr.Number), false
		}
		return []map[string]string{}, statusCodeServerError, false // default
	}

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	data := []map[string]string{}
	rest := map[string]string{}
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rest = map[string]string{}
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("selectProcess", err)
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				return []map[string]string{}, sqlErrorCheck(driverErr.Number), false
			}
			return []map[string]string{}, statusCodeServerError, false // default
		}

		for i, raw := range rawResult {
			if raw == nil {
				rest[cols[i]] = ""
			} else {
				rest[cols[i]] = string(raw)
			}
		}

		data = append(data, rest)
	}
	return data, statusCodeOk, true
}

func selectProcessNoLogging(SQLQuery string) ([]map[string]string, string, bool) {

	// SQLQuery = html.EscapeString(SQLQuery)
	rows, err := db.Query(SQLQuery)
	if err != nil {
		fmt.Println("selectProcess", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr != nil {
			return []map[string]string{}, sqlErrorCheck(driverErr.Number), false
		}
		return []map[string]string{}, statusCodeServerError, false // default
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("selectProcess", err)
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			return []map[string]string{}, sqlErrorCheck(driverErr.Number), false
		}
		return []map[string]string{}, statusCodeServerError, false // default
	}

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	data := []map[string]string{}
	rest := map[string]string{}
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rest = map[string]string{}
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("selectProcess", err)
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				return []map[string]string{}, sqlErrorCheck(driverErr.Number), false
			}
			return []map[string]string{}, statusCodeServerError, false // default
		}

		for i, raw := range rawResult {
			if raw == nil {
				rest[cols[i]] = ""
			} else {
				rest[cols[i]] = string(raw)
			}
		}

		data = append(data, rest)
	}
	return data, statusCodeOk, true
}
