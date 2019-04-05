## 服务注册与启动
hydra已支持6种服务类型:`http api`服务，`rpc`服务，`websocket`,`mqc`消息消费服务，`cron`定时任务,`web`服务. 

### 一. 服务注册

`hydra`实例提供了8个函数进行服务注册, 不同的函数可注册到不同的服务器,见下表:

| 注册函数 | api | rpc | web | ws  | mqc | cron |
| -------- | --- | --- | --- | --- | --- | ---- |
| Micro    | √   | √   | √   | ×   | ×   | ×    |
| Flow     | ×   | ×   | ×   | ×   | √   | √    |
| API      | √   | ×   | ×   | ×   | ×   | ×    |
| RPC      | ×   | √   | ×   | ×   | ×   | ×    |
| WEB      | ×   | ×   | √   | ×   | ×   | ×    |
| WS       | ×   | ×   | ×   | √   | ×   | ×    |
| MQC      | ×   | ×   | ×   | ×   | √   | ×    |
| CRON     | ×   | ×   | ×   | ×   | ×   | √    |

代码示例:
```go
    app.API("/hello",hello)
    app.MQC("/hello",hello)

    func hello(ctx *context.Context) (r interface{}) {
	    return "hello world"
    }   
```

注册函数支持两种类型:
   * 1. 函数注册: 服务实现代码放在函数中,函数签名格式为:`(*context.Context) (interface{})`,示例:
```go
        func hello(ctx *context.Context) (r interface{}) {
            return "hello world"
    }
```
   * 2. 实例注册: 服务实现代码放到`struct`中,传入`struct`实例的构造函数
  
        示例:
  ```go
             app.API("/hello",token.NewQueryHandler)
```

        添加服务实现文件`query.handler.go`

```go

            package token

            import (
                "github.com/micro-plat/hydra/component"
                "github.com/micro-plat/hydra/context"
            )

            type QueryHandler struct {
                container component.IContainer
            }


            //NewQueryHandler 创建服务
            func NewQueryHandler(container component.IContainer) (u *QueryHandler) {
                return &QueryHandler{
                    container: container,
                }
            }
            func (u *QueryHandler) Handle(ctx *context.Context) (r interface{}) {
                var result struct {
                    ErrCode int64  `json:"errcode"`
                    ErrMsg  string `json:"errmsg"`
                }
                result.ErrCode = 0
                result.ErrMsg = "success"
                return result
            }
  ```
该`struct`需具备两个条件:

1. 服务构造函数`NewQueryHandler`,只能是两种格式之一:
   `(container component.IContainer) (*QueryHandler) ` 或
   `(container component.IContainer) (*QueryHandler,error) `

2. 对象中至少包含一个命名为`...Handle`的函数,且签名为:
   `(*context.Context) (interface{})`格式


### 二. 服务启动
一个应用程序实例可启动6种服务器的任意组合,只需使用`-`连接,可通过代码或命令行指定:

#### 1. 代码中指定
 * 启动`api`,`rpc`服务器实例
  ```go
     hydra.WithServerTypes("api-rpc"),
  ```


 * 启动`api`,`cron`,`mqc`服务器实例
  ```go
     hydra.WithServerTypes("api-rpc-mqc"),
  ```

#### 2. 命令行中指定
启动`api`和`rpc`实例
```sh
$ sudo ./helloserver run -r fs://../ -c test -S "api-rpc"
```

#### 3. 服务启动
   可使用命令`run`和`start`启动服务,区别是:

> `run` 直接运行服务. 所有日志输出到控制台, 并根据级别显示不同颜色,便于调试,一般开发时使用此命令

> `start` 服务安装后可使用`start`命令启动, 服务将在在后台运行, 异常关闭或服务器重启会自动启动应用. 日志存入日志文件或远程日志归集系统, 控制台不显示日志. 可使用`stop`停止服务,`status`查看服务是否运行,`remove`卸载服务.

* 修改任何配置,请重新执行`install`命令
 > 执行`install`时返回`Service has already been installed`错误,则需执行`remove`命令

 ```sh
$ sudo ./helloserver install -r fs://../ -c test -S "api-rpc"
Service has already been installed

$ sudo ./helloserver remove
Removing helloserver:					[  OK  ]
 ```

再次执行`install`命令

```sh
$ sudo ./helloserver install -r fs://../ -c test -S "api-rpc"
	-> 创建注册中心配置数据?如存在则不安装(1),如果存在则覆盖(2),删除所有配置并重建(3),退出(n|no):2
		修改配置: /myplat_debug/helloserver/api/test/conf
		创建配置: /myplat_debug/helloserver/rpc/test/conf
Install helloserver:					[  OK  ]
```
* `run`启动服务
  
  一般`run`命令参数与`install`一致(`start`时不需要任何参数)
  
```sh
$ sudo ./helloserver run -r fs://../ -c test -S "api-rpc"

[2019/04/05 19:12:01.339456][i][7ab175312]初始化 /myplat_debug/helloserver/api-rpc/test
[2019/04/05 19:12:01.341163][i][8edb58733]开始启动...
[2019/04/05 19:12:01.341623][i][8edb58733][启用 静态文件]
[2019/04/05 19:12:01.341635][d][8edb58733][未启用 header设置]
[2019/04/05 19:12:01.341648][d][8edb58733][未启用 熔断设置]
[2019/04/05 19:12:01.341655][d][8edb58733][未启用 jwt设置]
[2019/04/05 19:12:01.341658][d][8edb58733][未启用 ajax请求限制设置]
[2019/04/05 19:12:01.341661][d][8edb58733][未启用 metric设置]
[2019/04/05 19:12:01.341664][d][8edb58733][未启用 host设置]
[2019/04/05 19:12:01.843111][i][8edb58733]启动成功(http://192.168.1.8:8090,1)
[2019/04/05 19:12:01.843367][i][7650a8ecf]开始启动...
[2019/04/05 19:12:01.843827][d][7650a8ecf][未启用 jwt设置]
[2019/04/05 19:12:01.843841][d][7650a8ecf][未启用 header设置]
[2019/04/05 19:12:01.843846][d][7650a8ecf][未启用 metric设置]
[2019/04/05 19:12:01.843849][d][7650a8ecf][未启用 host设置]
[2019/04/05 19:12:02.345952][i][7650a8ecf]启动成功(tcp://192.168.1.8:8081,1)
  ```
控制台打印出了两次`启动成功`,分别是`api`服务器(http协议),`rpc`服务器(甚于grpc,tcp协议),包含服务提供地址和启动的服务个数

同一个服务器的日志可根据`session_id`(当前启动实例为:`8edb58733`,`7650a8ecf`)查看上下文日志

* `start`启动服务
```sh
$ sudo ./helloserver start
Starting helloserver:					[  OK  ]
```
控制台只会输出启动成功,不会显示运行时日志



#### 4. 服务状态查询

* 服务发布信息
  
  服务器启动时会自动将当前服务器和服务添加到注册中心, 便于监控服务运行状况和服务发现者查找服务.

* 服务器监控节点

 `api`服务: `/myplat_debug/helloserver/api/test/servers/192.168.1.8:8090`
`rpc`服务:`/myplat_debug/helloserver/rpc/test/servers/192.168.1.8:8081`



* 服务提供者节点

`api`服务: `/myplat_debug/services/api/helloserver/hello/providers/192.168.1.8:8090`

`rpc`服务: `/myplat_debug/services/rpc/helloserver/hello/providers/192.168.1.8:8081`

如使用`fs://../`指定的注册中心,则运行以下命令查看:
```sh
$ cd ../myplat_debug/helloserver/api/test/servers/

$ ls
192.168.1.8:8090

$ cat 192.168.1.8\:8090 
{"service":"http://192.168.1.8:8090"}
```

其它节点内容请自行查询