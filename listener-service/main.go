package main

// 这段代码是用 Go 编写的程序，它通过 RabbitMQ 来监听和消费消息队列中的消息。代码的核心逻辑是与 RabbitMQ 建立连接，并创建消费者来处理消息。

// This code is written in Go and interacts with RabbitMQ to listen for and consume messages from the queue. The main logic is to establish a connection with RabbitMQ and create a consumer to handle the messages.

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// github.com/rabbitmq/amqp091-go: RabbitMQ 的 Go 语言客户端，用于与 RabbitMQ 服务器进行通信和操作队列。

// github.com/rabbitmq/amqp091-go: Go client for RabbitMQ used to communicate with the RabbitMQ server and operate on queues.

func main() {
	// 尝试连接到 RabbitMQ
	// Try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1) // 连接失败时，退出程序
	}
	defer rabbitConn.Close() // 在程序结束时关闭 RabbitMQ 连接

	// 开始监听消息
	// Start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// 创建消费者
	// Create a consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err) // 处理创建消费者时的错误
	}

	// 监听队列并消费事件
	// Watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err) // 处理消费消息时的错误
	}
}

// 主函数分析 (Main Function Analysis)
// 连接 RabbitMQ (Connecting to RabbitMQ)：

// 调用 connect() 函数尝试连接 RabbitMQ，如果连接失败，输出错误信息并退出程序。
// Calls the connect() function to try connecting to RabbitMQ. If the connection fails, it logs the error and exits the program.
// 创建消费者 (Creating a Consumer)：

// 调用 event.NewConsumer(rabbitConn) 创建一个消费者对象，用于从指定队列中消费消息。
// Calls event.NewConsumer(rabbitConn) to create a consumer object to consume messages from the specified queues.
// 监听消息队列 (Listening to Message Queue)：

// 调用 consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"}) 来监听队列中的指定类型消息（log.INFO、log.WARNING、log.ERROR）。
// Calls consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"}) to listen for specific message types in the queue (log.INFO, log.WARNING, log.ERROR).

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// 等待 RabbitMQ 准备好再继续连接
	// Don't continue until RabbitMQ is ready

	/*
		amqp://guest:guest@rabbitmq 是用于连接 RabbitMQ 服务器的 AMQP（Advanced Message Queuing Protocol，高级消息队列协议）地址。

		### 链接结构解析 (Connection URL Structure Analysis)
		AMQP 链接的通用格式为：amqp://username:password@host:port/vhost

		在这个链接中：

		- **`amqp://`**:
		- 指定了连接协议类型为 AMQP（通常用于与 RabbitMQ 等消息队列系统通信）。
		- Specifies that the connection protocol is AMQP, commonly used for communication with message queue systems like RabbitMQ.

		- **`guest:guest`**:
		- `guest` 是默认的用户名，后面的 `guest` 是该用户的密码。
		- 默认情况下，RabbitMQ 的用户名和密码都是 `guest`。这是 RabbitMQ 安装后用于本地测试和开发环境的默认配置，但在生产环境中使用时，通常需要更改此默认用户名和密码以确保安全性。
		- `guest` is the default username, and the second `guest` is the password for that user.
		- By default, RabbitMQ uses `guest` as both the username and password. This is a standard configuration for local testing and development environments after installing RabbitMQ. However, for production environments, it is advisable to change this default username and password for security reasons.

		- **`rabbitmq`**:
		- `rabbitmq` 是指 RabbitMQ 服务器的主机地址（hostname）。
		- 在大多数情况下，这个 `rabbitmq` 主机名可以是：
			- **`localhost`**：表示连接到本地机器上运行的 RabbitMQ 实例。
			- **IP 地址**：如 `127.0.0.1` 或者远程服务器的 IP 地址（如 `192.168.1.10`）。
			- **域名**：如 `rabbitmq.example.com`，指向一个指定的 RabbitMQ 服务器。
		- 在 Docker 环境中，通常会使用服务名（如 `rabbitmq`）来作为主机地址，因为容器之间可以使用服务名进行通信，而无需使用 IP 地址。
		- `rabbitmq` is the hostname of the RabbitMQ server.
		- In most cases, this `rabbitmq` hostname can be:
			- **`localhost`**: Connects to a RabbitMQ instance running on the local machine.
			- **IP address**: Such as `127.0.0.1` or the IP address of a remote server (e.g., `192.168.1.10`).
			- **Domain name**: Such as `rabbitmq.example.com`, pointing to a specified RabbitMQ server.
		- In Docker environments, it is common to use the service name (e.g., `rabbitmq`) as the hostname because containers can communicate with each other using service names without needing to know the IP addresses.

		- **省略的 `port` 和 `vhost`（Omitted `port` and `vhost`）**:
		- 在 `amqp://guest:guest@rabbitmq` 这个链接中，没有明确指定端口和虚拟主机（vhost）。
		- 默认情况下，RabbitMQ 使用 **`5672`** 作为 AMQP 的端口。
		- 默认虚拟主机（vhost）为 **`/`**。
		- If `port` and `vhost` are omitted in the `amqp://guest:guest@rabbitmq` URL, RabbitMQ will use the default configurations.
		- The default port for RabbitMQ's AMQP protocol is **`5672`**.
		- The default virtual host (`vhost`) is **`/`**.

		### 实际连接到哪 (Where Does It Actually Connect To)
		- **本地开发环境 (Local Development Environment)**:
		- 如果 `rabbitmq` 主机名指向 `localhost`，则该链接会连接到本地计算机上运行的 RabbitMQ 实例。
		- If the hostname `rabbitmq` resolves to `localhost`, then this URL will connect to a RabbitMQ instance running on the local machine.

		- **Docker 环境 (Docker Environment)**:
		- 在 Docker Compose 或者 Docker Swarm 中，`rabbitmq` 通常是指向 RabbitMQ 容器的服务名。因为 Docker 容器内部可以通过服务名互相访问，因此该链接可能连接到同一网络下的 RabbitMQ 容器实例。
		- In Docker Compose or Docker Swarm setups, `rabbitmq` often refers to the service name of a RabbitMQ container. Since containers within the same Docker network can communicate using service names, this URL likely connects to a RabbitMQ container instance in the same network.

		- **远程服务器 (Remote Server)**:
		- 如果 `rabbitmq` 被解析为一个远程服务器的 IP 地址或域名（如 `rabbitmq.example.com`），那么该链接将连接到远程的 RabbitMQ 实例。
		- If `rabbitmq` resolves to an IP address or domain name of a remote server (e.g., `rabbitmq.example.com`), then this URL will connect to a remote RabbitMQ instance.

		### 代码中的应用场景 (Use Case in the Code)
		在该代码中，`amqp://guest:guest@rabbitmq` 很可能用于连接一个 Docker 环境中的 RabbitMQ 实例。通常在 Docker Compose 中，RabbitMQ 容器的服务名被命名为 `rabbitmq`，因此可以直接通过 `amqp://guest:guest@rabbitmq` 链接到该 RabbitMQ 服务，而无需指定 IP 地址。

		In the code provided, `amqp://guest:guest@rabbitmq` is likely used to connect to a RabbitMQ instance running in a Docker environment. Typically, in a Docker Compose setup, the RabbitMQ container's service is named `rabbitmq`, allowing the code to connect to it directly without specifying an IP address.
	*/

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq") // 尝试连接 RabbitMQ
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...") // 如果连接失败，输出提示信息
			counts++                                 // 记录连接失败的次数
		} else {
			log.Println("Connected to RabbitMQ!") // 连接成功
			connection = c
			break // 跳出循环
		}

		// 如果连接失败次数超过 5 次，返回错误并退出
		// If connection attempts exceed 5 times, return error and exit
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// 计算指数退避时间
		// Calculate exponential backoff time
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...") // 输出退避日志
		time.Sleep(backOff)           // 按照退避时间等待
		continue
	}

	return connection, nil
}

// 连接函数分析 (Connect Function Analysis)
// 定义变量 (Define Variables)：

// var counts int64: 用于记录连接 RabbitMQ 的失败次数。

// var backOff = 1 * time.Second: 初始退避时间为 1 秒。

// var connection *amqp.Connection: 存储成功连接的 RabbitMQ 连接对象。

// var counts int64: Used to count the number of failed connection attempts.

// var backOff = 1 * time.Second: Initial backoff time set to 1 second.

// var connection *amqp.Connection: Stores the successfully established RabbitMQ connection object.

// 循环尝试连接 (Loop to Try Connecting)：

// 使用 for 循环不断尝试与 RabbitMQ 建立连接，直到连接成功或达到最大重试次数。
// Uses a for loop to continuously attempt to establish a connection with RabbitMQ until successful or the maximum retry count is reached.
// 处理连接错误 (Handle Connection Error)：

// 如果连接失败，打印提示信息，并将 counts 增加 1。
// If the connection fails, prints a message and increments counts by 1.
// 退避策略 (Backoff Strategy)：

// 每次连接失败后，使用 math.Pow(float64(counts), 2) 计算指数退避时间（例如，1 次失败等待 1 秒，2 次失败等待 4 秒，3 次失败等待 9 秒，以此类推），然后调用 time.Sleep(backOff) 暂停一段时间。
// After each failed attempt, uses math.Pow(float64(counts), 2) to calculate exponential backoff time (e.g., 1st failure waits 1 second, 2nd failure waits 4 seconds, 3rd failure waits 9 seconds, and so on), and then calls time.Sleep(backOff) to wait for the calculated duration.
// 连接成功与失败处理 (Handling Success and Failure)：

// 如果连接成功，打印连接成功信息，退出循环，并返回连接对象。
// 如果连接失败次数超过 5 次，打印错误信息并返回 nil。
// If connection is successful, logs the success message, exits the loop, and returns the connection object.
// If failed attempts exceed 5 times, logs the error message and returns nil.
