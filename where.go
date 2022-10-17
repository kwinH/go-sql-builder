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

func (b *Builder) WhereIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "IN"

	args = append(args, value...)

	return b.Where(args...)
}

func (b *Builder) WhereNotIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT IN"

	args = append(args, value...)

	return b.Where(args...)
}

func (b *Builder) OrWhereIn(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "IN"

	args = append(args, value...)

	return b.OrWhere(args...)
}

func (b *Builder) OrWhereNotIn(field string, value ...interface{}) *Builder {

	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT IN"

	args = append(args, value...)

	return b.OrWhere(args...)
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

func (b *Builder) WhereBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "BETWEEN"

	args = append(args, value...)

	return b.Where(args...)
}

func (b *Builder) OrWhereBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "BETWEEN"

	args = append(args, value...)

	return b.OrWhere(args...)
}

func (b *Builder) WhereNotBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT BETWEEN"

	args = append(args, value...)

	return b.Where(args...)
}

func (b *Builder) OrWhereNotBetween(field string, value ...interface{}) *Builder {
	args := make([]interface{}, 2)
	args[0] = field
	args[1] = "NOT BETWEEN"

	args = append(args, value...)

	return b.OrWhere(args...)
}
