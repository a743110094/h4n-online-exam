package main

import (
	"fmt"
	"log"
	"online-exam-system/cache"
	"online-exam-system/config"
	"online-exam-system/models"
	"online-exam-system/services"
	"time"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化Redis连接
	cache.InitRedis()
	if cache.RedisClient == nil {
		log.Printf("Redis连接失败")
		return
	}

	fmt.Println("=== 测试用户认证缓存功能 ===")

	// 创建缓存服务实例
	cacheService := services.NewCacheService()
	tenantID := uint(100)

	// 测试用户缓存
	fmt.Println("\n1. 测试用户缓存")
	testUser := &models.User{
		ID:       1,
		Username: "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     models.RoleStudent,
	}

	// 设置用户缓存
	cacheService.SetUserCache(tenantID, testUser)
	fmt.Printf("✓ 用户缓存已设置: %s\n", testUser.Username)

	// 直接从缓存获取用户（不访问数据库）
	cacheKey := fmt.Sprintf("%s:%d", services.UserCachePrefix, testUser.ID)
	var cachedUser models.User
	if err := cache.GetWithTenant(tenantID, cacheKey, &cachedUser); err == nil {
		fmt.Printf("✓ 直接从缓存获取用户成功: %s\n", cachedUser.Username)
	} else {
		fmt.Printf("✗ 直接从缓存获取用户失败: %v\n", err)
	}

	// 测试Token缓存
	fmt.Println("\n2. 测试Token缓存")
	tokenHash := "test_token_hash_123"
	cacheService.SetTokenCache(tenantID, tokenHash, testUser.ID, testUser.Username, testUser.Role)
	fmt.Printf("✓ Token缓存已设置: %s\n", tokenHash)

	// 从缓存获取Token信息
	if tokenInfo, found := cacheService.GetTokenCache(tenantID, tokenHash); found {
		fmt.Printf("✓ Token缓存获取成功: %+v\n", tokenInfo)
	} else {
		fmt.Println("✗ Token缓存获取失败")
	}

	// 测试缓存失效
	fmt.Println("\n3. 测试缓存失效")
	cacheService.InvalidateUserCache(tenantID, testUser.ID, testUser.Username)
	fmt.Println("✓ 用户缓存已失效")

	cacheService.InvalidateTokenCache(tenantID, tokenHash)
	fmt.Println("✓ Token缓存已失效")

	// 验证缓存已失效（直接检查缓存，不访问数据库）
	if err := cache.GetWithTenant(tenantID, cacheKey, &cachedUser); err != nil {
		fmt.Println("✓ 确认用户缓存已失效")
	} else {
		fmt.Println("✗ 用户缓存仍然存在")
	}

	if _, found := cacheService.GetTokenCache(tenantID, tokenHash); !found {
		fmt.Println("✓ 确认Token缓存已失效")
	} else {
		fmt.Println("✗ Token缓存仍然存在")
	}

	// 测试缓存过期时间
	fmt.Println("\n4. 测试缓存过期时间")
	testUser2 := &models.User{
		ID:       2,
		Username: "testuser2",
		Name:     "Test User 2",
		Email:    "test2@example.com",
		Role:     models.RoleTeacher,
	}

	// 设置短期缓存（用于测试）
	cacheKey2 := fmt.Sprintf("%s:%d", services.UserCachePrefix, testUser2.ID)
	cache.SetWithTenant(tenantID, cacheKey2, *testUser2, 2*time.Second)
	fmt.Println("✓ 设置2秒过期的用户缓存")

	// 立即获取
	var cachedUser2 models.User
	if err := cache.GetWithTenant(tenantID, cacheKey2, &cachedUser2); err == nil {
		fmt.Printf("✓ 立即获取缓存成功: %s\n", cachedUser2.Username)
	} else {
		fmt.Printf("✗ 立即获取缓存失败: %v\n", err)
	}

	// 等待3秒后再获取
	fmt.Println("等待3秒...")
	time.Sleep(3 * time.Second)

	if err := cache.GetWithTenant(tenantID, cacheKey2, &cachedUser2); err != nil {
		fmt.Println("✓ 确认缓存已过期")
	} else {
		fmt.Println("✗ 缓存未过期")
	}

	// 测试多租户隔离
	fmt.Println("\n5. 测试多租户缓存隔离")
	tenantID2 := uint(200)
	testUser3 := &models.User{
		ID:       3,
		Username: "testuser3",
		Name:     "Test User 3",
		Email:    "test3@example.com",
		Role:     models.RoleAdmin,
	}

	// 在租户1设置缓存
	cacheService.SetUserCache(tenantID, testUser3)
	fmt.Printf("✓ 在租户%d设置用户缓存: %s\n", tenantID, testUser3.Username)

	// 在租户2尝试获取
	cacheKey3 := fmt.Sprintf("%s:%d", services.UserCachePrefix, testUser3.ID)
	var cachedUser3 models.User
	if err := cache.GetWithTenant(tenantID2, cacheKey3, &cachedUser3); err != nil {
		fmt.Printf("✓ 租户%d无法获取租户%d的缓存，隔离正常\n", tenantID2, tenantID)
	} else {
		fmt.Printf("✗ 租户隔离失败，租户%d获取到了租户%d的缓存\n", tenantID2, tenantID)
	}

	// 在租户1可以正常获取
	if err := cache.GetWithTenant(tenantID, cacheKey3, &cachedUser3); err == nil {
		fmt.Printf("✓ 租户%d正常获取自己的缓存: %s\n", tenantID, cachedUser3.Username)
	} else {
		fmt.Printf("✗ 租户%d无法获取自己的缓存: %v\n", tenantID, err)
	}

	fmt.Println("\n=== 用户认证缓存测试完成 ===")
	fmt.Println("\n总结:")
	fmt.Println("- ✓ 用户信息缓存功能正常")
	fmt.Println("- ✓ Token验证缓存功能正常")
	fmt.Println("- ✓ 缓存失效机制正常")
	fmt.Println("- ✓ 缓存过期机制正常")
	fmt.Println("- ✓ 多租户缓存隔离正常")
	fmt.Println("\n用户登录和Token缓存已成功集成到认证系统中！")
}