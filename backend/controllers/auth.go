package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"online-exam-system/database"
	"online-exam-system/middleware"
	"online-exam-system/models"
	"online-exam-system/services"
	"online-exam-system/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 全局缓存服务实例
var cacheService = services.NewCacheService()

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string          `json:"username" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Name     string          `json:"name" binding:"required"`
	Role     models.UserRole `json:"role"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := middleware.GetTenantID(c)

	// 尝试从缓存获取用户信息（包含密码用于验证）
	user, err := cacheService.GetUserByUsernameWithCache(tenantID, req.Username)
	if err != nil {
		log.Println("用户名或密码错误")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 缓存token信息（用于快速验证）
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	cacheService.SetTokenCache(tenantID, tokenHash, user.ID, user.Username, user.Role)

	// 更新用户缓存（登录成功后刷新缓存）
	cacheService.SetUserCache(tenantID, user)

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}

// 用户注册（仅管理员可创建教师和学生账号）
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查权限
	currentRole := middleware.GetCurrentUserRole(c)
	if currentRole != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以创建账号"})
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

	// 设置租户ID
	utils.SetTenantID(&user, tenantID)

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户创建成功",
		"user":    user,
	})
}

// 获取当前用户信息
func GetProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	tenantID := middleware.GetTenantID(c)

	// 尝试从缓存获取用户信息
	user, err := cacheService.GetUserWithCache(tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, *user)
}

// 更新用户信息
func UpdateProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 保存原用户名用于缓存失效
	oldUsername := user.Username

	// 更新用户信息
	user.Name = req.Name
	user.Email = req.Email
	user.Avatar = req.Avatar

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	// 使旧缓存失效
	cacheService.InvalidateUserCache(tenantID, userID, oldUsername)

	// 设置新缓存
	cacheService.SetUserCache(tenantID, &user)

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// 修改密码
func ChangePassword(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	tenantID := middleware.GetTenantID(c)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := utils.WithTenant(database.DB, tenantID).First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码错误"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码更新失败"})
		return
	}

	// 密码修改后，使用户缓存失效（因为密码已变更）
	cacheService.InvalidateUserCache(tenantID, userID, user.Username)

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
