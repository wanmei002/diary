## RPC流
在 RPC 系统中，服务端会实现一组可以远程调用的方法。客户端会生成一个存根，该存根为服务端的
方法提供抽象。这样一来，客户端应用程序可以直接调用存根方法，进而调用服务端应用程序的远程方法。

一般步骤如下: `grpc 使用的是 HTTP/2`
 1. 客户端进程通过生成的存根调用服务端对应的方法。
 2. 客户端存根使用已编码的消息创建 `HTTP POST` 请求。在 gRPC中，所有的请求都是 `HTTP POST` 请求，
 并且 `Content-Type` 前缀为 `application/grpc`。 要调用的远程方法是以单独的 `HTTP`
 头信息的形式发送的。
 3. HTTP请求消息通过网络发送到服务端。
 4. 当接收到消息后，服务器端检查消息头信息，从而确定需要调用的服务方法，然后将消息传递给服务器端
 5. 服务器端将消息字节解析成特定语言的数据结构。
 6. 借助解析后的消息，服务发起对方法的本地调用。
 7. 服务端响应-编码-http/2 返回给客户端
 
 
## gRPC编码方式 protocol buffers 的消息格式
### 消息体
 ```proto
    message GoodsId {
        string value = 1;
    }
 ```

   | 标签 | \[length\]字段值 | 标签 | \[length\]字段值 | ··· |  \[length\]字段值 |
   | ---- | ---- | ---- | ---- | ----|---- |

 - `标签的值` : 由两个值构成 : 字段索引和线路类型。
    1. `字段索引` 就是在 proto 文件中定义消息时，为每个消息字段所设置的唯一数字
    2. `线路类型` 是基于字段类型的，如 `string` 线路类型为`2`, int线路类型是 `0`
 
    value 对应的值是  `value = (field_index << 3) | wire_type` 
    > field_index-字段索引， wire_type-线路类型
    
 - `字段的值` : 
    不同类型的会采用不同的编码类型，string 采用 utf-8, int类型会采用 varint
    
 - `length` : 
    如果编码类型长度是可变的, 会有length值，如 utf-8(string用 utf8), length表示有几个 8 位块
    
> 8位表示标签值，这 8 位中后 3 位为线路类型，线路类型又确定了后面值的编码类型和编码长度等

### 消息头部
前一个字节是压缩标记，后面4个字节是消息长度(意味着 gRPC 通信可以处理大小不超过 `4GB` 的所有消息)

## gRPC 的消息格式
### 客户端发送数据的主要格式
#### 头信息
```header
HEADERS (flags = END_HEADERS)
:method : POST // 请求的方法
:scheme : http // http/https
:path : /Comment/AddComment // 微服务路径 /{服务名}/{服务方法}
:authority : baidu.com // 目标 URL 的虚拟主机名
te : trailers  // gRPC 中 这个值必须为 trailers
grpc-timeout : 15 // 调用的超时时间
Content-Type : application/grpc // grpc 必须以 applaction/grpc 开头
grpc-encoding : gzip // 压缩类型
authorization : ***** 
```

#### 消息体
长度作为前缀以数据帧的方式发送数据，数据过大会分成多帧发送，

#### 消息结束
请求消息结束会在最后一个数据帧上添加 `END_STREAM` 或发送一个带有 `END_STREAM` 空的数据帧。

> END_STREAM 也 称作  EOS(end of stream)

### 响应消息
#### 头信息
```http request
HEADERS (flags = END_HEADERS)
:status : 200
grpc-encoding : gzip
Content-Type: application/grpc
```

#### 消息体
以长度作为消息，会以一个帧或多个帧来发送消息

#### 消息结束
通过发送 trailer 来提醒客户端响应消息已发送。
```http request
HEADERS (flags = END_STREAM, END_HEADERS)
grpc-status : 0 // OK / grcp 状态码
grpc-message : **** // 对错误的描述
```



 
 
## gRPC 支持的通信模式
 1. 一元 RPC 模式
 2. 服务端流 RPC 模式
 3. 客户端流 RPC 模式
 4. 双向流 RPC 模式

> 后续会介绍各个通信模式

 
