package model

import "gorm.io/gorm"

// Role 角色表
type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"` // Added varchar length and unique index
	Description string `gorm:"type:varchar(255)" json:"description"`           // Added varchar length
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"` // For easier loading of permissions
}

// Permission 权限表
type Permission struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"` // Added varchar length and unique index
	Description string `gorm:"type:varchar(255)" json:"description"`           // Added varchar length
}

// RolePermission 角色-权限关联表 (Many2Many join table)
// GORM can often manage this implicitly with `many2many` tag,
// but explicit model can be useful for direct queries or additional fields.
// For this task, we'll rely on GORM's implicit handling via Role.Permissions
// unless direct manipulation of the join table is strictly needed by the middleware logic.
// The issue's example middleware queries RolePermission directly, so we define it.
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"role_id"`
	PermissionID uint `gorm:"primaryKey" json:"permission_id"`
}
// TableName overrides the default table name if needed
func (RolePermission) TableName() string {
	return "role_permissions"
}


// UserRole 用户-角色关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`
	RoleID uint `gorm:"primaryKey" json:"role_id"`
	// User User `gorm:"foreignKey:UserID"` // Optional: define relations
	// Role Role `gorm:"foreignKey:RoleID"` // Optional: define relations
}
// TableName overrides the default table name if needed
func (UserRole) TableName() string {
	return "user_roles"
}

// Helper function to get RBAC models for AutoMigrate
func GetRBACModels() []interface{} {
   return []interface{}{&Role{}, &Permission{}, &RolePermission{}, &UserRole{}}
}
