#### protobuf3 语法学习 基于golang

> 借鉴于 https://www.cnblogs.com/tohxyblog/p/8974763.html

##### 基本的运行原理
 - 首先安装 [protoc 编译器]( https://github.com/google/protobuf/releases)
 - 安装生成go代码的插件 `github.com/golang/protobuf/protoc-gen-go` 此处会生成一个可执行的程序 `protoc-gen-go` , 需要把 $GOPATH/bin 目录加入环境变量中
 - 运行命令生成go代码 : `protoc --go_out=. hello.proto` hello.proto 是protobuf 文件
 > protobuf 的 protoc 编译器是通过插件机制实现对不同语言的支持。比如 protoc 命令出现 --xxx_out 格式的参数, 那么 protoc 将首先查询是否有内置的 xxx 插件, 如果没有
 内置的 xxx 插件那么将继续查询当前系统中是否存在 protoc-gen-xxx 命名的可执行程序, 最终通过查询到的插件生成代码。 对于 go 语言的 protoc-gen-go 插件来说，里面又实现了一层静态插件系统。比如 protoc-gen-go 
 内置了一个 gRPC 插件，用户可以通过 --go_out=plugins=grpc 参数来生成 gRPC 相关代码, 否则只会针对 message 生成相关代码

##### 第一个例子
```proto
syntax = "proto3";

message SearchRequest {
    string query = 1;
    int32 page_number = 2;
    int32 result_per_page = 3;
}
``` 
- 对上面的例子进行讲解分析
    + 第一行 `syntax = "proto3";`, 这一行指定了正在使用 `proto3` 语法(指定语法行必须是文件的非空非注释的第一行). 如果你没有指定这个, 编译器会使用 `proto2` 语法解析文件
    + SearchRequest 消息格式有3个字段, 在消息中承载的数据分别对应于每一个字段。其中每个字段都有一个名字和一种类型
    + 每个字段都有唯一的一个数字标识符，这些标识符是用来在消息的二进制格式中识别各个字段的，一旦开始使用就不能再改变, [1,15]之内的标识符再编码的时候会占用一个字节。
    [16, 2047] 之内的标识符则占用2个字节, 所以应该为那些频繁出现的消息元素保留[1, 15]之内的标识符
    + 最小的标识符可以从 1 开始, 最大到 2^29-1, or 536,870,911. 不可以使用其中的[19000-19999]
    
- 指定字段规则, 指定的消息字段修饰符必须是如下之一:
    + `singular`: 一个格式良好的消息应该有 `0`个 或 `1`个这种字段(不能超过`1`个).
    + `repeated`: 在一个格式良好的消息中, 这种字段可以重复任意多次(包括`0`次). 重复的值的顺序会被保留
    在`proto3`中, repeated 的标量域默认情况下使用 `packed`
 
 
- 枚举
    ```proto
    message SearchRequest {
        string query = 1;
        int32 page_number = 2;
        int32 result_per_page = 3;
        enum Corpus {
            UNIVERSAL = 0;
            WEB = 1;
            IMAGES = 2;
            LOCAL = 3;
            NEWS = 4;
            
        }
        Corpus corpus = 4;
    }
    ```
    + 每个枚举必须将其第一个类型映射为0， 这是因为:
        + 必须有一个 0 值, 我们可以用这个 0 值作为默认值
        + 这个零值必须为第一个元素, 为了兼容 `proto2` 语法, 枚举类的第一个值总是默认值
        + 你可以通过将不同的枚举常量指定相同的值。如果这样做你需要将 `allow_alias` 设定为 `true`, 否则编译器会在别的地方产生一个错误信息
        ```proto3
        enum EnumAllowingAlias {
          option allow_alias = true;
          UNKNOW = 0;
          STARTED = 1;
          RUNNING = 1;
        }
        ```
- 使用其他消息类型
    ```proto3
    message SearchResponse {
        repeated Result results = 1;
    }
    
    message Result {
        string url = 1;
        string title = 2;
        repeated string snippets = 3;
    }
    ```
    
- 嵌套类型
    ```proto3
    message SearchResponse {
        message Result {
            string url = 1;
            string title = 2;
            repeated string snippets = 3;
        }
        repeated Result results = 1;
    }
    
    // 如果你想在它的父消息类型的外部重用这个消息类型, 你需要以 Parent.Type 的形式使用它，如:
    message SomeOtherMessage {
        SearchResponse.Result result = 1;
    }
    ```
    
- 在一个 .proto 文件中可以定义多个消息类型
- 向 .proto 文件添加注释, 可以使用 `//` 语法
- 保留标识符(Reserved) , protocol buffer 的编译器会警告未来尝试使用这些域标识符的用户
    ```proto
    message Foo {
  	    reserved 2,15,9 to 11;
  	    reserved "foo", "bar";
    }
    ```
    
    > 不要在同一行 reserved 声明中同时声明域名字和标识号
    
- 对 go 来说, 编译器会为每个消息类型生成一个 .pd.go 文件

- 以下是类型对照表
    ```
    double => float64
    float  => float32
    int32 => int32
    uint32 => uint32
    uint64 => uint64
    sint32 => int32  // 变长编码, 这种编码在负值时比 int32 高效的多
    sint64 => int64
    fixed32 => uint32
    fixed64 => uint64
    sfixed32 => int32
    sfixed64 => int64
    bool => bool
    string => string
    bytes => []byte
    ```
    
- 

