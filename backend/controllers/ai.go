package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"online-exam-system/config"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AIChatRequest struct {
	Message string `json:"message" binding:"required"`
	Context string `json:"context"` // 可选的上下文信息
}

type AIChatResponse struct {
	Message   string `json:"message"`
	Response  string `json:"response"`
	Timestamp time.Time `json:"timestamp"`
}

type AIAPIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	MaxTokens int        `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIAPIResponse struct {
	Choices []struct {
		Message AIMessage `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// AI问答
func ChatWithAI(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID := middleware.GetCurrentUserID(c)
	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 构建AI请求
	messages := []AIMessage{
		{
			Role:    "system",
			Content: "你是一个专业的在线考试系统AI助手。你可以帮助用户解答关于学习、考试、题目等相关问题。请用中文回答，保持专业和友好的语调。",
		},
	}

	// 如果有上下文，添加到消息中
	if req.Context != "" {
		messages = append(messages, AIMessage{
			Role:    "user",
			Content: "上下文信息：" + req.Context,
		})
	}

	messages = append(messages, AIMessage{
		Role:    "user",
		Content: req.Message,
	})

	aiRequest := AIAPIRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	// 调用AI API
	response, err := callAIAPI(aiRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI服务暂时不可用"})
		return
	}

	// 保存聊天记录
	chatRecord := models.AIChat{
		UserID:    currentUserID,
		Message:   req.Message,
		Response:  response,
		Context:   req.Context,
		CreatedAt: time.Now(),
	}

	// 设置租户ID
	utils.SetTenantID(&chatRecord, tenantID)

	if err := database.DB.Create(&chatRecord).Error; err != nil {
		// 记录保存失败不影响响应
		fmt.Printf("保存AI聊天记录失败: %v\n", err)
	}

	c.JSON(http.StatusOK, AIChatResponse{
		Message:   req.Message,
		Response:  response,
		Timestamp: time.Now(),
	})
}

// 获取AI聊天历史
func GetAIChatHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	offset := (page - 1) * size
	currentUserID := middleware.GetCurrentUserID(c)
	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	var chats []models.AIChat
	var total int64

	// 获取总数
	utils.WithTenant(database.DB, tenantID).Model(&models.AIChat{}).Where("user_id = ?", currentUserID).Count(&total)

	// 获取聊天记录
	if err := utils.WithTenant(database.DB, tenantID).Where("user_id = ?", currentUserID).Offset(offset).Limit(size).Order("created_at DESC").Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取聊天历史失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// 清空AI聊天历史
func ClearAIChatHistory(c *gin.Context) {
	currentUserID := middleware.GetCurrentUserID(c)
	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	if err := utils.WithTenant(database.DB, tenantID).Where("user_id = ?", currentUserID).Delete(&models.AIChat{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空聊天历史失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "聊天历史已清空"})
}

// AI题目解析
func AnalyzeQuestion(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	// 获取租户ID
	tenantID := middleware.GetTenantID(c)

	// 获取题目信息
	var question models.Question
	if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").First(&question, uint(questionID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题目不存在"})
		return
	}

	// 构建分析请求
	var options []string
	if question.Options != "" {
		json.Unmarshal([]byte(question.Options), &options)
	}

	optionsText := ""
	if len(options) > 0 {
		for i, option := range options {
			optionsText += fmt.Sprintf("%c. %s\n", 'A'+i, option)
		}
	}

	prompt := fmt.Sprintf(`请分析以下题目：

科目：%s
题目类型：%s
题目：%s
%s
%s

请提供：
1. 题目知识点分析
2. 解题思路
3. 答案解析
4. 相关知识点扩展

请用中文回答，内容要专业且易于理解。`,
		question.Subject.Name,
		getQuestionTypeText(question.Type),
		question.Title,
		question.Content,
		optionsText)

	messages := []AIMessage{
		{
			Role:    "system",
			Content: "你是一个专业的教育AI助手，擅长分析各种学科的题目，提供详细的解题思路和知识点解析。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	aiRequest := AIAPIRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		MaxTokens:   1500,
		Temperature: 0.3,
	}

	// 调用AI API
	response, err := callAIAPI(aiRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI分析服务暂时不可用"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"question": question,
		"analysis": response,
	})
}

// 调用AI API
func callAIAPI(request AIAPIRequest) (string, error) {
	// 如果没有配置AI API，返回模拟响应
	if config.AppConfig.AIAPIKey == "" || config.AppConfig.AIURL == "" {
		return "抱歉，AI服务暂时不可用。这是一个模拟响应，实际使用时需要配置AI API密钥和URL。", nil
	}

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", config.AppConfig.AIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.AIAPIKey)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 解析响应
	var aiResponse AIAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		return "", err
	}

	// 检查错误
	if aiResponse.Error.Message != "" {
		return "", fmt.Errorf("AI API错误: %s", aiResponse.Error.Message)
	}

	// 返回响应内容
	if len(aiResponse.Choices) > 0 {
		return aiResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("AI API返回空响应")
}

// 获取题目类型文本
func getQuestionTypeText(questionType models.QuestionType) string {
	switch questionType {
	case models.SingleChoice:
		return "单选题"
	case models.MultipleChoice:
		return "多选题"
	case models.TrueFalse:
		return "判断题"
	case models.ShortAnswer:
		return "简答题"
	default:
		return "未知类型"
	}
}