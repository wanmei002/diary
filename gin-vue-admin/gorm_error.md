### unsupported Scan, storing driver.Value type []uint8 into type *t ime.Time
连接数据库的时候加上 parseTime=true
```go
dsn := fmt.Sprintf(
    "%v:%v@tcp(%v)/%v?charset=utf8&parseTime=true",
    c.User, c.Pwd, c.Host+":"+ strconv.Itoa(c.Port), c.DB,
)
```
