# 简介[![Go Reference](https://pkg.go.dev/badge/github.com/kwinH/go-sql-builder.svg)](https://pkg.go.dev/github.com/kwinH/go-sql-builder)

**一款超好用Golang版SQL构造器，支持链式操作拼接SQL语句，单元测试覆盖率100%，详细用法请查看测试用例。**

## 使用
```bash
go get -v github.com/kwinH/go-sql-builder
```

# SQL构造器

```go
user := Builder{
TableName: "user",
}
```

## Select

> 查询字段 默认查询所有，即`*`
>
> 支持多个参数，单个参数，切片方式传值

```go

// SELECT `id`,`name` FROM `user`
user.Select("id", "name").ToSql()

user.Select("id,name").ToSql()

user.Select([]string{"id","name"}).ToSql()
```

## 原生表达式

> 有时候你可能需要在查询中使用原生表达式。你可以使用 `sqlBuilder.Raw` 创建一个原生表达式：

```go
// SELECT count(*) as c FROM `user` []
sql, params = user.Select(Raw("count(*) as c")).ToSql()
```

```go 
// SELECT * FROM `user` WHERE price > IF(state = 'TX', 200, 100) []
sql, params = user.Where(Raw("price > IF(state = 'TX', 200, 100)")).ToSql()
```

## Table

> 指定查询表名

```go
//SELECT * FROM `users`
user.Table("users").ToSql()
```

### 子查询

```go
// SELECT * FROM (SELECT * FROM (SELECT `sex`,count(*) as c FROM m_users GROUP BY `sex`) as `tmp2`) as `tmp1` []
sql, params = user.Table(func (m *Builder) {
m.Table(func (m *Builder) {
m.Table("m_users").Select("sex", Raw("count(*) as c")).Group("sex")
})
```

## Where

### 简单where语句

> 在构造 where 查询实例中，你可以使用 where 方法。调用 where 最基本的方式是需要传递三个参数：第一个参数是列名，第二个参数是任意一个数据库系统支持的运算符，第三个是该列要比较的值。

```go
// SELECT `id`,`name` FROM `users` WHERE  `id` = ?  [1]
sql, params = user.
Where("id", "=", 1).
ToSql()
```

> 为了方便，如果你只是简单比较列值和给定数值是否相等，可以将数值直接作为 where 方法的第二个参数：

```go
// SELECT `id`,`name` FROM `users` WHERE  `id` = ? [1]
sql, params = user.
Where("id", 1).
ToSql()
```

### OrWhere语句

> orWhere 方法和 where 方法接收的参数一样：

```go
// SELECT * FROM `user` WHERE  `id` = ? OR  `name` like ? [1 %q%]
sql, params = user.
Where("id", 1).
OrWhere("name", "like", "%q%").
ToSql()
```

### 参数分组

> 如果需要在括号内对 or 条件进行分组，将闭包作为 orWhere 方法的第一个参数也是可以的：

```go
// SELECT `id` FROM `user` WHERE  `id` <> ? OR  ( `age` > ? AND  `name` like ?) [1 18 %q%]
sql, params := user.Where("id", "<>", 1).
OrWhere(func (m *Builder) {
m.Where("age", ">", 18).
Where("name", "like", "%q%")
}).ToSql()
```

### 子查询 Where 语句

```go
// SELECT * FROM `user` WHERE  `id` <> ? AND  `id` in (SELECT `id` FROM `user_old` WHERE  `age` > ? AND  `name` like ?) [1 18 %q%]
sql, params = user.Where("id", "<>", 1).
WhereIn("id", func (m *Builder) {
m.Select("id").
Table("user_old").
Where("age", ">", 18).
Where("name", "like", "%q%")
}).ToSql()
```

### whereIn / whereNotIn / orWhereIn / orWhereNotIn

> `WhereIn` 方法验证给定列的值是否包含在给定数组中：

```go
// SELECT * FROM `user` WHERE `sex` = ? AND `id` IN (?,?) [1 100 200]
sql, params = user.Where("sex", 1).
WhereIn("id", []int{100, 200}).ToSql()

// SELECT * FROM `user` WHERE `sex` = ? OR `id` IN (?,?) [1 100 200]
sql, params = user.Where("sex", 1).
OrWhereIn("id", []int{100, 200}).ToSql()
```

> `WhereNotIn` 方法验证给定列的值是否`不存在`给定的数组中：

```go
// SELECT * FROM `user` WHERE `sex` = ? AND `id` NOT IN (?,?) [1 100 200]
sql, params = user.Where("sex", 1).
WhereNotIn("id", []int{100, 200}).ToSql()

// SELECT * FROM `user` WHERE `sex` = ? OR `id` NOT IN (?,?) [1 100 200]
sql, params = user.Where("sex", 1).
OrWhereNotIn("id", []int{100, 200}).ToSql()
```

### ＷhereNull / ＷhereNotNull / ＯrWhereNull / ＯrWhereNotNull

＞`ＷhereNull` 方法验证指定的字段`必须是 NULL`:

```go
// SELECT * FROM `user` WHERE `sex` = ? AND `deleted_at` IS NULL [1]
sql, params = user.Where("sex", 1).
WhereNull("deleted_at").ToSql()

// SELECT * FROM `user` WHERE `sex` = ? OR `deleted_at` IS NULL [1]
sql, params = user.Where("sex", 1).
OrWhereNull("deleted_at").ToSql()

```

> `WhereNotNull` 方法验证指定的字段`肯定不是 NULL`:

```go
// SELECT * FROM `user` WHERE `sex` = ? AND `deleted_at` IS NOT NULL [1]
sql, params = user.Where("sex", 1).
WhereNotNull("deleted_at").ToSql()

// SELECT * FROM `user` WHERE `sex` = ? OR `deleted_at` IS NOT NULL [1]
sql, params = user.Where("sex", 1).
OrWhereNotNull("deleted_at").ToSql()

```

## Order

> `Order`方法允许你通过给定字段对结果集进行排序。 `order`
> 的第一个参数应该是你希望排序的字段，第二个参数控制排序的方向，可以是 `asc` 或 `desc`,也可以省略，默认是`desc`

```go
// SELECT `id`,`name` FROM `user` ORDER BY `id` DESC []
sql, params = user.Select("id", "name").
Order("id", "desc").
ToSql()
```

> 如果你需要使用多个字段进行排序，你可以多次调用 `Order`

```go
// SELECT `id`,`name` FROM `user` ORDER BY `id` DESC,`age` ASC []
sql, params = user.Select("id", "name").
Order("id").
Order("age", "asc").
ToSql()

```

## groupBy / Having

> groupBy 和 having 方法用于将结果分组。 having 方法的使用与 where 方法十分相似：

```go
// SELECT `age`,count(*) as c FROM `user` GROUP BY `age` HAVING  `c` > ? [10]
sql, params = user.Select("age", Raw("count(*) as c")).Group("age").Having("c", ">", 10).ToSql()

// SELECT `age`,`sex`,count(*) as c FROM `user` GROUP BY `age`,`sex` HAVING  `c` > ? [10]
sql, params = user.Select("age", 'sex', Raw("count(*) as c")).Group("age", "sex").Having("c", ">", 10).ToSql()
```

## Limit

```go
// SELECT `id`,`name` FROM `user` LIMIT 10 []
sql, params = user.Select("id", "name").Limit(10).ToSql()

// SELECT `id`,`name` FROM `user` LIMIT 1,10 []
sql, params = user.Select("id", "name").Limit(1, 10).ToSql()

```

## Page

```go
//SELECT `id`,`name` FROM `user` LIMIT 0,10 []
sql, params = user.Select("id", "name").Page(1, 10).ToSql()
```

## Joins

### Inner Join 语句

> 查询构造器也可以编写 join 方法。若要执行基本的「内链接」，你可以在查询构造器实例上使用 Join 方法。传递给 Join
> 方法的第一个参数是你需要连接的表的名称，第二个参数是指定连接的字段约束，而其他的则是绑定参数。你还可以在单个查询中连接多个数据表：

```go
// SELECT `id`,`name` FROM `user` INNER JOIN `order` as `o` o.user_id=u.user_id and o.type=? INNER JOIN `contacts` as `c` c.user_id=u.user_id [1]
sql, params = user.Select("id", "name").
Join("order o", "o.user_id=u.user_id and o.type=?", 1).
Join("contacts c", "c.user_id=u.user_id").
ToSql()
```

### Left Join / Right Join 语句

> 如果你想使用 「左连接」或者 「右连接」代替「内连接」 ，可以使用 ＬeftJoin 或者 ＲightJoin 方法。这两个方法与 Join 方法用法相同：

```go
// SELECT `id`,`name` FROM `user` RIGHT JOIN `contacts` as `c` c.user_id=u.user_id []
sql, params = user.Select("id", "name").
LeftJoin("contacts c", "c.user_id=u.user_id").
ToSql()

// SELECT `id`,`name` FROM `user` LEFT JOIN `contacts` as `c` c.user_id=u.user_id []
sql, params = user.Select("id", "name").
RightJoin("contacts c", "c.user_id=u.user_id").
ToSql()
```

### 关联子查询

```go
// SELECT `id`,`name` FROM `user` as `u` INNER JOIN (SELECT * FROM `contacts` WHERE `id` > ?) as `tmp1` tmp1.user_id=u.user_id [100]
sql, params = user.Table("user u").Select("id", "name").
Join(func(b *Builder) {
b.Table("contacts").Where("id", ">", 100)
}, "tmp1.user_id=u.user_id").
ToSql()
```

## 插入

> 查询构造器还提供了 `insert` 方法用于插入记录到数据库中。 `insert` 方法接收数组形式的字段名和字段值进行插入操作：

```go
// INSERT INTO `user` (`name`,`age`) VALUES(?,?) [张三 18]
sql, params, err = user.Insert(map[string]interface{}{
"name": "张三",
"age":  18,
})
```

> 你甚至可以传递多个map给 insert 方法，依次将多个记录插入到表中：

> 注意：多个map参数要一致，以第一个为准，否则会省略后面不一致的map

```go
// INSERT INTO `user` (`name`,`age`) VALUES(?,?),(?,?) [张三 18 李四 30]
sql, params, err = user.Insert(map[string]interface{}{
"name": "张三",
"age":  18,
}, map[string]interface{}{
"name": "李四",
"age":  30,
})
```

## 更新

```go
// UPDATE `user` SET `name`=?,`age`=? WHERE `id` = ? [test 18 1]
sql, params = user.Table("user").Where("id", 1).Update(map[string]interface{}{
"name": "test",
"age":  18,
})
```

## 删除

```go
// delete from `user` WHERE `id` = ? [1]
sql, params = user.Select("id", "name").Where("id", 1).Delete()

```
