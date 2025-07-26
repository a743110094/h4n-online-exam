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

type SubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type SubjectListResponse struct {
	Subjects []models.Subject `json:"subjects"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
}

// 获取科目列表
func GetSubjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	search := c.Query("search")
	tenantID := middleware.GetTenantID(c)

	offset := (page - 1) * size

	query := utils.WithTenant(database.DB, tenantID).Model(&models.Subject{})

	// 搜索筛选
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取科目列表
	var subjects []models.Subject
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取科目列表失败"})
		return
	}

	c.JSON(http.StatusOK, SubjectListResponse{
		Subjects: subjects,
		Total:    total,
		Page:     page,
		Size:     size,
	})
}

// 获取所有科目（不分页）
func GetAllSubjects(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var subjects []models.Subject
	if err := utils.WithTenant(database.DB, tenantID).Order("name ASC").Find(&subjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取科目列表失败"})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

// 获取单个科目
func GetSubject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "科目不存在"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// 创建科目
func CreateSubject(c *gin.Context) {
	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	// 检查科目名称是否已存在
	var existingSubject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).Where("name = ?", req.Name).First(&existingSubject).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目名称已存在"})
		return
	}

	// 创建科目
	subject := models.Subject{
		Name:        req.Name,
		Description: req.Description,
	}
	utils.SetTenantID(&subject, tenantID)

	if err := database.DB.Create(&subject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建科目失败"})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

// 更新科目
func UpdateSubject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var req SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "科目不存在"})
		return
	}

	// 检查科目名称是否已存在（排除当前科目）
	var existingSubject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).Where("name = ? AND id != ?", req.Name, uint(id)).First(&existingSubject).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "科目名称已存在"})
		return
	}

	// 更新科目
	subject.Name = req.Name
	subject.Description = req.Description

	if err := database.DB.Save(&subject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新科目失败"})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// 删除科目
func DeleteSubject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "科目不存在"})
		return
	}

	// 检查是否有题目关联此科目
	var questionCount int64
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("subject_id = ?", uint(id)).Count(&questionCount)
	if questionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该科目下还有题目，无法删除"})
		return
	}

	// 检查是否有试卷关联此科目
	var paperCount int64
	utils.WithTenant(database.DB, tenantID).Model(&models.Paper{}).Where("subject_id = ?", uint(id)).Count(&paperCount)
	if paperCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该科目下还有试卷，无法删除"})
		return
	}

	if err := database.DB.Delete(&subject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除科目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "科目删除成功"})
}

// 获取科目统计信息
func GetSubjectStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的科目ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var subject models.Subject
	if err := utils.WithTenant(database.DB, tenantID).First(&subject, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "科目不存在"})
		return
	}

	var stats struct {
		Subject       models.Subject `json:"subject"`
		QuestionCount int64          `json:"question_count"`
		PaperCount    int64          `json:"paper_count"`
		ExamCount     int64          `json:"exam_count"`
	}

	stats.Subject = subject

	// 题目数量
	utils.WithTenant(database.DB, tenantID).Model(&models.Question{}).Where("subject_id = ?", uint(id)).Count(&stats.QuestionCount)

	// 试卷数量
	utils.WithTenant(database.DB, tenantID).Model(&models.Paper{}).Where("subject_id = ?", uint(id)).Count(&stats.PaperCount)

	// 考试数量
	utils.WithTenant(database.DB, tenantID).Model(&models.Exam{}).Joins("JOIN papers ON exams.paper_id = papers.id").Where("papers.subject_id = ?", uint(id)).Count(&stats.ExamCount)

	c.JSON(http.StatusOK, stats)
}