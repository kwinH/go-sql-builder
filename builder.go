package sqlBuilder

import (
	"fmt"
)

type Raw string

type methods struct {
	field  []interface{}
	where  []string
	order  []string
	limit  string
	group  []string
	having []string
	join   []string
}

type Builder struct {
	TableName            string
	tmpTable             string
	tmpTableClosureCount uint8
	params               map[string][]interface{}

	// 链式操作方法列表
	methods methods
}

func NewBuilder(tableName string) *Builder {
	obj := &Builder{TableName: tableName}
	obj.initialize()
	return obj
}

func (b *Builder) GetField() []interface{} {
	return b.methods.field
}

func (b *Builder) GetWhere() []string {
	return b.methods.where
}

func (b *Builder) GetTable() string {
	if b.tmpTable != "" {
		if b.tmpTableClosureCount == 0 {
			return b.strEscapeId(b.tmpTable, "")
		} else {
			return b.tmpTable
		}
	} else {
		if b.TableName == "" {
			return ""
		}

		return fmt.Sprintf("`%s`", b.TableName)
	}
}

func (b *Builder) GetOrder() []string {
	return b.methods.order
}

func (b *Builder) GetLimit() string {
	return b.methods.limit
}

func (b *Builder) GetGroup() []string {
	return b.methods.group
}

func (b *Builder) GetHaving() []string {
	return b.methods.having
}

func (b *Builder) GetJoin() []string {
	return b.methods.join
}
