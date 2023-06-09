package clause

import (
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators = map[Type]generator{}

func init() {
	generators[INSERT] = genInsert
	generators[VALUES] = genValues
	generators[SELECT] = genSelect
	generators[LIMIT] = genLimit
	generators[ORDERBY] = genOrderBy
	generators[WHERE] = genWhere
	generators[UPDATE] = genUpdate
	generators[COUNT] = genCount
	generators[DELETE] = genDelete
}

func genBindVar(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func genInsert(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ", ")
	return "INSERT INTO " + tableName + " (" + fields + ")", []interface{}{}
}

func genValues(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")

	for i, value := range values {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString("(")
		sql.WriteString(genBindVar(len(value.([]interface{}))))
		sql.WriteString(")")
		vars = append(vars, value.([]interface{})...)
	}

	return sql.String(), vars
}

func genSelect(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ", ")
	return "SELECT " + fields + " FROM " + tableName, []interface{}{}
}

func genLimit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func genOrderBy(values ...interface{}) (string, []interface{}) {
	return "ORDER BY " + values[0].(string), []interface{}{}
}

func genWhere(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("WHERE ")

	for i, expr := range values {
		if i%2 == 0 {
			if i > 0 {
				sql.WriteString(" AND ")
			}
			sql.WriteString(expr.(string))
		} else {
			vars = append(vars, expr)
		}
	}
	return sql.String(), vars
}

// genUpdate generates the UPDATE clause
func genUpdate(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	m := values[1].(map[string]interface{})
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("UPDATE " + tableName + " SET ")
	for k, v := range m {
		sql.WriteString(k + "=?, ")
		vars = append(vars, v)
	}
	var sqlStr = sql.String()
	sqlStr = sqlStr[:len(sqlStr)-2]
	return sqlStr, vars
}

// genDelete generates the DELETE clause
func genDelete(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	return "DELETE FROM " + tableName, []interface{}{}
}

// genCount generates the COUNT clause
func genCount(values ...interface{}) (string, []interface{}) {
	return genSelect(values[0], []string{"COUNT(*)"})
}
