## customizeLogger

一个线程安全的支持异步写入分区日志文件的通用日志框架

### Features

* 支持分区写入日志文件

* 自动拆分日志文件（按天拆分）

* 线程安全，异步写入
* WriterPool写入池保证了写入不同日志文件的效率

### 可配置环境变量
- LOG_PATH：日志文件保存路径，默认当前目录
- LOG_LEVEL：日志级别，默认info

## Get Start

#### 控制台输出
初始化一个日志对象
```go
log := CustomizeLogger{}
log.Init()
```

使用

```go
log.Info("info....")
log.Error("error....")
```
输出
```
{"level":"info","message":"info....","time":"2024-06-24T20:41:03.095674+08:00"}
{"level":"error","message":"error....","time":"2024-06-24T20:41:03.095913+08:00"}
```
#### 控制台与日志文件共同输出

初始化一个带写日志文件的对象
```go
log := CustomizeLogger{}
log.SetOutPutNew(true)
log.Init()
```

使用
```go
log.Info("info....")
log.Error("error....")
```
控制台输出
```
{"level":"info","message":"info....","time":"2024-06-24T20:41:03.095674+08:00"}
{"level":"error","message":"error....","time":"2024-06-24T20:41:03.095913+08:00"}
```
日志文件输出：
```
.
├── README.md
├── dataStructure
│   ├── map
│   │   ├── concurrent.go
│   │   └── concurrent_test.go
│   └── pool
│       └── writer.go
├── go.mod
├── go.sum
├── log
│   ├── hooks.go
│   ├── log.go
│   ├── log_test.go
│   ├── structs.go
│   └── utils.go
├── logs
│   └── 2024-06-24
│       └── log.log
└── main.go
```
当前目录下将自动生成`logs/yyyy-mm-dd/log.log`日志文件

指定日志分区
```go
log.WithFieldUID("log-1").Info("info...")
log.WithFieldUID("log-1").Error("info...")
```
在控制台输出的同时，会在`logs/yyyy-mm-dd/`目录下生成`log-1.log`日志文件

同样支持向不同的分区写入日志而不用考虑性能和线程安全的问题，借助`WriterPool`对象，实现了文件writer流的复用

```go
log.WithFieldUID("log-partition-1").Info("info...")
log.WithFieldUID("log-partition-2").Error("info...")
```
在控制台输出的同时，会在`logs/yyyy-mm-dd/`目录下分别生成`log-partition-1`和`log-partition-2`两个日志文件


#### 关闭日志：

若是未使用输出文件功能，则无需关闭操作。

以下针对使用日志文件输出做关闭操作说明：
```go
// 初始化
log := CustomizeLogger{}
log.SetOutPutNew(true)
log.Init()

// 关闭
log.Cleanup() // 清理写入池所有writer写入流
log.Wg.Wait() // 等待异步写入完成
```

#### 服务中使用（以GRPC服务为例）
```go
const (
	Address string = ":50051"
	Network string = "tcp"
)


func main() {
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	
	......
	
	customizeLogger := clog.CustomizeLogger{}
	customizeLogger.SetOutPutNew(true)
	customizeLogger.Init()

	// 通过调用定时任务动态清理过期的writer流对象保证了整体的性能
	c := cron.New()
	_, err = c.AddFunc("*/1 * * * *", func() {
		customizeLogger.ServiceHookNew.WriterPool.CleanExpiredWriter()
	})
	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return
	}
	c.Start()

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
```

