package main

// 这段代码是一个使用Go语言编写的服务端程序，它连接到MongoDB数据库，并启动了一个基于HTTP协议的Web服务器。在这段代码中，作者使用了net/http包创建了HTTP服务器，并使用mongo-driver包连接到MongoDB数据库。

// 主要代码解释：
// main() 函数是程序的入口，首先连接到MongoDB数据库，然后启动HTTP服务器。
// 在main()中定义了上下文（context）用于管理MongoDB的连接生命周期。
// connectToMongo() 函数负责创建连接MongoDB的客户端并返回该客户端的指针。
// 程序会监听在指定的端口（webPort），在该端口上启动HTTP服务。

import (
	"context"          // 引入上下文包，用于控制goroutine的生命周期
	"fmt"              // 引入格式化输出包
	"log"              // 引入日志包
	"log-service/data" // 引入本地包log-service/data，用于管理数据库模型
	"net/http"         // 引入HTTP服务器包
	"time"             // 引入时间处理包

	"go.mongodb.org/mongo-driver/mongo"         // 引入MongoDB驱动包
	"go.mongodb.org/mongo-driver/mongo/options" // 引入MongoDB连接选项包
)

// 这段代码导入了多个标准库和外部库，如context、log、http，以及MongoDB相关的mongo和options。
// This section imports several standard libraries and external libraries like context, log, http, and MongoDB-related mongo and options.
// context用于管理上下文，log用于记录日志，http用于处理Web服务器，mongo用于连接MongoDB数据库。
// context is used to manage contexts, log is for logging, http is for handling the web server, and mongo is for connecting to the MongoDB database.

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

// 这段代码定义了几个常量，用于指定Web服务和MongoDB的连接配置。
// This section defines several constants for specifying web service and MongoDB connection configurations.
// webPort 是Web服务器监听的端口，mongoURL 是MongoDB的连接URL。
// webPort is the port the web server listens on, and mongoURL is the URL for connecting to MongoDB.

var client *mongo.Client

// 定义了一个全局变量client，用于存储MongoDB客户端的指针，以便在整个应用程序中使用它。
// A global variable client is defined to store the MongoDB client pointer for use throughout the application.

type Config struct {
	Models data.Models
}

// Config 是一个结构体类型（struct），用于存储数据库模型等配置信息。它包含一个字段 Models，类型是 data.Models，用于管理数据库操作。
// Config is a struct type that stores configuration information such as database models. It has a field Models of type data.Models, used for managing database operations.

// 关于 type 和 struct（About type and struct）:

// type：在Go中，type 关键字用于定义新类型，可以是结构体、接口或别名等。
// In Go, the type keyword is used to define new types, which can be structs, interfaces, or aliases.
// struct：struct 是Go语言中的结构体类型，用于将一组数据字段（属性）组合在一起。
// struct is a type in Go used to group a set of fields (attributes) together.

// 可以改变吗？（Can it be modified?）
// 可以更改 Config 结构体中的字段，如添加、删除或更改字段的类型和名称。但更改后需要修改相应的代码。
// Yes, the fields in the Config struct can be modified, such as adding, deleting, or changing field types and names. However, the corresponding code needs to be updated accordingly.

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	// go app.serve()
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}

}

// 连接到MongoDB（Connect to MongoDB）:
// 首先通过 connectToMongo() 函数连接到MongoDB数据库，并将返回的客户端指针保存到全局变量 client 中。
// First, it connects to the MongoDB database using the connectToMongo() function and saves the returned client pointer to the global variable client.

// 上下文和取消函数（Context and Cancel Function）:

// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 创建一个带有超时时间（15秒）的上下文，用于在超时时自动取消操作。
// Creates a context with a timeout (15 seconds) to automatically cancel the operation if it times out.

// defer cancel()
// defer 用于延迟调用指定函数，cancel() 是用于释放上下文资源的函数。
// defer is used to delay the execution of a function until the surrounding function (here, main()) returns. cancel() is a function used to release context resources.
// 上下文（Context）是什么？（What is Context?）

// context 是Go中的一个包，用于控制 goroutine 的生命周期。它可以传递取消信号、超时和截止时间等信息。
// context is a package in Go used to control the lifecycle of goroutines. It can carry cancellation signals, timeouts, and deadlines.
// 在这段代码中，context.WithTimeout() 用于创建一个带有超时控制的上下文 ctx，防止数据库操作超过指定时间。
// In this code, context.WithTimeout() is used to create a context ctx with a timeout control to prevent database operations from exceeding the specified time.
// Web服务器启动（Starting the Web Server）:
// 代码创建了一个 http.Server 实例，并指定了监听地址和路由处理器，随后使用 srv.ListenAndServe() 启动Web服务器。
// The code creates an http.Server instance, specifies the listening address and route handler, and then starts the web server using srv.ListenAndServe().

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr: fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}

// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic()
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}

// 该函数创建了MongoDB客户端连接选项，并通过 mongo.Connect() 连接到指定的MongoDB服务器（mongoURL），最后返回连接的客户端指针。
// This function creates MongoDB client connection options and connects to the specified MongoDB server (mongoURL) using mongo.Connect(). It then returns the connected client pointer.

// 总结（Summary）：
// type 和 struct 定义了自定义数据类型，并将其用作配置和模型管理，可以根据需要修改 Config 结构体的字段。
// type and struct define custom data types used for configuration and model management. The fields in the Config struct can be modified as needed.

// defer cancel() 用于在 main() 函数结束时自动释放上下文资源。
// defer cancel() is used to automatically release context resources when the main() function ends.

// context 是管理 goroutine 生命周期的机制，支持取消信号、超时和截止时间控制。
// context is a mechanism for managing goroutine lifecycles, supporting cancellation signals, timeouts, and deadline control.
