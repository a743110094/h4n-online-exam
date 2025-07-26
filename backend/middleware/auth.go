package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"online-exam-system/config"
	"online-exam-system/models"
	"online-exam-system/services"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 全局缓存服务实例
var cacheService = services.NewCacheService()

type Claims struct {
	UserID   uint             `json:"user_id"`
	Username string           `json:"username"`
	Role     models.UserRole  `json:"role"`
	jwt.RegisteredClaims
}

// 生成JWT Token
func GenerateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWTSecret))
}

// 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// 获取租户ID
		tenantID := GetTenantID(c)

		// 尝试从缓存获取token信息
		tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(tokenString)))
		if tokenInfo, found := cacheService.GetTokenCache(tenantID, tokenHash); found {
			// 从缓存中获取到token信息，直接使用
			if userID, ok := tokenInfo["user_id"].(float64); ok {
				c.Set("user_id", uint(userID))
			}
			if username, ok := tokenInfo["username"].(string); ok {
				c.Set("username", username)
			}
			if roleStr, ok := tokenInfo["role"].(string); ok {
				c.Set("role", models.UserRole(roleStr))
			}
			c.Next()
			return
		}

		// 缓存未命中，解析JWT token
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 缓存token信息
		cacheService.SetTokenCache(tenantID, tokenHash, claims.UserID, claims.Username, claims.Role)

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// 角色权限中间件
func RoleMiddleware(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role := userRole.(models.UserRole)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// 获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	userID, _ := c.Get("user_id")
	return userID.(uint)
}

// 获取用户ID (别名函数)
func GetUserID(c *gin.Context) uint {
	return GetCurrentUserID(c)
}

// 获取当前用户角色
func GetCurrentUserRole(c *gin.Context) models.UserRole {
	role, _ := c.Get("role")
	return role.(models.UserRole)
}