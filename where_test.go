package sqlBuilder

import (
	"reflect"
	"testing"
)

func TestBuilder_Where(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id", "name").
		Table("users").
		Where("id", 1).
		OrWhere("name", "like", "%q%").
		ToSql()

	if sql == "SELECT `id`,`name` FROM `users` WHERE `id` = ? OR `name` like ?" &&
		reflect.DeepEqual(params, []interface{}{1, "%q%"}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WHERERaw(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where(Raw("price > IF(state = 'TX', 200, 100)")).ToSql()
	if sql == "SELECT * FROM `user` WHERE price > IF(state = 'TX', 200, 100)" &&
		reflect.DeepEqual(params, []interface{}{}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_Where_Closure(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Select("id").
		Where("id", "<>", 1).
		OrWhere(func(m *Builder) {
			m.Where("age", ">", 18).
				Where("name", "like", "%q%")
		}).ToSql()

	if sql == "SELECT `id` FROM `user` WHERE `id` <> ? OR (  `age` > ? AND `name` like ?)" &&
		reflect.DeepEqual(params, []interface{}{1, 18, "%q%"}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

/**
where 子查询
*/
func TestBuilder_Where_SubQuery(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("id", "<>", 1).
		WhereIn("id", func(m *Builder) {
			m.Select("id").
				Table("user_old").
				Where("age", ">", 18).
				Where("name", "like", "%q%")
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `id` <> ? AND `id` IN (SELECT `id` FROM `user_old` WHERE `age` > ? AND `name` like ?)" &&
		reflect.DeepEqual(params, []interface{}{1, 18, "%q%"}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereIn(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereIn("id", []int{100, 200}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `id` IN (?,?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100, 200}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereIn("id", func(m *Builder) {
			m.Select("id").
				Table("user_old").
				Where("id", "<", 100)
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `id` IN (SELECT `id` FROM `user_old` WHERE `id` < ?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereIn(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereIn("id", []int{100, 200}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `id` IN (?,?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100, 200}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereIn("id", func(m *Builder) {
			m.Select("id").
				Table("user_old").
				Where("id", "<", 100)
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `id` IN (SELECT `id` FROM `user_old` WHERE `id` < ?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereNotIn(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereNotIn("id", []int{100, 200}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `id` NOT IN (?,?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100, 200}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereNotIn("id", func(m *Builder) {
			m.Select("id").
				Table("user_old").
				Where("id", "<", 100)
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `id` NOT IN (SELECT `id` FROM `user_old` WHERE `id` < ?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereNotIn(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereNotIn("id", []int{100, 200}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `id` NOT IN (?,?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100, 200}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereNotIn("id", func(m *Builder) {
			m.Select("id").
				Table("user_old").
				Where("id", "<", 100)
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `id` NOT IN (SELECT `id` FROM `user_old` WHERE `id` < ?)" &&
		reflect.DeepEqual(params, []interface{}{1, 100}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereExists(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereExists(func(m *Builder) {
			m.Select(Raw("1")).
				Table("user_old").
				Where(Raw("user_old.id = user.id"))
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND EXISTS (SELECT 1 FROM `user_old` WHERE user_old.id = user.id)" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereExists(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereExists(func(m *Builder) {
			m.Select(Raw("1")).
				Table("user_old").
				Where(Raw("user_old.id = user.id"))
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR EXISTS (SELECT 1 FROM `user_old` WHERE user_old.id = user.id)" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereNotExists(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereNotExists(func(m *Builder) {
			m.Select(Raw("1")).
				Table("user_old").
				Where(Raw("user_old.id = user.id"))
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND NOT EXISTS (SELECT 1 FROM `user_old` WHERE user_old.id = user.id)" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereNotExists(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereNotExists(func(m *Builder) {
			m.Select(Raw("1")).
				Table("user_old").
				Where(Raw("user_old.id = user.id"))
		}).ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR NOT EXISTS (SELECT 1 FROM `user_old` WHERE user_old.id = user.id)" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereNull(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereNull("deleted_at").ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `deleted_at` IS NULL" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereNotNull(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereNotNull("deleted_at").ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `deleted_at` IS NOT NULL" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereNull(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereNull("deleted_at").ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `deleted_at` IS NULL" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereNotNull(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereNotNull("deleted_at").ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `deleted_at` IS NOT NULL" &&
		reflect.DeepEqual(params, []interface{}{1}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_WhereBetween(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		WhereBetween("attribute", []int{2, 3}).
		ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? AND `attribute` BETWEEN ? AND ?" &&
		reflect.DeepEqual(params, []interface{}{1, 2, 3}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}

func TestBuilder_OrWhereBetween(t *testing.T) {
	var (
		sql    string
		params []interface{}
	)

	sql, params = NewBuilder("user").Where("sex", 1).
		OrWhereBetween("attribute", []int{2, 3}).
		ToSql()

	if sql == "SELECT * FROM `user` WHERE `sex` = ? OR `attribute` BETWEEN ? AND ?" &&
		reflect.DeepEqual(params, []interface{}{1, 2, 3}) {
		t.Log(sql, params)
	} else {
		t.Error(sql, params)
	}
}
