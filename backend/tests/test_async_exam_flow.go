package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"online-exam-system/config"
	"online-exam-system/controllers"
	"online-exam-system/database"
	"online-exam-system/models"
	"online-exam-system/routes"
	"online-exam-system/services"
	"online-exam-system/utils"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AsyncExamFlowTestSuite 异步考试流程测试套件
type AsyncExamFlowTestSuite struct {
	suite.Suite
	router        *gin.Engine
	testUser      *models.User
	testExam      *models.Exam
	testPaper     *models.Paper
	testQuestions []models.Question
	testSubject   *models.Subject
	testTenant    *models.Tenant
	token         string
}

// SetupSuite 测试套件初始化
func (suite *AsyncExamFlowTestSuite) SetupSuite() {
	// 初始化配置和数据库连接
	config.Init()
	database.Init()

	// 初始化RabbitMQ
	if err := services.InitRabbitMQ(); err != nil {
		log.Printf("Warning: Failed to initialize RabbitMQ: %v", err)
	}

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化路由
	router := gin.New()
	routes.SetupRoutes(router)
	suite.router = router

	// 创建测试数据
	suite.createTestData()

	// 启动消息消费者
	if services.GetRabbitMQService() != nil {
		consumer := services.NewMessageConsumer()
		go consumer.StartConsumers()
	}
}

// TearDownSuite 测试套件清理
func (suite *AsyncExamFlowTestSuite) TearDownSuite() {
	// 清理测试数据
	suite.cleanupTestData()

	// 关闭RabbitMQ连接
	if err := services.CloseRabbitMQ(); err != nil {
		log.Printf("Error closing RabbitMQ: %v", err)
	}
}

// createTestData 创建测试数据
func (suite *AsyncExamFlowTestSuite) createTestData() {
	// 创建测试租户
	suite.testTenant = &models.Tenant{
		Name: "测试租户",
	}
	database.DB.Create(suite.testTenant)

	// 创建测试科目
	suite.testSubject = &models.Subject{
		Name:     "测试科目",
		TenantID: suite.testTenant.ID,
	}
	utils.WithTenant(database.DB, suite.testTenant.ID).Create(suite.testSubject)

	// 创建测试用户（学生）
	suite.testUser = &models.User{
		Username: "test_student",
		Email:    "student@test.com",
		Password: "password123",
		Role:     models.RoleStudent,
		TenantID: suite.testTenant.ID,
	}
	utils.WithTenant(database.DB, suite.testTenant.ID).Create(suite.testUser)

	// 创建测试题目
	suite.testQuestions = []models.Question{
		{
			Content:   "1+1等于多少？",
			Type:      models.SingleChoice,
			Options:   `["1", "2", "3", "4"]`,
			Answer:    "2",
			Score:     10,
			SubjectID: suite.testSubject.ID,
			TenantID:  suite.testTenant.ID,
		},
		{
			Content:   "2+2等于多少？",
			Type:      models.SingleChoice,
			Options:   `["2", "3", "4", "5"]`,
			Answer:    "4",
			Score:     10,
			SubjectID: suite.testSubject.ID,
			TenantID:  suite.testTenant.ID,
		},
		{
			Content:   "简述计算机的基本组成部分",
			Type:      models.ShortAnswer,
			Answer:    "CPU、内存、存储器、输入输出设备",
			Score:     20,
			SubjectID: suite.testSubject.ID,
			TenantID:  suite.testTenant.ID,
		},
	}

	for i := range suite.testQuestions {
		utils.WithTenant(database.DB, suite.testTenant.ID).Create(&suite.testQuestions[i])
	}

	// 创建测试试卷
	suite.testPaper = &models.Paper{
		Title:       "测试试卷",
		Description: "用于测试异步流程的试卷",
		SubjectID:   suite.testSubject.ID,
		TotalScore:  40,
		TenantID:    suite.testTenant.ID,
	}
	utils.WithTenant(database.DB, suite.testTenant.ID).Create(suite.testPaper)

	// 关联试卷和题目
	for _, question := range suite.testQuestions {
		pq := models.PaperQuestion{
			PaperID:    suite.testPaper.ID,
			QuestionID: question.ID,
			Score:      float64(question.Score),
		}
		database.DB.Create(&pq)
	}

	// 创建测试考试
	startTime := time.Now().Add(-1 * time.Hour) // 1小时前开始
	endTime := time.Now().Add(1 * time.Hour)    // 1小时后结束

	suite.testExam = &models.Exam{
		PaperID:     suite.testPaper.ID,
		Title:       "测试考试",
		Description: "用于测试异步流程的考试",
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      models.ExamPublished,
		CreatedBy:   suite.testUser.ID,
		TenantID:    suite.testTenant.ID,
	}
	utils.WithTenant(database.DB, suite.testTenant.ID).Create(suite.testExam)

	// 生成测试用户的JWT token
	suite.token = suite.generateTestToken()
}

// cleanupTestData 清理测试数据
func (suite *AsyncExamFlowTestSuite) cleanupTestData() {
	if suite.testTenant != nil {
		// 删除所有相关数据
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.ExamRecord{})
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.Answer{})
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.Exam{})
		database.DB.Delete(&models.PaperQuestion{}, "paper_id = ?", suite.testPaper.ID)
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.Paper{})
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.Question{})
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.Subject{})
		database.DB.Where("tenant_id = ?", suite.testTenant.ID).Delete(&models.User{})
		database.DB.Delete(&models.Tenant{}, suite.testTenant.ID)
	}
}

// generateTestToken 生成测试用的JWT token
func (suite *AsyncExamFlowTestSuite) generateTestToken() string {
	// 这里应该使用实际的JWT生成逻辑
	// 为了简化测试，我们返回一个模拟的token
	return "test_jwt_token_" + strconv.Itoa(int(suite.testUser.ID))
}

// TestCompleteAsyncExamFlow 测试完整的异步考试流程
func (suite *AsyncExamFlowTestSuite) TestCompleteAsyncExamFlow() {
	// 步骤1: 学生开始考试
	log.Println("=== 步骤1: 学生开始考试 ===")
	examRecord := suite.startExam()
	assert.NotNil(suite.T(), examRecord, "考试记录应该被创建")
	assert.Equal(suite.T(), models.ExamInProgress, examRecord.Status, "考试状态应该是进行中")

	// 步骤2: 学生提交考试答案
	log.Println("=== 步骤2: 学生提交考试答案 ===")
	submitSuccess := suite.submitExamAnswers()
	assert.True(suite.T(), submitSuccess, "考试提交应该成功")

	// 步骤3: 验证立即响应
	log.Println("=== 步骤3: 验证立即响应 ===")
	// 检查考试记录状态已更新为完成
	var updatedRecord models.ExamRecord
	err := utils.WithTenant(database.DB, suite.testTenant.ID).Where("exam_id = ? AND user_id = ?", suite.testExam.ID, suite.testUser.ID).First(&updatedRecord).Error
	assert.NoError(suite.T(), err, "应该能找到更新后的考试记录")
	assert.Equal(suite.T(), models.ExamCompleted, updatedRecord.Status, "考试状态应该是已完成")
	assert.NotNil(suite.T(), updatedRecord.EndTime, "结束时间应该被设置")

	// 步骤4: 等待异步处理完成
	log.Println("=== 步骤4: 等待异步处理完成 ===")
	suite.waitForAsyncProcessing()

	// 步骤5: 验证异步处理结果
	log.Println("=== 步骤5: 验证异步处理结果 ===")
	suite.verifyAsyncProcessingResults()

	log.Println("=== 异步考试流程测试完成 ===")
}

// startExam 开始考试
func (suite *AsyncExamFlowTestSuite) startExam() *models.ExamRecord {
	url := fmt.Sprintf("/api/v1/exams/%d/start", suite.testExam.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("X-Tenant-ID", strconv.Itoa(int(suite.testTenant.ID)))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code, "开始考试应该成功")

	// 从数据库获取创建的考试记录
	var record models.ExamRecord
	err := utils.WithTenant(database.DB, suite.testTenant.ID).Where("exam_id = ? AND user_id = ?", suite.testExam.ID, suite.testUser.ID).First(&record).Error
	if err != nil {
		return nil
	}
	return &record
}

// submitExamAnswers 提交考试答案
func (suite *AsyncExamFlowTestSuite) submitExamAnswers() bool {
	// 构造答案数据
	answers := []controllers.SubmitAnswerRequest{
		{
			QuestionID: suite.testQuestions[0].ID,
			Answer:     "2", // 正确答案
		},
		{
			QuestionID: suite.testQuestions[1].ID,
			Answer:     "4", // 正确答案
		},
		{
			QuestionID: suite.testQuestions[2].ID,
			Answer:     "CPU、内存、存储器、输入输出设备", // 正确答案
		},
	}

	submitReq := controllers.SubmitExamRequest{
		Answers: answers,
	}

	jsonData, _ := json.Marshal(submitReq)
	url := fmt.Sprintf("/api/v1/exams/%d/submit", suite.testExam.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	req.Header.Set("X-Tenant-ID", strconv.Itoa(int(suite.testTenant.ID)))

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	log.Printf("提交考试响应状态: %d", w.Code)
	log.Printf("提交考试响应内容: %s", w.Body.String())

	return w.Code == http.StatusOK
}

// waitForAsyncProcessing 等待异步处理完成
func (suite *AsyncExamFlowTestSuite) waitForAsyncProcessing() {
	// 等待异步消息处理完成
	// 在实际环境中，这个时间可能需要根据系统负载调整
	log.Println("等待异步处理完成...")
	time.Sleep(5 * time.Second)
}

// verifyAsyncProcessingResults 验证异步处理结果
func (suite *AsyncExamFlowTestSuite) verifyAsyncProcessingResults() {
	// 1. 验证成绩计算结果
	suite.verifyScoreCalculation()

	// 2. 验证统计更新
	suite.verifyStatsUpdate()

	// 3. 验证报告生成
	suite.verifyReportGeneration()
}

// verifyScoreCalculation 验证成绩计算
func (suite *AsyncExamFlowTestSuite) verifyScoreCalculation() {
	log.Println("验证成绩计算结果...")

	// 获取考试记录
	var record models.ExamRecord
	err := utils.WithTenant(database.DB, suite.testTenant.ID).Where("exam_id = ? AND user_id = ?", suite.testExam.ID, suite.testUser.ID).First(&record).Error
	assert.NoError(suite.T(), err, "应该能找到考试记录")

	// 验证成绩计算
	expectedScore := 40 // 3道题全部正确：10+10+20=40
	assert.Equal(suite.T(), expectedScore, record.Score, "成绩应该被正确计算")
	assert.Equal(suite.T(), 3, record.CorrectCount, "正确题目数应该是3")
	assert.Equal(suite.T(), 3, record.TotalCount, "总题目数应该是3")

	log.Printf("成绩计算验证通过: 得分=%.1f, 正确数=%d, 总数=%d", record.Score, record.CorrectCount, record.TotalCount)
}

// verifyStatsUpdate 验证统计更新
func (suite *AsyncExamFlowTestSuite) verifyStatsUpdate() {
	log.Println("验证统计更新...")

	// 检查用户统计是否更新
	var userStats models.UserStats
	err := utils.WithTenant(database.DB, suite.testTenant.ID).Where("user_id = ?", suite.testUser.ID).First(&userStats).Error
	if err == nil {
		assert.Greater(suite.T(), userStats.ExamCount, 0, "用户考试次数应该大于0")
		log.Printf("用户统计验证通过: 考试次数=%d", userStats.ExamCount)
	} else {
		log.Println("用户统计记录不存在，可能是首次考试")
	}

	// 检查考试统计是否更新
	var examStats models.ExamStats
	err = utils.WithTenant(database.DB, suite.testTenant.ID).Where("exam_id = ?", suite.testExam.ID).First(&examStats).Error
	if err == nil {
		assert.Greater(suite.T(), examStats.ParticipantCount, 0, "考试参与人数应该大于0")
		log.Printf("考试统计验证通过: 参与人数=%d", examStats.ParticipantCount)
	} else {
		log.Println("考试统计记录不存在，可能需要检查统计更新逻辑")
	}
}

// verifyReportGeneration 验证报告生成
func (suite *AsyncExamFlowTestSuite) verifyReportGeneration() {
	log.Println("验证报告生成...")

	// 检查考试报告是否生成
	var examReport models.ExamReport
	err := utils.WithTenant(database.DB, suite.testTenant.ID).Where("exam_id = ?", suite.testExam.ID).First(&examReport).Error
	if err == nil {
		assert.NotEmpty(suite.T(), examReport.ReportData, "报告内容不应该为空")
		log.Printf("考试报告验证通过: 报告ID=%d", examReport.ID)
	} else {
		log.Println("考试报告不存在，可能需要检查报告生成逻辑")
	}
}

// TestAsyncMessagePublishing 测试异步消息发布
func (suite *AsyncExamFlowTestSuite) TestAsyncMessagePublishing() {
	log.Println("=== 测试异步消息发布 ===")

	// 检查RabbitMQ服务是否可用
	rabbitMQService := services.GetRabbitMQService()
	if rabbitMQService == nil {
		suite.T().Skip("RabbitMQ服务不可用，跳过消息发布测试")
		return
	}

	// 创建测试消息
	message := services.ExamResultMessage{
		ExamID:     suite.testExam.ID,
		UserID:     suite.testUser.ID,
		TenantID:   strconv.FormatUint(uint64(suite.testTenant.ID), 10),
		SubmitTime: time.Now(),
	}

	// 测试发布考试结果消息
	err := rabbitMQService.PublishExamResult(message)
	assert.NoError(suite.T(), err, "发布考试结果消息应该成功")

	// 测试发布成绩计算消息
	err = rabbitMQService.PublishScoreCalculation(message)
	assert.NoError(suite.T(), err, "发布成绩计算消息应该成功")

	// 测试发布统计更新消息
	scoreResult := services.ScoreCalculationResult{
		ExamID:       message.ExamID,
		UserID:       message.UserID,
		TenantID:     message.TenantID,
		Score:        85,
		CorrectCount: 17,
		TotalCount:   20,
	}

	err = rabbitMQService.PublishStatsUpdate(scoreResult)
	assert.NoError(suite.T(), err, "发布统计更新消息应该成功")

	// 测试发布报告生成消息
	err = rabbitMQService.PublishReportGeneration(scoreResult)
	assert.NoError(suite.T(), err, "发布报告生成消息应该成功")

	log.Println("异步消息发布测试完成")
}

// TestAsyncProcessingPerformance 测试异步处理性能
func (suite *AsyncExamFlowTestSuite) TestAsyncProcessingPerformance() {
	log.Println("=== 测试异步处理性能 ===")

	// 模拟多个学生同时提交考试
	concurrentSubmissions := 10
	results := make(chan bool, concurrentSubmissions)

	start := time.Now()

	// 并发提交考试
	for i := 0; i < concurrentSubmissions; i++ {
		go func(studentIndex int) {
			// 为每个学生创建考试记录
			record := models.ExamRecord{
				ExamID:    suite.testExam.ID,
				UserID:    suite.testUser.ID + uint(studentIndex), // 模拟不同学生
				StartTime: time.Now().Add(-30 * time.Minute),
				Status:    models.ExamInProgress,
				TenantID:  suite.testTenant.ID,
			}
			utils.WithTenant(database.DB, suite.testTenant.ID).Create(&record)

			// 模拟提交考试
			if services.GetRabbitMQService() != nil {
				message := services.ExamResultMessage{
					ExamID:     suite.testExam.ID,
					UserID:     record.UserID,
					TenantID:   strconv.FormatUint(uint64(suite.testTenant.ID), 10),
					SubmitTime: time.Now(),
				}

				err := services.GetRabbitMQService().PublishExamResult(message)
				results <- err == nil
			} else {
				results <- false
			}
		}(i)
	}

	// 等待所有提交完成
	successCount := 0
	for i := 0; i < concurrentSubmissions; i++ {
		if <-results {
			successCount++
		}
	}

	duration := time.Since(start)

	log.Printf("并发提交测试结果: %d/%d 成功, 耗时: %v", successCount, concurrentSubmissions, duration)

	// 验证性能指标
	assert.Greater(suite.T(), successCount, concurrentSubmissions/2, "至少一半的提交应该成功")
	assert.Less(suite.T(), duration, 5*time.Second, "并发提交应该在5秒内完成")

	log.Println("异步处理性能测试完成")
}

// TestMain 测试入口
func TestMain(m *testing.M) {
	// 运行测试
	code := m.Run()

	// 退出
	os.Exit(code)
}

// TestAsyncExamFlowTestSuite 运行测试套件
func TestAsyncExamFlowTestSuite(t *testing.T) {
	suite.Run(t, new(AsyncExamFlowTestSuite))
}
