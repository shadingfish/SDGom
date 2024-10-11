package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
Consumer 结构体定义了一个 RabbitMQ 消费者，包含连接和队列名称属性。

### 字段说明 (Field Description)
- `conn *amqp.Connection`：RabbitMQ 的连接对象，用于在消费者和 RabbitMQ 服务器之间建立连接。
- `queueName string`：队列名称，用于指定消费者监听的队列。

### Struct Description
The `Consumer` struct defines a RabbitMQ consumer with connection and queue name properties.

- `conn *amqp.Connection`: RabbitMQ connection object used to establish communication between the consumer and the RabbitMQ server.
- `queueName string`: The name of the queue that the consumer is listening to.
*/

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

/*
NewConsumer 函数用于创建一个新的 `Consumer` 实例，并设置交换机。

### 函数参数 (Function Parameters)
- `conn *amqp.Connection`：RabbitMQ 连接对象，用于在消费者和 RabbitMQ 服务器之间建立连接。
  - `conn *amqp.Connection`: RabbitMQ connection object used to establish communication between the consumer and the RabbitMQ server.

### 返回值 (Return Value)
- `Consumer`：新创建的 `Consumer` 实例。
  - `Consumer`: Newly created `Consumer` instance.
- `error`：设置交换机时发生的任何错误。
  - `error`: Any error that occurs during exchange setup.

### 函数描述 (Function Description)
该函数初始化 `Consumer` 结构体，并调用 `setup` 函数来创建和配置 RabbitMQ 交换机。返回配置好的消费者对象或错误信息。
This function initializes the `Consumer` struct and calls the `setup` function to create and configure the RabbitMQ exchange. It returns the configured consumer object or an error.
*/

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	// 调用 setup 函数进行交换机设置 (Call setup function for exchange setup)
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

/*
setup 函数用于配置 RabbitMQ 交换机。

### 返回值 (Return Value)
- `error`：在设置交换机时发生的任何错误。
  - `error`: Any error that occurs during exchange setup.

### 函数描述 (Function Description)
该函数创建 RabbitMQ 通道，并使用 `declareExchange` 函数在通道上声明一个交换机。
This function creates a RabbitMQ channel and uses the `declareExchange` function to declare an exchange on the channel.
*/

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

/*
Payload 结构体定义了消息的载荷格式。

### 字段说明 (Field Description)
- `Name string`：消息名称，用于标识消息的类型。
- `Data string`：消息内容，存储具体的消息数据。

### Struct Description
The `Payload` struct defines the format of the message payload.

- `Name string`: The name of the message used to identify the type of message.
- `Data string`: The content of the message storing the actual data.
*/
type Payload struct {
	Name string `json:"name"` // 消息名称 (Message name)
	Data string `json:"data"` // 消息内容 (Message content)
}

/*
Listen 函数用于监听指定主题的消息队列。

### 函数参数 (Function Parameters)
- `topics []string`：要监听的消息主题列表。
  - `topics []string`: List of message topics to listen to.

### 返回值 (Return Value)
- `error`：监听或消费消息时发生的任何错误。
  - `error`: Any error that occurs during message listening or consumption.

### 函数描述 (Function Description)
该函数使用消费者的 RabbitMQ 连接创建通道，并在通道上声明一个随机队列。随后将指定的主题绑定到该队列，并开始消费消息。每条消息都被解析为 `Payload` 类型，并传递给 `handlePayload` 函数进行处理。

This function creates a channel using the consumer's RabbitMQ connection and declares a random queue on that channel. It then binds the specified topics to the queue and starts consuming messages. Each message is unmarshalled into a `Payload` type and passed to the `handlePayload` function for processing.
*/

func (consumer *Consumer) Listen(topics []string) error {
	// 创建 RabbitMQ 通道 (Create a RabbitMQ channel)
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close() // 函数结束时关闭通道 (Close the channel when function ends)

	// 声明一个随机名称的队列 (Declare a random queue)
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// 将指定的主题绑定到队列 (Bind specified topics to the queue)
	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,       // 队列名称 (Queue name)
			s,            // 路由键 (Routing key)
			"logs_topic", // 交换机名称 (Exchange name)
			false,        // 是否等待服务器响应 (No-wait?)
			nil,          // 额外参数 (Arguments)
		)

		if err != nil {
			return err
		}
	}

	// 开始消费消息 (Start consuming messages)
	messages, err := ch.Consume(
		q.Name, // 队列名称 (Queue name)
		"",     // 消费者标签 (Consumer tag)
		true,   // 自动确认消息 (Auto-acknowledge?)
		false,  // 是否为独占 (Exclusive?)
		false,  // 是否在本地消费 (No-local?)
		false,  // 是否等待服务器响应 (No-wait?)
		nil,    // 额外参数 (Arguments)
	)
	if err != nil {
		return err
	}

	// 使用 Goroutine 异步处理消息 (Use goroutine to process messages asynchronously)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			// 异步处理每条消息 (Asynchronously handle each message)
			go handlePayload(payload)
		}
	}()

	// 输出等待消息的提示信息 (Print waiting message information)
	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever // 阻塞程序，以保持监听状态 (Block the program to keep listening)

	return nil
}

/*
handlePayload 函数用于处理接收到的消息载荷。

### 函数参数 (Function Parameters)
- `payload Payload`：接收到的消息载荷对象。
  - `payload Payload`: The received message payload object.

### 函数描述 (Function Description)
根据 `payload.Name` 的不同值执行相应的逻辑。当前实现支持以下操作：
- `log`, `event`：将消息记录到日志中。
- `auth`：进行认证操作（示例中未实现）。
- 其他：将消息记录到日志中。

Executes corresponding logic based on the value of `payload.Name`. The current implementation supports the following operations:
- `log`, `event`: Logs the message.
- `auth`: Performs authentication operation (not implemented in the example).
- Others: Logs the message.
*/
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// 记录消息 (Log the message)
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// 处理认证 (Handle authentication)

	// 可以根据需求添加更多处理逻辑 (You can add more cases as needed)

	default:
		// 默认处理逻辑 (Default handling logic)
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

/*
logEvent 函数用于将消息记录到日志服务中。

### 函数参数 (Function Parameters)
- `entry Payload`：要记录的消息载荷对象。
  - `entry Payload`: The message payload object to be logged.

### 返回值 (Return Value)
- `error`：记录消息时发生的任何错误。
  - `error`: Any error that occurs during logging the message.

### 函数描述 (Function Description)
该函数将消息载荷对象序列化为 JSON 格式，并将其发送到日志服务 `http://logger-service/log`。如果 HTTP 响应状态码不是 `202 Accepted`，则返回错误。
This function serializes the message payload object into JSON format and sends it to the log service `http://logger-service/log`. If the HTTP response status code is not `202 Accepted`, it returns an error.
*/
func logEvent(entry Payload) error {
	// 将消息载荷对象序列化为 JSON 格式 (Serialize the message payload object into JSON format)
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log" // 日志服务的 URL (Log service URL)

	// 创建 HTTP POST 请求 (Create HTTP POST request)
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头为 JSON 格式 (Set the request header as JSON format)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{} // 创建 HTTP 客户端 (Create HTTP client)

	// 发送 HTTP 请求并获取响应 (Send the HTTP request and get response)
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close() // 函数结束时关闭响应 (Close the response when function ends)

	// 检查响应状态码是否为 202 Accepted (Check if response status code is 202 Accepted)
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
