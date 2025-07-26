package controllers

import (
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserListResponse struct {
	Users []models.User `json:"users"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	role := c.Query("role")
	search := c.Query("search")
	tenantID := middleware.GetTenantID(c)

	offset := (page - 1) * size

	query := utils.WithTenant(database.DB, tenantID).Model(&models.User{})

	// 角色筛选
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// 搜索筛选
	if search != "" {
		query = query.Where("username ILIKE ? OR name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取用户列表
	var users []models.User
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, UserListResponse{
		Users: users,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// 获取单个用户信息
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// 创建用户
func CreateUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	// 检查用户名是否已存在
	var existingUser models.User
	if err := utils.WithTenant(database.DB, tenantID).Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := utils.WithTenant(database.DB, tenantID).Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     req.Role,
		IsActive: true,
	}
	utils.SetTenantID(&user, tenantID)

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

// 更新用户信息
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var req struct {
		Username string          `json:"username"`
		Email    string          `json:"email"`
		Name     string          `json:"name"`
		Role     models.UserRole `json:"role"`
		IsActive bool            `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查用户名是否已被其他用户使用
	if req.Username != user.Username {
		var existingUser models.User
		if err := utils.WithTenant(database.DB, tenantID).Where("username = ? AND id != ?", req.Username, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != user.Email {
		var existingUser models.User
		if err := utils.WithTenant(database.DB, tenantID).Where("email = ? AND id != ?", req.Email, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
			return
		}
	}

	// 更新用户信息
	user.Username = req.Username
	user.Email = req.Email
	user.Name = req.Name
	user.Role = req.Role
	user.IsActive = req.IsActive

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 不能删除管理员账号
	if user.Role == models.RoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除管理员账号"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// 重置用户密码
func ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码重置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// 批量导入用户
func BatchImportUsers(c *gin.Context) {
	var req struct {
		Users []RegisterRequest `json:"users" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantID := middleware.GetTenantID(c)

	var successCount int
	var errors []string

	for _, userReq := range req.Users {
		// 检查用户名是否已存在
		var existingUser models.User
		if err := utils.WithTenant(database.DB, tenantID).Where("username = ?", userReq.Username).First(&existingUser).Error; err == nil {
			errors = append(errors, "用户名 "+userReq.Username+" 已存在")
			continue
		}

		// 检查邮箱是否已存在
		if err := utils.WithTenant(database.DB, tenantID).Where("email = ?", userReq.Email).First(&existingUser).Error; err == nil {
			errors = append(errors, "邮箱 "+userReq.Email+" 已存在")
			continue
		}

		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
		if err != nil {
			errors = append(errors, "用户 "+userReq.Username+" 密码加密失败")
			continue
		}

		// 创建用户
		user := models.User{
			Username: userReq.Username,
			Email:    userReq.Email,
			Password: string(hashedPassword),
			Name:     userReq.Name,
			Role:     userReq.Role,
			IsActive: true,
		}
		utils.SetTenantID(&user, tenantID)

		if err := database.DB.Create(&user).Error; err != nil {
			errors = append(errors, "用户 "+userReq.Username+" 创建失败")
			continue
		}

		successCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "批量导入完成",
		"success_count": successCount,
		"errors":        errors,
	})
}