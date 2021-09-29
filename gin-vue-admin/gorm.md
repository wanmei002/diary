```go
import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
)

func init() {
  db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

  // 大部分 CRUD API 都是兼容的
  db.AutoMigrate(&Product{})// 生成 products table
  // 添加数据
  db.Create(&user)
  // 找到 id=1 的数据
  db.First(&user, 1)
  // 更新 age=18
  db.Model(&user).Update("Age", 18)
  // 
  db.Model(&user).Omit("Role").Updates(map[string]interface{}{"Name": "jinzhu", "Role": "admin"})
 
  db.Delete(&user)
}
```
#### 批量插入
```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)
// 插入完后 自动填充 users 里的id数据
for _, user := range users {
  user.ID // 1,2,3
}
```
#### 指定创建的数量
```go
var users = []User{name: "jinzhu_1"}, ...., {Name: "jinzhu_10000"}}
//指定创建的数据为  100  条
db.CreateInBatches(&users, 100)
```

### 预编译模式
> 预编译模式会预编译 sql 语句, 以加速后续执行速度
```go
// 全局模式，所有的操作都会创建并缓存预编译语句，以加速后续执行速度
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{PrepareStmt: true})

// 会话模式，当前会话中的操作会创建并缓存预编译语句
tx := db.Session(&Session{PrepareStmt: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```
### DryRun 模式
> DryRun 模式会生成但不执行 SQL, 可以用于检查、测试生成的sql
```go
stmt := db.Session(&Session{DryRun: true}).Find(&user, 1).Statement
stmt.SQL.String()  // mysql 会生成 SELECT * FROM `users` WHERE `id` = ?
stmt.Vars          // []interface{}{1}
```
### find to map
> scan 结果到 `map[string]interface{}` 或 `[]map[string]interface{}`
```go
// TODO map 必需分配内存
result := map[string]interface{}{}
db.Model(&User{}).First(&result, "id = ?", 1)
```
### create from map
> 根据 `map[string]interface{}` 或 `[]map[string]interface{}` create data
```go
db.Model(&User{}).Create(map[string]interface{}{"Name": "zzh", "Age": 18})

dataList := []map[string]interface{}{
    {"Name": "zzh1", "Age": 17},
    {"Name": "zzh2", "Age": 16},
}
db.Model(&User{}).Create(dataList)
```

### FindInBatches
> 用于批量查询并处理记录
```go
result := db.Where("age>?", 13).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
    if btch == 1 {
        log.Println(第一次查询) // limit 0 ,100
    }
    return nil
})
```
### 嵌套事务
```go
user1 := &User{}
    DB.Transaction(func(tx *gorm.DB) error {
        res1 := tx.Table(UserTable).Select("Name", "PassWord", "Phone", "Salt").Create(user1)
        log.Printf("res1:%+v; user:%+v\n", res1, user1)
    
        user2 := &User{Name: "zzh2",PassWord: "123457",Phone: "123456799",Salt: "abc2"}
    
        tx.Transaction(func(tx2 *gorm.DB) error {
            // 用方法传进来的 tx2
            res2 := tx2.Table(UserTable).Select("Name", "PassWord", "Phone", "Salt").Create(user2)
            log.Printf("res2:%+v; user:%+v\n", res2, user2)
            // 返回的 error!=nil 就会回滚
            return errors.New("i want rollback")
        })
        user3 := &User{}
        tx.Transaction(func(tx3 *gorm.DB) error {
        // TODO 用方法传进来的 tx3
            res := tx3.Table(UserTable).Select("Name", "PassWord", "Phone", "Salt").Create(user3)
            log.Printf("res3:%+v; user:%+v\n", res, user3)
            return nil
        })
        
        return nil
    })
```

### SavePoint, RollbackTo
```go
tx := db.Begin()
tx.Create(&user1)
tx.SavePoint("sp1")

tx.Create(&user2)
tx.RollbackTo("sp1")// rollback user2

tx.Commit()
```

### 命名参数
```go
db.Where("name1 = @name OR name2=@name", sql.Named("name", "zzh")).Find(&user)
// SELECT * FORM `users` WHERE name1="zzh" OR name2="zzh"

db.Where("name1=@name OR name2=@name", map[string]interface{}{"name": "zzh"}).First(&user)
// SELECT * FROM `users` WHERE name1="zzh" OR name2="zzh" LIMIT 1

db.Row("SELECT * FROM `users` WHERE name1=@name OR name2=@name", map[string]interface{}{"name": "zzh"}).Find(&user)// &user 可以用 map 替换，但是 map 必须 make 分配内存

db.Exec("UPDATE users SET name1=@name, name2=@name2",
    map[string]interface{}{"name": "zzh1", "name2": "zzh2"}
)
// UPDATE users SET name="zzh1", name2="zzh2"
```

### 分组条件
```go
db.Where(
    db.Where("pizza = ?", "pepperoni").Where(db.Where("size=?", "small").Or("size=?", "medium")),
).Or(
    db.Where("pizza = ?", "hawaiian").Where("size=?", "xlarge"),
).Find(&pizzas)
```

### 支持多个字段追踪 create/update 时间(time、unix(毫/纳)秒)
```go
type User struct {
    CreateedAt time.Time    // 在创建时, 如果该字段值为零值，则使用当前时间填充
    UpdatedAt   int         // 在创建时该字段为零值或者在更新时, 使用当前时间戳的秒数填充
    Updated     int64 `gorm:"autoUpdateTime:nano"`       // 使用时间戳的秒数填充更新时间
    Updated2    int64 `gorm:"autoUpdateTime:milli"`     // 使用时间戳的毫秒数填充更新时间
    Created     int64   `` 
}
```

