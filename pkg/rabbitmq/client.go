package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitmqClient *amqp.Connection
var RabbitmqChannel *amqp.Channel

// NewRabbitmq 初始化 RabbitMQ 连接
func NewRabbitmq(host string, port int) error {
	// 检查 MQ 配置是否为空
	if host == "" {
		fmt.Println("⏭️  RabbitMQ 配置为空，跳过初始化")
		return nil
	}

	fmt.Println("正在初始化 RabbitMQ 连接...")
	amqpHost := fmt.Sprintf("amqp://guest:guest@%s:%d/", host, port)
	conn, err := amqp.Dial(amqpHost)
	if err != nil {
		return fmt.Errorf("RabbitMQ 连接失败: %v", err)
	}
	RabbitmqClient = conn

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("RabbitMQ Channel 创建失败: %v", err)
	}
	RabbitmqChannel = ch

	fmt.Println("✅ RabbitMQ 连接成功")
	return nil
}

// Close 关闭 RabbitMQ 连接
func Close() {
	if RabbitmqChannel != nil {
		RabbitmqChannel.Close()
	}
	if RabbitmqClient != nil {
		RabbitmqClient.Close()
	}
}

// Send 发送消息到指定队列
func Send(queueName string, data string) {
	ch := RabbitmqChannel

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable - 非持久化
		true,      // delete when unused - 自动删除
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", data)
}

// StartQueue 启动队列监听，支持自定义队列名和处理函数
func StartQueue(queueName string, handler func([]byte) error) {
	ch := RabbitmqChannel

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable - 非持久化
		true,      // delete when unused - 自动删除
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message from queue [%s]: %s", queueName, d.Body)

			// 调用业务处理函数
			if err := handler(d.Body); err != nil {
				log.Printf("Error processing message: %v", err)
				d.Nack(false, true) // 消息处理失败，重新入队
			} else {
				d.Ack(false) // 消息处理成功，确认
				log.Printf("Message processed successfully")
			}
		}
	}()

	log.Printf(" [*] Listening on queue: %s", queueName)
}

// ListenQueue 启动队列监听
func ListenQueue() {
	// 检查 RabbitMQ 是否已初始化
	if RabbitmqChannel == nil {
		return
	}

	// 启动示例队列监听
	// 使用示例: StartQueue("your_queue_name", YourHandler)
	// StartQueue("demo_queue", DemoHandler)

	fmt.Println("✅ 队列监听已启动")
}

// DemoHandler 示例消息处理函数
func DemoHandler(body []byte) error {
	// TODO: 在这里添加你的业务逻辑
	log.Printf("Processing message: %s", string(body))
	return nil
}

// failOnError 辅助函数，记录错误日志
func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
