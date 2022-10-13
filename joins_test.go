package sqlBuilder

import (
	"reflect"
	"testing"
)

func TestBuilder_Join(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").
		Join("order o", "o.user_id=u.user_id and o.type=?", 1).
		Join("contacts c", "c.user_id=u.user_id").
		ToSql()

	if sql == "SELECT `id`,`name` FROM `user` INNER JOIN `order` as `o` o.user_id=u.user_id and o.type=? INNER JOIN `contacts` as `c` c.user_id=u.user_id" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

}

func TestBuilder_Join_SubQuery(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)
	sql, params = NewBuilder("user").Table("user u").Select("id", "name").
		Join(func(b *Builder) {
			b.Table("contacts").Where("id", ">", 100)
		}, "tmp1.user_id=u.user_id").
		ToSql()

	if sql == "SELECT `id`,`name` FROM `user` as `u` INNER JOIN (SELECT * FROM `contacts` WHERE `id` > ?) as `tmp1` tmp1.user_id=u.user_id" &&
		reflect.DeepEqual(params, []interface{}{100}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_LeftJoin(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").
		LefJoin("contacts c", "user.user_id=u.user_id").
		ToSql()

	if sql == "SELECT `id`,`name` FROM `user` LEFT JOIN `contacts` as `c` user.user_id=u.user_id" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}
