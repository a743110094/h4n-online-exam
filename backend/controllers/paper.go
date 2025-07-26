package controllers

import (
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaperRequest struct {
	SubjectID   uint   `json:"subject_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Duration    int    `json:"duration" binding:"required"` // 考试时长（分钟）
	TotalScore  int    `json:"total_score"`
	Questions   []uint `json:"questions" binding:"required"` // 题目ID列表
}

type AutoPaperRequest struct {
	SubjectID      uint                            `json:"subject_id" binding:"required"`
	Title          string                          `json:"title" binding:"required"`
	Description    string                          `json:"description"`
	Duration       int                             `json:"duration" binding:"required"`
	QuestionConfig []AutoPaperQuestionConfig      `json:"question_config" binding:"required"`
}

type AutoPaperQuestionConfig struct {
	Type       models.QuestionType `json:"type" binding:"required"`
	Count      int                 `json:"count" binding:"required"`
	Difficulty int                 `json:"difficulty"`
	Score      int                 `json:"score" binding:"required"`
}

type PaperListResponse struct {
	Papers []models.Paper `json:"papers"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
}

type PaperDetailResponse struct {
	Paper     models.Paper      `json:"paper"`
	Questions []models.Question `json:"questions"`
}

// 获取试卷列表
func GetPapers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	subjectID := c.Query("subject_id")
	search := c.Query("search")
	tenantID := middleware.GetTenantID(c)

	offset := (page - 1) * size

	query := utils.WithTenant(database.DB, tenantID).Model(&models.Paper{}).Preload("Subject").Preload("Creator")

	// 科目筛选
	if subjectID != "" {
		query = query.Where("subject_id = ?", subjectID)
	}

	// 搜索筛选
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取试卷列表
	var papers []models.Paper
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&papers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取试卷列表失败"})
		return
	}

	c.JSON(http.StatusOK, PaperListResponse{
		Papers: papers,
		Total:  total,
		Page:   page,
		Size:   size,
	})
}

// 获取单个试卷详情
func GetPaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的试卷ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var paper models.Paper
	if err := utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&paper, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "试卷不存在"})
		return
	}

	// 获取试卷题目
	var questions []models.Question
	if err := database.DB.Preload("Subject").Model(&paper).Association("Questions").Find(&questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取试卷题目失败"})
		return
	}

	c.JSON(http.StatusOK, PaperDetailResponse{
		Paper:     paper,
		Questions: questions,
	})
}

// 创建试卷
func CreatePaper(c *gin.Context) {
	var req PaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	// 验证科目是否存在
	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, req.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目不存在"})
		return
	}

	// 验证题目是否存在且属于该科目
	var questions []models.Question
	if err := utils.WithTenant(database.DB, tenantID).Where("id IN ? AND subject_id = ?", req.Questions, req.SubjectID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证题目失败"})
		return
	}

	if len(questions) != len(req.Questions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "部分题目不存在或不属于该科目"})
		return
	}

	// 计算总分
	totalScore := req.TotalScore
	if totalScore == 0 {
		for _, q := range questions {
			totalScore += q.Score
		}
	}

	// 创建试卷
	paper := models.Paper{
		SubjectID:   req.SubjectID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		TotalScore:  totalScore,
		CreatedBy:   middleware.GetCurrentUserID(c),
	}
	utils.SetTenantID(&paper, tenantID)

	if err := database.DB.Create(&paper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建试卷失败"})
		return
	}

	// 关联题目
	if err := database.DB.Model(&paper).Association("Questions").Append(questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "关联题目失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&paper, paper.ID)

	c.JSON(http.StatusCreated, paper)
}

// 自动组卷
func AutoCreatePaper(c *gin.Context) {
	var req AutoPaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	// 验证科目是否存在
	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, req.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目不存在"})
		return
	}

	var selectedQuestions []uint
	var totalScore int

	// 根据配置选择题目
	for _, config := range req.QuestionConfig {
		query := utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("subject_id = ? AND type = ?", req.SubjectID, config.Type)

		// 如果指定了难度
		if config.Difficulty > 0 {
			query = query.Where("difficulty = ?", config.Difficulty)
		}

		var questions []models.Question
		if err := query.Order("RANDOM()").Limit(config.Count).Find(&questions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "选择题目失败"})
			return
		}

		if len(questions) < config.Count {
			c.JSON(http.StatusBadRequest, gin.H{"error": "题目数量不足，无法完成自动组卷"})
			return
		}

		for _, q := range questions {
			selectedQuestions = append(selectedQuestions, q.ID)
			// 使用配置中的分数或题目原有分数
			if config.Score > 0 {
				totalScore += config.Score
			} else {
				totalScore += q.Score
			}
		}
	}

	if len(selectedQuestions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有选择到合适的题目"})
		return
	}

	// 获取选中的题目对象
	var selectedQuestionObjs []models.Question
	if err := utils.WithTenant(database.DB, tenantID).Where("id IN ?", selectedQuestions).Find(&selectedQuestionObjs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取题目失败"})
		return
	}

	// 创建试卷
	paper := models.Paper{
		SubjectID:   req.SubjectID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		TotalScore:  totalScore,
		CreatedBy:   middleware.GetCurrentUserID(c),
	}
	utils.SetTenantID(&paper, tenantID)

	if err := database.DB.Create(&paper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建试卷失败"})
		return
	}

	// 关联题目
	if err := database.DB.Model(&paper).Association("Questions").Append(selectedQuestionObjs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "关联题目失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&paper, paper.ID)

	c.JSON(http.StatusCreated, gin.H{
		"paper":            paper,
		"selected_count":   len(selectedQuestions),
		"total_score":      totalScore,
	})
}

// 更新试卷
func UpdatePaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的试卷ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var req PaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paper models.Paper
	if err := utils.WithTenant(database.DB, tenantID).First(&paper, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "试卷不存在"})
		return
	}

	// 检查权限（只有创建者或管理员可以修改）
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)
	if paper.CreatedBy != currentUserID && currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此试卷"})
		return
	}

	// 检查试卷是否已被使用（有考试关联）
	var examCount int64
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Where("paper_id = ?", uint(id)).Count(&examCount)
	if examCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "试卷已被使用，无法修改"})
		return
	}

	// 验证科目是否存在
	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, req.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目不存在"})
		return
	}

	// 验证题目是否存在且属于该科目
	var questions []models.Question
	if err := utils.WithTenant(database.DB, tenantID).Where("id IN ? AND subject_id = ?", req.Questions, req.SubjectID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证题目失败"})
		return
	}

	if len(questions) != len(req.Questions) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "部分题目不存在或不属于该科目"})
		return
	}

	// 计算总分
	totalScore := req.TotalScore
	if totalScore == 0 {
		for _, q := range questions {
			totalScore += q.Score
		}
	}

	// 更新试卷
	paper.SubjectID = req.SubjectID
	paper.Title = req.Title
	paper.Description = req.Description
	paper.Duration = req.Duration
	paper.TotalScore = totalScore



	// 更新题目关联
	if err := database.DB.Model(&paper).Association("Questions").Replace(questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新题目关联失败"})
		return
	}

	if err := database.DB.Save(&paper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新试卷失败"})
		return
	}

	// 预加载关联数据
	utils.WithTenant(database.DB, tenantID).Preload("Subject").Preload("Creator").First(&paper, paper.ID)

	c.JSON(http.StatusOK, paper)
}

// 删除试卷
func DeletePaper(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的试卷ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var paper models.Paper
	if err := utils.WithTenant(database.DB, tenantID).First(&paper, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "试卷不存在"})
		return
	}

	// 检查权限（只有创建者或管理员可以删除）
	currentUserID := middleware.GetCurrentUserID(c)
	currentRole := middleware.GetCurrentUserRole(c)
	if paper.CreatedBy != currentUserID && currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此试卷"})
		return
	}

	// 检查试卷是否已被使用（有考试关联）
	var examCount int64
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Where("paper_id = ?", uint(id)).Count(&examCount)
	if examCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "试卷已被使用，无法删除"})
		return
	}

	if err := database.DB.Delete(&paper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除试卷失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "试卷删除成功"})
}