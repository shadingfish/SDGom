/*
Package event 定义了与 RabbitMQ 消息队列相关的事件处理功能。

Package event defines event handling functionalities related to RabbitMQ message queues.
*/

package event

import (
	// 引入 RabbitMQ 的 Go 客户端包，用于与 RabbitMQ 服务器进行通信
	// Importing RabbitMQ's Go client package for communication with the RabbitMQ server.
	amqp "github.com/rabbitmq/amqp091-go"
)

/*
declareExchange 函数用于在 RabbitMQ 中声明一个主题交换机。

### 函数参数 (Function Parameters)
- `ch *amqp.Channel`：RabbitMQ 通道，用于与 RabbitMQ 服务器进行通信。
  - `ch *amqp.Channel`: RabbitMQ channel used for communication with the RabbitMQ server.

### 返回值 (Return Value)
- `error`：返回声明交换机时发生的任何错误。
  - `error`: Returns any error that occurs when declaring the exchange.

### 函数描述 (Function Description)
该函数使用 `ExchangeDeclare` 方法在指定的 RabbitMQ 通道中创建一个名为 `logs_topic` 的主题（topic）类型交换机。
交换机的属性如下：
1. **name**: `logs_topic` 表示交换机的名称。
2. **type**: `topic` 表示交换机的类型是主题交换机，可以根据消息的路由键来将消息发送到不同的队列中。
3. **durable**: `true` 表示交换机在 RabbitMQ 重启时依然会保存。
4. **auto-deleted**: `false` 表示当没有队列绑定到该交换机时不会自动删除。
5. **internal**: `false` 表示该交换机不是内部使用的，允许客户端发送消息到该交换机。
6. **no-wait**: `false` 表示等待服务器确认交换机的声明成功。
7. **arguments**: `nil` 表示没有额外的参数。

This function uses the `ExchangeDeclare` method to create an exchange named `logs_topic` of type `topic` on the specified RabbitMQ channel.
The properties of the exchange are as follows:
1. **name**: `logs_topic` represents the name of the exchange.
2. **type**: `topic` indicates that the exchange is a topic exchange, which can route messages to different queues based on the routing key.
3. **durable**: `true` indicates that the exchange will persist even if RabbitMQ restarts.
4. **auto-deleted**: `false` means that the exchange will not be automatically deleted when no queues are bound to it.
5. **internal**: `false` means that the exchange is not for internal use, allowing clients to publish messages to it.
6. **no-wait**: `false` means that the server will wait for the exchange declaration to be completed before returning.
7. **arguments**: `nil` indicates that no additional arguments are provided.
*/

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguements?
	)
}

/*
declareRandomQueue 函数用于在 RabbitMQ 中声明一个随机生成名称的队列。

### 函数参数 (Function Parameters)
- `ch *amqp.Channel`：RabbitMQ 通道，用于与 RabbitMQ 服务器进行通信。
  - `ch *amqp.Channel`: RabbitMQ channel used for communication with the RabbitMQ server.

### 返回值 (Return Value)
- `amqp.Queue`：声明的队列对象。
  - `amqp.Queue`: The declared queue object.
- `error`：返回声明队列时发生的任何错误。
  - `error`: Returns any error that occurs when declaring the queue.

### 函数描述 (Function Description)
该函数使用 `QueueDeclare` 方法在指定的 RabbitMQ 通道中创建一个随机名称的队列。队列的属性如下：
1. **name**: `""` 表示队列名称为空，RabbitMQ 会自动生成一个唯一的队列名称。
2. **durable**: `false` 表示队列不会在 RabbitMQ 重启时保存（非持久化）。
3. **delete when unused**: `false` 表示队列在不使用时不会被自动删除。
4. **exclusive**: `true` 表示该队列只对当前连接可见，并在连接断开时自动删除。
5. **no-wait**: `false` 表示等待服务器确认队列的声明成功。
6. **arguments**: `nil` 表示没有额外的参数。

This function uses the `QueueDeclare` method to create a queue with a random name on the specified RabbitMQ channel. The properties of the queue are as follows:
1. **name**: `""` indicates that the queue name is empty, and RabbitMQ will automatically generate a unique queue name.
2. **durable**: `false` means that the queue will not be saved if RabbitMQ restarts (non-durable).
3. **delete when unused**: `false` means that the queue will not be automatically deleted when not in use.
4. **exclusive**: `true` means that the queue is only visible to the current connection and will be automatically deleted when the connection is closed.
5. **no-wait**: `false` means that the server will wait for the queue declaration to be completed before returning.
6. **arguments**: `nil` indicates that no additional arguments are provided.
*/

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguments?
	)
}
