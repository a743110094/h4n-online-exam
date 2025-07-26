package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TenantMiddleware 租户中间件
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取租户ID
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		
		// 如果没有提供租户ID，使用默认值100
		if tenantIDStr == "" {
			tenantIDStr = "100"
		}
		
		// 转换为整数
		tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
		if err != nil {
			// 如果转换失败，使用默认值100
			tenantID = 100
		}
		
		// 将租户ID存储到上下文中
		c.Set("tenant_id", uint(tenantID))
		
		c.Next()
	}
}

// GetTenantID 从上下文中获取租户ID
func GetTenantID(c *gin.Context) uint {
	if tenantID, exists := c.Get("tenant_id"); exists {
		if id, ok := tenantID.(uint); ok {
			return id
		}
	}
	return 100 // 默认租户ID
}

// RequireTenant 要求必须提供有效的租户ID
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "租户ID不能为空",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}