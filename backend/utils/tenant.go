package utils

import (
	"reflect"
	"gorm.io/gorm"
)

// WithTenant 为查询添加租户过滤条件
func WithTenant(db *gorm.DB, tenantID uint) *gorm.DB {
	return db.Where("tenant_id = ?", tenantID)
}

// SetTenantID 为模型设置租户ID
func SetTenantID(model interface{}, tenantID uint) {
	switch m := model.(type) {
	case map[string]interface{}:
		m["tenant_id"] = tenantID
	}
}

// TenantScope 租户作用域，用于GORM的Scopes
func TenantScope(tenantID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}

// BeforeCreateTenant GORM钩子，在创建记录前自动设置租户ID
type TenantModel struct {
	TenantID uint `json:"tenant_id" gorm:"not null;index;default:100"`
}

// SetTenantIDForCreate 为创建操作设置租户ID
func SetTenantIDForCreate(db *gorm.DB, tenantID uint) *gorm.DB {
	session := db.Session(&gorm.Session{})
	session.Callback().Create().Before("gorm:create").Register("set_tenant_id", func(db *gorm.DB) {
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("TenantID"); field != nil {
				if db.Statement.ReflectValue.Kind() == reflect.Struct {
					field.Set(db.Statement.Context, db.Statement.ReflectValue, tenantID)
				}
			}
		}
	})
	return session
}