package sqlBuilder

import (
	"fmt"
	"strings"
)

func (b *Builder) Select(args ...interface{}) *Builder {
	b.methods.field = make([]interface{}, 0)

	if len(args) == 1 {
		fieldArr := make([]string, 0)
		if field, ok := args[0].(string); ok {
			fieldArr = strings.Split(field, ",")
		} else if field, ok := args[0].(Raw); ok {
			b.methods.field = append(b.methods.field, field)
			return b
		} else if field, ok := args[0].([]string); ok {
			fieldArr = field
		}

		for _, v := range fieldArr {
			b.methods.field = append(b.methods.field, v)
		}

	} else {
		b.methods.field = append(b.methods.field, args...)
	}

	return b
}

func (b *Builder) Table(table interface{}) *Builder {
	b.initialize()
	b.tmpTableClosureCount, b.tmpTable, b.params["table"] = b.setTable(table)
	return b
}

func (b *Builder) builderHaving(sql string) (string, []interface{}) {
	params := make([]interface{}, 0)

	having := strings.Join(b.methods.having, "")
	if having != "" {
		having = strings.Trim(having, " ")
		sql += fmt.Sprintf(" HAVING %s", having)
	}

	if havingParams, ok := b.params["having"]; ok {
		params = append(params, havingParams...)
	}

	return sql, params
}

func (b *Builder) builderOrder(sql string) string {
	if len(b.methods.order) > 0 {
		sql += " ORDER BY " + strings.Join(b.methods.order, ",")
	}

	return sql
}

func (b *Builder) Group(group ...string) *Builder {
	b.methods.group = append(b.methods.group, group...)

	return b
}

func (b *Builder) Having(args ...interface{}) *Builder {
	var boolean string

	if len(b.methods.having) > 0 {
		boolean = "AND"
	}

	b.conditions("having", boolean, args...)

	return b
}

func (b *Builder) OrHaving(args ...interface{}) *Builder {
	var boolean = ""
	if len(b.methods.having) > 0 {
		boolean = "OR"
	}

	b.conditions("having", boolean, args...)

	return b
}

func (b *Builder) Order(args ...interface{}) *Builder {
	var (
		field string
		value string
	)
	field, ok := args[0].(string)
	if !ok {
		return b
	}

	if len(args) == 1 {
		value = "DESC"
	} else if value, ok = args[1].(string); !ok {
		return b
	}

	value = strings.ToUpper(value)

	b.methods.order = append(b.methods.order, b.escapeId(field)+" "+value)

	return b
}

//
// Limit
// @Description: 指定查询数量
// @receiver b
// @param int64 offset 起始位置
// @param int64 length 查询数量
// @return *Builder
//
func (b *Builder) Limit(args ...int64) *Builder {

	switch len(args) {
	case 1:
		b.methods.limit = fmt.Sprintf(" LIMIT %d", args[0])
	case 2:
		b.methods.limit = fmt.Sprintf(" LIMIT %d,%d", args[0], args[1])
	}

	return b
}

// Page 指定分页
// param int64 page 页数
// param int64 listRows 每页数量
// return *Builder
func (b *Builder) Page(page int64, listRows int64) *Builder {
	b.methods.limit = fmt.Sprintf(" LIMIT %d,%d", (page-1)*listRows, listRows)
	return b
}

func (b *Builder) ToSql() (string, []interface{}) {
	defer b.cleanLastSql()

	params := make([]interface{}, 0)

	fieldStr := ""
	if len(b.methods.field) == 0 {
		fieldStr = "*"
	} else {
		fieldStr = b.escapeId(b.methods.field)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s", fieldStr, b.GetTable())

	if tableParams, ok := b.params["table"]; ok {
		params = append(params, tableParams...)
	}

	if len(b.methods.join) > 0 {
		sql += " " + strings.Join(b.methods.join, " ")

		if joinParams, ok := b.params["join"]; ok {
			params = append(params, joinParams...)
		}
	}

	sql, whereParams := b.builderWhere(sql)
	params = append(params, whereParams...)

	if len(b.methods.group) > 0 {
		sql += " GROUP BY " + b.escapeId(b.methods.group)
	}

	sql, havingParams := b.builderHaving(sql)
	params = append(params, havingParams...)

	sql = b.builderOrder(sql)

	if b.methods.limit != "" {
		sql += b.methods.limit
	}

	return sql, params
}
