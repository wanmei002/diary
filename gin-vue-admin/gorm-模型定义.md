### 模型定义
模型是标准的 `struct` , 由 go 的基本数据类型、实现了 `Scanner`、`Valuer` 接口的自定义类型及其指针或别名组成。

```go
type User struct {
    ID uint
    Name string
    Birthday    time.Time
}
```
### 约定
默认情况下 `GORM` 使用 `ID` 作为主键, 使用结构体名的 `蛇形复数` 作为表名, 字段名的 蛇形 作为列名，并使用 `CreateAt` 、`UpdatedAt` 字段追踪创建、更新时间。

#### 表名
GORM 使用结构体名的 `蛇形命名` 作为表名。对于结构体 `User`, 根据约定其表名为 `users`

#### 修改默认表名
```go
func (User) TableName() string {
    return "user"
}
```
#### 动态修改表名
```go
func UserTable(user User) func (tx *gorm.DB) *gorm.DB {
    return func(tx *gorm.DB) *gorm.DB {
        // 根据用户id 取余 来获取表后缀
        prefix := user.ID%100
        return fmt.Sprintf("user_%v", prefix)
    }
}
```
#### 临时指定表名
```go
// 根据 User 的字段创建  deleted_users 表
db.Table("deleted_users").AutoMigrate(&User{})

// 从 tbUser 中查询数据
var u &User
db.Table("tbUser").Find(u)
```


#### 主键
默认 ID 字段为主键, 可以加 `gorm:"primaryKey"` 将其它字段设为主键


### GORM 配置
#### 跳过默认事务
为了确保数据一致性, gorm 会在事务里执行写入操作(创建、更新、删除). 如果没有这方面的要求, 您可以在初始化时禁用它
```go
db, err := gorm.Open(msql.Open(dsn), &gorm.Config{
    SkipDefaultTransaction: true,
})
```
#### 命名策略
GORM 允许用户通过覆盖默认的 `命名策略` 更改默认的命名约定，这需要实现接口`Namer`
```go
type Name interface {
    TableName(table string) string
    ColumnName(table, cloumn string) string
    JoinTableName(table string) string
    RelationshipFKName(Relationship) string
    CheckerName(table, column string) string
    IndexName(table, column string) string
}
```

默认 `NamingStrategy` 也提供了几个选项， 如:
```go
db, err := gorm.Open(nysql.Open(dsn), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
        TablePrefix: "tb",
        SingularTable: true,
    }
})
```

#### PrepareStmt
`PrepareStmt` 在执行任何 SQL 时都会创建一个 Prepared statement 并将其缓存, 以提高后续的效率
```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    PrepareStmt: false,
})
```
#### 取消外键约束
> `DisableForeignKeyConstraintWhenMigrating`

在 `AutoMigrate` 或 `CreateTable` 时, GORM 会自动创建外键约束, 若要禁用该特性, 可将其设置为 `true`

```go
db, err := gorm.Open(mysql.Open(dsn), &grom.Config{
    DisableForeignKeyConstraintWhenMigrating: true,
})
```

### table 命名策略
#### 列名
根据约定，数据表的列名使用的是 `struct` 字段名的 `蛇形命名`
```go
type User struct {
    ID  uint    // 列名是 id
    Name string // 列名是 name
    Birthday time.Time // 列名是 birthday
    CreatedAt time.Time // 列名是 created_at
}
```
可以使用 @@column@@ 标签或 命名策略 来覆盖列名
```go
type Animal struct {
    AnimalID  int64  `gorm:"column:beast_id"`  // 将列设置为 beast_id
    Birthday  time.Time `gorm:"column:day_of_the_beast"` // 将列设置为 day_of_the_beast
}
```

#### 时间戳追踪
##### CreatedAt
对于有 CreatedAt 字段的模型，创建记录时，如果该字段值为零值，则将该字段的值设为当前时间
```go
db.Create(&user)
user2 := User{Name: "zzh", CreatedAt: time.Now()}
db.Create(&user2)
// 想要修改值
```

