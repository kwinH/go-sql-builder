package sqlBuilder

import (
	"reflect"
	"testing"
)

func TestBuilder_Insert(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Insert(map[string]interface{}{
		"name": "张三",
		"age":  18,
	})

	if (sql == "INSERT INTO `user` (`name`,`age`) VALUES(?,?)" &&
		reflect.DeepEqual(params, []interface{}{"张三", 18})) ||
		(sql == "INSERT INTO `user` (`age`,`name`) VALUES(?,?)" &&
			reflect.DeepEqual(params, []interface{}{18, "张三"})) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Insert_Multi(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Insert(map[string]interface{}{
		"name": "张三",
		"age":  18,
	}, map[string]interface{}{
		"name": "李四",
		"age":  30,
	})
	if (sql == "INSERT INTO `user` (`name`,`age`) VALUES(?,?),(?,?)" &&
		reflect.DeepEqual(params, []interface{}{"张三", 18, "李四", 30})) ||
		(sql == "INSERT INTO `user` (`age`,`name`) VALUES(?,?),(?,?)" &&
			reflect.DeepEqual(params, []interface{}{18, "张三", 30, "李四"})) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Insert_Select(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Table("user").Insert([]string{"id", "name"}, func(m *Builder) {
		m.Select(Raw("id"), "name").Table("user_old").Where("name", "like", "%q%")
	})
	if sql == "INSERT INTO `user` (`id`,`name`) SELECT id,`name` FROM `user_old` WHERE `name` like ?" &&
		reflect.DeepEqual(params, []interface{}{"%q%"}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Delete(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").Where("id", 1).Delete()
	if sql == "delete from `user` WHERE `id` = ?" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Update(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Table("user").Where("id", 1).Update(map[string]interface{}{
		"name": "test",
		"age":  18,
	})
	if (sql == "UPDATE `user` SET `name`=?,`age`=? WHERE `id` = ?" &&
		reflect.DeepEqual(params, []interface{}{"test", 18, 1})) ||
		(sql == "UPDATE `user` SET `age`=?,`name`=? WHERE `id` = ?" &&
			reflect.DeepEqual(params, []interface{}{18, "test", 1})) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Replace(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Replace(map[string]interface{}{
		"name": "张三",
		"age":  18,
	})

	if (sql == "REPLACE INTO `user` (`name`,`age`) VALUES(?,?)" &&
		reflect.DeepEqual(params, []interface{}{"张三", 18})) ||
		(sql == "REPLACE INTO `user` (`age`,`name`) VALUES(?,?)" &&
			reflect.DeepEqual(params, []interface{}{18, "张三"})) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_DuplicateKey(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").DuplicateKey(map[string]interface{}{
		"age": 18,
	}).Insert(map[string]interface{}{
		"name": "张三",
		"age":  18,
	})

	if (sql == "INSERT INTO `user` (`name`,`age`) VALUES(?,?) ON DUPLICATE KEY UPDATE age=?" &&
		reflect.DeepEqual(params, []interface{}{"张三", 18, 18})) ||
		(sql == "INSERT INTO `user` (`age`,`name`) VALUES(?,?) ON DUPLICATE KEY UPDATE age=?" &&
			reflect.DeepEqual(params, []interface{}{18, "张三", 18})) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}
