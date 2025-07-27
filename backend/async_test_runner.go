package main

import (
	"fmt"
	"log"
	"online-exam-system/config"
	"online-exam-system/database"
	"online-exam-system/services"
	"os"
	"strconv"
	"time"
)

// SimpleAsyncFlowTest 简化的异步流程测试
type SimpleAsyncFlowTest struct {
	testTenantID uint
	testUserID   uint
	testExamID   uint
}

// RunSimpleAsyncTest 运行简化的异步测试
func RunSimpleAsyncTest() {
	log.Println("=== 开始异步流程测试 ===")

	// 初始化配置
	config.Init()

	// 初始化数据库
	database.Init()

	// 初始化RabbitMQ
	if err := services.InitRabbitMQ(); err != nil {
		log.Printf("RabbitMQ初始化失败: %v", err)
		log.Println("请确保RabbitMQ服务正在运行")
		return
	}

	test := &SimpleAsyncFlowTest{
		testTenantID: 1, // 假设存在租户ID为1
		testUserID:   1, // 假设存在用户ID为1
		testExamID:   1, // 假设存在考试ID为1
	}

	// 启动消息消费者
	consumer := services.NewMessageConsumer()
	go consumer.StartConsumers()
	log.Println("消息消费者已启动")

	// 等待消费者启动
	time.Sleep(2 * time.Second)

	// 测试1: 消息发布测试
	test.testMessagePublishing()

	// 测试2: 模拟考试提交流程
	test.testExamSubmissionFlow()

	// 测试3: 并发提交测试
	test.testConcurrentSubmissions()

	// 等待异步处理完成
	log.Println("等待异步处理完成...")
	time.Sleep(10 * time.Second)

	// 关闭连接
	if err := services.CloseRabbitMQ(); err != nil {
		log.Printf("关闭RabbitMQ连接时出错: %v", err)
	}

	log.Println("=== 异步流程测试完成 ===")
}

// testMessagePublishing 测试消息发布
func (t *SimpleAsyncFlowTest) testMessagePublishing() {
	log.Println("\n--- 测试1: 消息发布 ---")

	rabbitMQService := services.GetRabbitMQService()
	if rabbitMQService == nil {
		log.Println("❌ RabbitMQ服务不可用")
		return
	}

	// 创建测试消息
	message := services.ExamResultMessage{
		ExamID:     t.testExamID,
		UserID:     t.testUserID,
		Answers:    `{"1":"A","2":"B","3":"C"}`,
		TenantID:   strconv.FormatUint(uint64(t.testTenantID), 10),
		SubmitTime: time.Now(),
	}

	// 测试发布考试结果消息
	log.Println("发布考试结果消息...")
	if err := rabbitMQService.PublishExamResult(message); err != nil {
		log.Printf("❌ 发布考试结果消息失败: %v", err)
	} else {
		log.Println("✅ 考试结果消息发布成功")
	}

	// 测试发布成绩计算消息
	log.Println("发布成绩计算消息...")
	if err := rabbitMQService.PublishScoreCalculation(message); err != nil {
		log.Printf("❌ 发布成绩计算消息失败: %v", err)
	} else {
		log.Println("✅ 成绩计算消息发布成功")
	}

	// 测试发布统计更新消息
	scoreResult := services.ScoreCalculationResult{
		ExamID:       message.ExamID,
		UserID:       message.UserID,
		TenantID:     message.TenantID,
		Score:        85,
		CorrectCount: 17,
		TotalCount:   20,
	}

	log.Println("发布统计更新消息...")
	if err := rabbitMQService.PublishStatsUpdate(scoreResult); err != nil {
		log.Printf("❌ 发布统计更新消息失败: %v", err)
	} else {
		log.Println("✅ 统计更新消息发布成功")
	}

	// 测试发布报告生成消息
	log.Println("发布报告生成消息...")
	if err := rabbitMQService.PublishReportGeneration(scoreResult); err != nil {
		log.Printf("❌ 发布报告生成消息失败: %v", err)
	} else {
		log.Println("✅ 报告生成消息发布成功")
	}
}

// testExamSubmissionFlow 测试考试提交流程
func (t *SimpleAsyncFlowTest) testExamSubmissionFlow() {
	log.Println("\n--- 测试2: 考试提交流程 ---")

	rabbitMQService := services.GetRabbitMQService()
	if rabbitMQService == nil {
		log.Println("❌ RabbitMQ服务不可用")
		return
	}

	// 模拟学生提交考试
	log.Println("模拟学生提交考试...")
	start := time.Now()

	// 创建考试结果消息
	message := services.ExamResultMessage{
		ExamID:     t.testExamID,
		UserID:     t.testUserID,
		Answers:    `{"1":"A","2":"B","3":"C","4":"D","5":"A"}`,
		TenantID:   strconv.FormatUint(uint64(t.testTenantID), 10),
		SubmitTime: time.Now(),
	}

	// 发布到考试结果队列（这会触发后续的异步处理）
	if err := rabbitMQService.PublishExamResult(message); err != nil {
		log.Printf("❌ 考试提交失败: %v", err)
		return
	}

	responseTime := time.Since(start)
	log.Printf("✅ 考试提交成功，响应时间: %v", responseTime)
	log.Println("📤 异步消息已发送，后台正在处理...")

	// 验证响应时间
	if responseTime < 100*time.Millisecond {
		log.Println("✅ 响应时间优秀 (< 100ms)")
	} else if responseTime < 500*time.Millisecond {
		log.Println("✅ 响应时间良好 (< 500ms)")
	} else {
		log.Println("⚠️ 响应时间较慢 (> 500ms)")
	}
}

// testConcurrentSubmissions 测试并发提交
func (t *SimpleAsyncFlowTest) testConcurrentSubmissions() {
	log.Println("\n--- 测试3: 并发提交 ---")

	rabbitMQService := services.GetRabbitMQService()
	if rabbitMQService == nil {
		log.Println("❌ RabbitMQ服务不可用")
		return
	}

	concurrentCount := 20
	log.Printf("模拟 %d 个学生同时提交考试...", concurrentCount)

	start := time.Now()
	results := make(chan bool, concurrentCount)

	// 并发提交
	for i := 0; i < concurrentCount; i++ {
		go func(studentIndex int) {
			message := services.ExamResultMessage{
				ExamID:     t.testExamID,
				UserID:     t.testUserID + uint(studentIndex),
				Answers:    fmt.Sprintf(`{"1":"A","2":"B","3":"C","student":"%d"}`, studentIndex),
				TenantID:   strconv.FormatUint(uint64(t.testTenantID), 10),
				SubmitTime: time.Now(),
			}

			err := rabbitMQService.PublishExamResult(message)
			results <- err == nil
		}(i)
	}

	// 收集结果
	successCount := 0
	for i := 0; i < concurrentCount; i++ {
		if <-results {
			successCount++
		}
	}

	totalTime := time.Since(start)
	successRate := float64(successCount) / float64(concurrentCount) * 100

	log.Printf("✅ 并发测试完成:")
	log.Printf("   - 总提交数: %d", concurrentCount)
	log.Printf("   - 成功数: %d", successCount)
	log.Printf("   - 成功率: %.1f%%", successRate)
	log.Printf("   - 总耗时: %v", totalTime)
	log.Printf("   - 平均响应时间: %v", totalTime/time.Duration(concurrentCount))

	// 评估性能
	if successRate >= 95 {
		log.Println("✅ 并发性能优秀")
	} else if successRate >= 80 {
		log.Println("✅ 并发性能良好")
	} else {
		log.Println("⚠️ 并发性能需要优化")
	}
}

// main 函数
func main() {
	// 检查命令行参数
	if len(os.Args) > 1 && os.Args[1] == "test" {
		RunSimpleAsyncTest()
	} else {
		fmt.Println("异步流程测试工具")
		fmt.Println("")
		fmt.Println("使用方法:")
		fmt.Println("  go run async_test_runner.go test")
		fmt.Println("")
		fmt.Println("测试内容:")
		fmt.Println("  1. 消息发布测试")
		fmt.Println("  2. 考试提交流程测试")
		fmt.Println("  3. 并发提交性能测试")
		fmt.Println("")
		fmt.Println("注意: 请确保以下服务正在运行:")
		fmt.Println("  - PostgreSQL 数据库")
		fmt.Println("  - RabbitMQ 服务")
	}
}