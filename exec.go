package sqlBuilder

import (
	"fmt"
	"strings"
)

func (b *Builder) Delete() (string, []interface{}) {
	defer b.cleanLastSql()

	params := make([]interface{}, 0)

	sql := fmt.Sprintf("delete from %s", b.GetTable())
	sql, whereParams := b.builderWhere(sql)

	if len(b.methods.order) > 0 {
		sql += " ORDER BY " + strings.Join(b.methods.order, ",")
	}

	if b.methods.limit != "" {
		sql += b.methods.limit
	}

	params = append(params, whereParams...)

	return sql, params
}

func (b *Builder) Insert(args ...interface{}) (string, []interface{}) {
	params := make([]interface{}, 0)
	sql := ""
	defer b.cleanLastSql()

	var field []string
	var values [][]interface{}
	if len(args) == 2 {
		field, ok := args[0].([]string)
		if query, ok1 := args[1].(func(*Builder)); ok && ok1 {
			bw := NewBuilder("")
			query(bw)
			sql, params := bw.ToSql()
			sql = fmt.Sprintf("INSERT INTO %s (%s) %s", b.GetTable(), b.escapeId(field), sql)

			return sql, params
		}
	}

	for k, arg := range args {
		isContinue := false
		if arg, ok := arg.(map[string]interface{}); ok {

			if k == 0 {
				for f, _ := range arg {
					field = append(field, f)
				}
			}

			value := make([]interface{}, 0)

			for _, v := range field {
				if val, ok := arg[v]; ok {
					value = append(value, val)
				} else {
					isContinue = true
				}
			}

			if isContinue {
				continue
			}

			values = append(values, value)
		}
	}

	sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES", b.GetTable(), b.escapeId(field))

	comma := ""
	for k, value := range values {
		if k > 0 {
			comma = ","
		}
		sql += fmt.Sprintf("%s(%s)", comma, strings.Trim(strings.Repeat("?,", len(value)), ","))
		params = append(params, value...)
	}

	return sql, params
}

func (b *Builder) Update(data map[string]interface{}) (string, []interface{}) {
	defer b.cleanLastSql()

	params := make([]interface{}, 0)
	setVal := ""
	for k, v := range data {
		setVal += b.escapeId(k) + "=?,"
		params = append(params, v)
	}

	setVal = strings.Trim(setVal, ",")

	sql := fmt.Sprintf("UPDATE %s SET %s", b.GetTable(), setVal)
	sql, whereParams := b.builderWhere(sql)
	params = append(params, whereParams...)

	return sql, params
}
