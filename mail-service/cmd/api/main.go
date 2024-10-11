package main

import (
	"fmt"      // 引入fmt包，用于格式化字符串输出
	"log"      // 引入log包，用于记录日志信息
	"net/http" // 引入http包，用于启动HTTP服务器
	"os"       // 引入os包，用于读取环境变量
	"strconv"  // 引入strconv包，用于字符串和其他数据类型之间的转换
)

// 这段代码导入了几个标准库，fmt 用于格式化字符串，log 用于记录日志，http 用于启动HTTP服务器，os 用于读取环境变量，strconv 用于字符串和整数之间的转换。
// This code imports several standard libraries. fmt is for formatting strings, log is for logging, http is for starting an HTTP server, os is for reading environment variables, and strconv is for converting between strings and integers.

// 定义了一个配置结构体 Config，其中包含一个 Mailer 字段，用于存储邮件配置。
type Config struct {
	Mailer Mail
}

// 定义了一个 Config 结构体，包含一个 Mailer 字段，该字段是 Mail 类型，用于存储邮件服务的配置。
// Defines a Config struct with a Mailer field, which is of type Mail, and is used to store mail service configurations.

// 定义了一个常量 webPort，用于指定HTTP服务器的端口号
const webPort = "80"

// main 函数是程序的入口，初始化并启动HTTP服务器
func main() {
	// 创建一个 Config 实例，Mailer 字段通过 createMail() 函数来初始化
	app := Config{
		Mailer: createMail(),
	}

	// 记录日志信息，表示服务器启动
	log.Println("Starting mail service on port", webPort)

	// 创建一个 HTTP 服务器实例，并指定监听的端口和处理函数
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // 服务器监听的地址和端口
		Handler: app.routes(),                // HTTP 请求的路由处理函数
	}

	// 	创建了一个HTTP服务器实例 srv，指定了监听地址 Addr（使用 webPort），以及处理路由的处理器 Handler（通过 app.routes() 生成）。
	// Creates an HTTP server instance srv, specifying the listening address Addr (using webPort), and the request handler Handler (generated by app.routes()).

	// 启动HTTP服务器，并监听端口
	err := srv.ListenAndServe()
	if err != nil {
		// 如果服务器启动失败，则记录错误信息并终止程序
		log.Panic(err)
	}
}

// createMail 函数用于创建一个 Mail 实例，并从环境变量中读取配置
func createMail() Mail {
	// 从环境变量中读取 MAIL_PORT 的值，并转换为整数类型
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	// 创建一个 Mail 结构体实例，并从环境变量中读取其他邮件配置
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),     // 邮件服务的域名
		Host:        os.Getenv("MAIL_HOST"),       // 邮件服务器主机地址
		Port:        port,                         // 邮件服务器端口
		Username:    os.Getenv("MAIL_USERNAME"),   // 邮件账户用户名
		Password:    os.Getenv("MAIL_PASSWORD"),   // 邮件账户密码
		Encryption:  os.Getenv("MAIL_ENCRYPTION"), // 邮件加密方式（如TLS、SSL）
		FromName:    os.Getenv("FROM_NAME"),       // 发送者的名称
		FromAddress: os.Getenv("FROM_ADDRESS"),    // 发送者的邮箱地址
	}

	// 返回 Mail 结构体实例
	return m
}
