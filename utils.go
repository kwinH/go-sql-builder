package sqlBuilder

import (
	"fmt"
	"reflect"
	"strings"
)

func (b *Builder) initialize() {
	if b.params == nil {
		b.params = make(map[string][]interface{}, 4)
	}
}

func (b *Builder) setTable(table interface{}) (uint8, string, []interface{}) {
	var (
		tmpTable             string
		param                []interface{}
		tmpTableClosureCount uint8
	)

	switch table.(type) {
	case string:
		tmpTable = table.(string)
	case func(*Builder):
		bw := NewBuilder("")
		bw.tmpTableClosureCount = b.tmpTableClosureCount
		bw.tmpTableClosureCount++
		tmpTableClosureCount = bw.tmpTableClosureCount
		table.(func(*Builder))(bw)
		tmpTable, param = bw.ToSql()
		tmpTable = fmt.Sprintf("(%s) as `tmp%d`", tmpTable, tmpTableClosureCount)
	case func() *Builder:
		tmpTableClosureCount = b.tmpTableClosureCount + 1
		tmpTable, param = table.(func() *Builder)().ToSql()
		tmpTable = fmt.Sprintf("(%s) as `tmp%d`", tmpTable, tmpTableClosureCount)
	}

	return tmpTableClosureCount, tmpTable, param
}

func (b *Builder) placeholders(n int) string {
	var s strings.Builder
	for i := 0; i < n-1; i++ {
		s.WriteString("?,")
	}
	if n > 0 {
		s.WriteString("?")
	}
	return s.String()
}

func (b *Builder) escapeId(field interface{}) (fieldStr string) {
	if field, ok := field.(Raw); ok {
		fieldStr += fmt.Sprintf("%s", field)
		return
	}

	comma := ""
	var fieldArr []string
	if field, ok := field.(string); ok {
		fieldArr = strings.Split(field, ",")
	}

	if field, ok := field.([]string); ok {
		fieldArr = field
	}

	if len(fieldArr) > 0 {
		for k, v := range fieldArr {
			if k > 0 {
				comma = ","
			}
			fieldStr += b.strEscapeId(v, comma)
		}
		return
	}

	if field, ok := field.([]interface{}); ok {
		for k, v := range field {
			if k > 0 {
				comma = ","
			}
			switch v.(type) {
			case string:
				fieldStr += b.strEscapeId(v.(string), comma)
			case Raw:
				fieldStr += fmt.Sprintf("%s%s", comma, v)
			}
		}
		return
	}

	return
}

func (b *Builder) strEscapeId(field string, comma string) string {
	var alias, table string

	containsAs := strings.Contains(field, " as ")

	if containsAs || strings.Contains(field, " ") {
		var fieldArr []string
		if containsAs {
			fieldArr = strings.Split(field, " as ")
		} else {
			fieldArr = strings.Split(field, " ")
		}

		field = strings.Trim(fieldArr[0], " ")

		for i := 1; i < len(fieldArr); i++ {
			if fieldArr[i] == "" {
				continue
			}
			alias = " as `" + strings.Trim(fieldArr[i], " ") + "`"
			break
		}

	}

	if strings.Contains(field, ".") {
		fieldArr := strings.Split(field, ".")
		table = strings.Trim(fieldArr[0], " ")
		field = strings.Trim(fieldArr[1], " ")
	}

	if table != "" {
		table = "`" + table + "`."
	}

	leftBracketIndex := strings.Index(field, "(")
	rightBracketIndex := strings.Index(field, ")")

	if leftBracketIndex >= 0 && rightBracketIndex >= 0 {
		param := strings.Trim(field[leftBracketIndex+1:rightBracketIndex], " `")
		field = field[:leftBracketIndex]

		if param != "" && param != "*" {
			param = fmt.Sprintf("`%s`", param)
		}
		field = fmt.Sprintf("%s(%s)", field, param)
	} else {
		field = fmt.Sprintf("`%s`", field)
	}

	return fmt.Sprintf("%s%s%s%s", comma, table, field, alias)
}

func (b *Builder) convertInterfaceSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	vLen := v.Len()
	ret := make([]interface{}, vLen)
	for i := 0; i < vLen; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func (b *Builder) conditions(mode string, boolean string, args ...interface{}) *Builder {
	var conditions string

	b.initialize()

	argsLen := len(args)
	if argsLen == 1 {
		if query, ok := args[0].(func(*Builder)); ok {
			bw := NewBuilder("")
			query(bw)
			conditions = fmt.Sprintf(" %s (%s)", boolean, strings.Join(bw.methods.where, ""))
			b.params[mode] = append(b.params[mode], bw.params[mode]...)
		} else if condition, ok := args[0].(Raw); ok {
			conditions = fmt.Sprintf(" %s %s", boolean, condition)
		}
	} else if argsLen > 1 {
		field := ""
		operator := ""

		var value interface{}
		switch argsLen {
		case 2:
			field = args[0].(string)
			operator = "="
			value = args[1]
		case 3:
			field = args[0].(string)
			operator = args[1].(string)
			value = args[2]
		default:
			field = args[0].(string)
			operator = args[1].(string)
			value = args[2:]
		}

		valueKind := reflect.TypeOf(value).Kind()

		if strings.Contains(operator, "BETWEEN") {
			args = b.convertInterfaceSlice(value)

			b.params[mode] = append(b.params[mode], args[:2]...)
			conditions = fmt.Sprintf(" %s %s %s ? AND ?", boolean, b.escapeId(field), operator)
		} else {
			switch valueKind {
			case reflect.Array:
			case reflect.Slice:
				vi := b.convertInterfaceSlice(value)
				conditions = fmt.Sprintf(" %s %s %s (%s)", boolean, b.escapeId(field), operator, b.placeholders(len(vi)))
				b.params[mode] = append(b.params[mode], vi...)
			case reflect.Func:
				if query, ok := value.(func(*Builder)); ok {
					bw := NewBuilder("")
					query(bw)
					bwSql, bwParams := bw.ToSql()
					if field == "EXISTS" || field == "NOT EXISTS" {
						operator = field
						conditions = fmt.Sprintf(" %s %s (%s)", boolean, operator, bwSql)
					} else {
						conditions = fmt.Sprintf(" %s %s %s (%s)", boolean, b.escapeId(field), operator, bwSql)
					}

					b.params[mode] = append(b.params[mode], bwParams...)
				}
			default:
				if value == "NULL" || value == "NOT NULL" {
					conditions = fmt.Sprintf(" %s %s IS %s", boolean, b.escapeId(field), value)
				} else {
					conditions = fmt.Sprintf(" %s %s %s ?", boolean, b.escapeId(field), operator)
					b.params[mode] = append(b.params[mode], value)
				}
			}
		}
	}

	switch mode {
	case "where":
		b.methods.where = append(b.methods.where, conditions)
	case "having":
		b.methods.having = append(b.methods.having, conditions)
	}

	return b
}

func (b *Builder) cleanLastSql() {
	b.tmpTable = ""
	b.tmpTableClosureCount = 0
	b.methods = methods{}
	b.params = make(map[string][]interface{}, 4)
}
