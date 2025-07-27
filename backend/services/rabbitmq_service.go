package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQService RabbitMQ服务
type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// ExamResultMessage 考试结果消息
type ExamResultMessage struct {
	ExamID     uint      `json:"exam_id"`
	UserID     uint      `json:"user_id"`
	Answers    string    `json:"answers"`
	SubmitTime time.Time `json:"submit_time"`
	TenantID   string    `json:"tenant_id"`
}

// ScoreCalculationResult 成绩计算结果
type ScoreCalculationResult struct {
	ExamID       uint    `json:"exam_id"`
	UserID       uint    `json:"user_id"`
	Score        float64 `json:"score"`
	CorrectCount int     `json:"correct_count"`
	TotalCount   int     `json:"total_count"`
	TenantID     string  `json:"tenant_id"`
}

var rabbitMQService *RabbitMQService

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ() error {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}

	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "guest"
	}

	vhost := os.Getenv("RABBITMQ_VHOST")
	if vhost == "" {
		vhost = "/"
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, password, host, port, vhost)

	// 连接到RabbitMQ
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	rabbitMQService = &RabbitMQService{
		conn:    conn,
		channel: ch,
	}

	// 声明队列
	err = rabbitMQService.declareQueues()
	if err != nil {
		return fmt.Errorf("failed to declare queues: %v", err)
	}

	log.Println("RabbitMQ connected successfully")
	return nil
}

// GetRabbitMQService 获取RabbitMQ服务实例
func GetRabbitMQService() *RabbitMQService {
	return rabbitMQService
}

// declareQueues 声明队列
func (r *RabbitMQService) declareQueues() error {
	// 考试结果处理队列
	_, err := r.channel.QueueDeclare(
		"exam_results", // 队列名称
		true,           // 持久化
		false,          // 非独占
		false,          // 非自动删除
		false,          // 非等待
		nil,            // 参数
	)
	if err != nil {
		return err
	}

	// 成绩计算队列
	_, err = r.channel.QueueDeclare(
		"score_calculation",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 统计更新队列
	_, err = r.channel.QueueDeclare(
		"stats_update",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 报告生成队列
	_, err = r.channel.QueueDeclare(
		"report_generation",
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

// PublishExamResult 发布考试结果到队列
func (r *RabbitMQService) PublishExamResult(message ExamResultMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",             // exchange
		"exam_results", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)

	if err != nil {
		log.Printf("Failed to publish exam result: %v", err)
		return err
	}

	log.Printf("Published exam result for user %d, exam %d", message.UserID, message.ExamID)
	return nil
}

// PublishScoreCalculation 发布成绩计算任务
func (r *RabbitMQService) PublishScoreCalculation(message ExamResultMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",
		"score_calculation",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("Failed to publish score calculation: %v", err)
		return err
	}

	log.Printf("Published score calculation for user %d, exam %d", message.UserID, message.ExamID)
	return nil
}

// PublishStatsUpdate 发布统计更新任务
func (r *RabbitMQService) PublishStatsUpdate(result ScoreCalculationResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",
		"stats_update",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("Failed to publish stats update: %v", err)
		return err
	}

	log.Printf("Published stats update for user %d, exam %d", result.UserID, result.ExamID)
	return nil
}

// PublishReportGeneration 发布报告生成任务
func (r *RabbitMQService) PublishReportGeneration(result ScoreCalculationResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",
		"report_generation",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		log.Printf("Failed to publish report generation: %v", err)
		return err
	}

	log.Printf("Published report generation for user %d, exam %d", result.UserID, result.ExamID)
	return nil
}

// Close 关闭RabbitMQ连接
func (r *RabbitMQService) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}

// CloseRabbitMQ 关闭RabbitMQ服务
func CloseRabbitMQ() error {
	if rabbitMQService != nil {
		return rabbitMQService.Close()
	}
	return nil
}