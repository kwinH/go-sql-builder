package sqlBuilder

import (
	"fmt"
	"strings"
)

func (b *Builder) builderWhere(sql string) (string, []interface{}) {
	params := make([]interface{}, 0)

	where := strings.Join(b.methods.where, "")
	if where != "" {
		where = strings.Trim(where, " ")
		sql += fmt.Sprintf(" WHERE %s", where)
	}

	if whereParams, ok := b.params["where"]; ok {
		params = append(params, whereParams...)
	}

	return sql, params
}

func (b *Builder) Where(args ...interface{}) *Builder {
	var boolean string

	if len(b.methods.where) > 0 {
		boolean = "AND"
	}

	b.conditions("where", boolean, args...)

	return b
}

func (b *Builder) OrWhere(args ...interface{}) *Builder {
	var boolean string

	if len(b.methods.where) > 0 {
		boolean = "OR"
	}

	b.conditions("where", boolean, args...)

	return b
}

func (b *Builder) WhereExists(where func(*Builder)) *Builder {
	return b.Where("EXISTS", where)
}

func (b *Builder) WhereNotExists(where func(*Builder)) *Builder {
	return b.Where("NOT EXISTS", where)
}

func (b *Builder) OrWhereExists(where func(*Builder)) *Builder {
	return b.OrWhere("EXISTS", where)
}

func (b *Builder) OrWhereNotExists(where func(*Builder)) *Builder {
	return b.OrWhere("NOT EXISTS", where)
}

func (b *Builder) WhereIn(field string, condition interface{}) *Builder {
	return b.Where(field, "IN", condition)
}

func (b *Builder) WhereNotIn(field string, condition interface{}) *Builder {
	return b.Where(field, "NOT IN", condition)
}

func (b *Builder) OrWhereIn(field string, condition interface{}) *Builder {
	return b.OrWhere(field, "IN", condition)
}

func (b *Builder) OrWhereNotIn(field string, condition interface{}) *Builder {
	return b.OrWhere(field, "NOT IN", condition)
}

func (b *Builder) WhereNull(field string) *Builder {
	return b.Where(field, "NULL")
}

func (b *Builder) WhereNotNull(field string) *Builder {
	return b.Where(field, "NOT NULL")
}

func (b *Builder) OrWhereNull(field string) *Builder {
	return b.OrWhere(field, "NULL")
}

func (b *Builder) OrWhereNotNull(field string) *Builder {
	return b.OrWhere(field, "NOT NULL")
}

func (b *Builder) WhereBetween(field string, condition interface{}) *Builder {
	return b.Where(field, "BETWEEN", condition)
}

func (b *Builder) OrWhereBetween(field string, condition interface{}) *Builder {
	return b.OrWhere(field, "BETWEEN", condition)
}
