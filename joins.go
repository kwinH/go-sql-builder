package sqlBuilder

import "fmt"

func (b *Builder) Joins(table interface{}, condition string, joinType string, params ...interface{}) *Builder {
	b.initialize()

	isClosure, table, param, tableAlias := b.setTable(table)

	if isClosure == 0 {
		table = fmt.Sprintf("`%s`", table)

		if tableAlias != "" {
			table = fmt.Sprintf("%s as `%s`", table, tableAlias)
		}
	}

	b.params["join"] = append(b.params["join"], param...)
	b.params["join"] = append(b.params["join"], params...)

	b.methods.join = append(b.methods.join, fmt.Sprintf("%s JOIN %s %s", joinType, table, condition))
	return b
}

func (b *Builder) LefJoin(table interface{}, condition string, params ...interface{}) *Builder {
	return b.Joins(table, condition, "LEFT", params...)
}

func (b *Builder) RightJoin(table interface{}, condition string, params ...interface{}) *Builder {
	return b.Joins(table, condition, "RIGHT", params...)
}

func (b *Builder) Join(table interface{}, condition string, params ...interface{}) *Builder {
	return b.Joins(table, condition, "INNER", params...)
}
