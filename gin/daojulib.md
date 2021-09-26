#### 看代码对工作流程理解
 - djframework/app 文件中 Router,输入 appName(微服务名) path(url上对应的路径？瞎猜的) controller(微服务对应的控制器) mappingMethods(入口方法相当于PHP中的StartApp(可能不对))
 - 开始走控制器 初始化 BeginData 结构体 -\> BeginExec -\> 创建 runtime实例(这个很重要, 框架中很多方法用到它) -\> 写你自己的逻辑吧 
    + djbase/log/djlog 中有日志记录方法 djreport包中有上报方法
    +  djrpc 包中有 curl dbproxy(mysql)
 - 运行 app包中的 Run 方法 让应用跑起来
#### attr_api 感觉此包像是 redis 简化版 功能应该是上报之类的
 - Add(int, int)  对一个int 值进行累加行为 像是计数器
 - Set(int, int)  设置一个值，此键已经被设置则覆盖
 - Warning(int, string, ...) 向网关发送警告之类的信息

#### djapi/dipaas AMS paas 接口
 - Request 指定 L5 请求接口 并解析接口数据

#### djapp 道聚城订单接口(实际上只是读取了配置)
 - Instance 这个方法是主要方法 读取 ini 配置文件 并缓存 2min
    + 采用原子操作 一次性不中断读取配置文件
     	+ 引用 beego/config 包读取 ini 文件 数据保存在结构体 data 属性中
     		+ config 包采用接口 和 注册的方法 注册了 ini 配置解析的方法，实现 Configer 接口，可以扩展 config 解析的其他文件的功能
    + 用 config中 Int String 方法从 config.data 中取出相应的数据
     	+ 读 beego.config 源码 可知 `::` 分隔符读取二级map数据

#### djbase/cl5/djl5 读取 l5 配置文件
 - readL5Cfg 读取 l5 配置文件 缓存 2 min
 - GetL5(node string) 返回 mod cmd 节点的默认ip port

#### djbase/db/djdb 读取数据库节点配置文件的信息 并 返回传入节点的 info 结构体 里面有属性 ip port proxyip proxyport
 - Get(node string) 主要方法 还有读取配置文件方法 跟上面的逻辑一样 就不做过多介绍了

#### djbase/log/djlog 日志输出
 - 引用了 beego 框架 的 log 系统，这个主要往 os.Stdout 里输出日志信息

#### djcbase/report/amshelper 应该是上报
 - 把信息组装写入 buffer 缓存中，然后通过 udp 协议 发送给指定的服务器


#### djbase/uin/djuin uin 编解码
 - 通过位移 或 对uin 编解码

#### djbiz 读取相应 biz 的配置文件
 - Instance 主要方法，读取 跟 djapp逻辑 几乎一样
 - 返回自定的一个结构体 `DjBizInfo` , 这个结构体有很多获取主要信息的方法

#### djbiz/djyxzj 通过 curl 获取相应的用户 英雄 皮肤等信息 ，主要调用 paas 包里的方法

#### djbucket 取公共的令牌？

#### djcjf 查询云积分接口
 - RequestByMap 这个应该是入口方法 传入 mod cmd 和 url 参数 或 post 参数，通过 CalcParams, 重组 url 参数，生成 sign 查询

#### djdb 合成 insert update 语句
 - MysqlEscape 把值中特殊字符转义

#### djruntime 包里有些结构体 都是比较重要的
##### 主要结构体
 - RunContext :: host pid Ns(启动时间戳) HostName
 - RpcInfo :: 保存了 RPC 调用时的状态 和 返回信息(请原谅我的片面理解)
 - FmtInfo :: 保存输出的主要字段(....再次的)
 - InitInfo :: 初始状态信息(应该是这意思吧)
 - Runtime :: 看结构里面有 rpclist 属性 应该是 rpc 上级的状态信息吧
    + FmtInfo 把 Runtime 信息 json 格式化 抛出
 
##### 主要的函数
 - init 初始化函数 初始化 RunContext 结构体相关 IP PID 等属性 
    + rand.Seed 每次运行保证随机数不一样
 - NewEventId 创建eventId
 - NewRuntime 返回RunTime 结构类型 并初始化
 - NewRpcInfo 初始化RpcInfo 结构体 并往父级 Runtime::rpclist 中存入它
 - 剩下的应该是 ams 日志相关的函数
 
##### 流程推断 (如有驴嘴不对马脸的情况 请大佬轻蔑一笑 不要当真)
 1. init 初始化当前进程相关信息
 2. NewRuntime 初始化当前服务的信息
 3. NewRpcInfo 保存请求的其它微服务信息   
    
#### djframework 看字面意思应该是框架的信息 让我简单拜读下大佬的代码
##### 先从 stat 文件看起
 - 结构体
    + SvrStatInfo 接收 成功 失败 相关信息
    + statInfo 批量保存服务状态 顾客状态
    + statData 当前状态信息 : 错误 时间 Runtime
 - 方法
    + initStat 初始化管道 管道里保存 statData 结构体
    + doStat 如果 statData 不为空 服务; 为空 客户端服务
        - doRunTimeStat 新来服务 更新全局 接收数量和接收耗时 ; 更新成功数量 或 失败数量 和耗时
        - doCustomerStat 更新全局客户端相关的状态
    + AddRuntimeStat 添加当前进程到管道中(可能用管道做队列吧)
##### app 文件
 - 用 beego 框架的Router 注册路由 和 启动服务
##### controller 文件这个应该是特别重要的文件了
 - 结构体
    + BeginData:: AppName Biz AppId PlugId ActId EventId NotChkSign(0-强制校验 1-走网关时校验)
    + ControllerInterface 控制器接口
    + Controller 应该是父级控制器 实现 MC 结构, 属性主要有 : _appName _bizcode _loginUid _loginType _loginToken _rt(Runtime 结构体实例)等等
        - SetOutputHdr 设置输出头信息
        - DelOutputHdr 删除头信息
        - GetCookie 获取指定的 cookie 值
        - GetCookies 获取全部 cookie 串成字符串, 这里有个问题 为什么串联字符串 用 fmt.Sprint(v)
        - BeginExec 主要是是指 Controller 结构体中的属性值
            + biz appid plugId eventId NotChkSign 
            + 开始校验基本参数 bizcode appid actid plugId
            + 记录当前请求到日志 
            + 创建 Runtime 实例 传入 appname parent eventid 端口等信息
            + 是否 notChkSign 设置为1 并且不是网关 就不会校验
        - Runtime 返回 runtime 结构体
        - SetRetKey 设置是 result或ret 状态字段
        - SetTsKey 设置 输出的时间字段是哪个
        - *EscapeHtml 是否转码输出的 html 字段
        - SetEventId 请求其它微服务时可以传此字段 使日志连续 可以看到调用的微服务日志信息
        - SetAppName 
        - SetBiz
        - SetLoginUin GetLoginUin SetLoginUid GetLoginUid GetLoginType 
        - AddRptDetail
        - DoInitUid 
        - AddRptExt 增加日志输出？
        - Report
        - OutputAndStop 这个慎用吧？
        - OutputR OutputDS 输出的 data 是 map\[string\]\[string\]
        - OutputDX 输出的 data 是 map\[string\]interface{}
        - GetOutputString GetOutputStr
        - Output DoOutput OutputS 输出头信息和 body

##### framework 
 - Init 函数 初始化 日志 心跳等信息
##### heart
 - initHeart 1 分钟重载一下配置文件
##### login 校验登录态

#### djja/gjapool
##### dbproxy 
 - JsGoDb 注册方法？
 - InitDbConn 链接节点
 - JsGoDbQuery JsGoDbUpdate JsGoDbTrans 查询 更新 事物
##### redis
 - redis 操作
##### util 
#### djlogic/ams/djams/openid2opneid

#### djlol ApiGetLolMasterUid
#### djmemcache
 - memcache 相关操作
#### djpay
 - Intance 主要函数，解析配置文件并返回含有 beego.Config 的结构体
 
#### djrpc
##### curl
 - 两个结构体 CurlReq CurlRsp
 - Curl 函数 做的有点简单吧感觉
 	+ Client 设置超时时间
 	+ 创建并且设置请求体
 	+Client.Do 执行

##### curl_l5
 - 两个结构体 CurlReqL5 CurlRspL5  这两个结构体分别保存对应的 curl 里的两个结构体
 - CurlL5 根据 mod cmd 获取 host ，根据 CurlReqL5 里的 api QueryS 请求 url 数据保存到 CurlRspL5

##### dbproxy MysqlDbProxy 实现了 DbInterface 接口
 - NewMysqlDbproxy 创建 mysql 连接
 - Insert Query Update Trans Close 等方法

##### http 
 - 主要是 HttpCall  实际上还是调用了 CurlL5 方法

##### mysql 也实现了 DbInterface 接口
 - NewMysqlPool 创建 mysql 连接
 - Query Insert Update Trans Close

##### redis
 - NewRedisPool 创建 Redis 连接 

##### udp
 = RPC 是基于 tcp 协议的

##### djtools/djcurl
 - RequestByL5 l5 curl 请求
 - RequestByDomain 域名请求

##### djutil 工具包
 - 常用的函数 就不介绍了
 

    
    

 
 