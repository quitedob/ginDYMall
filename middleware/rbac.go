package middleware

import (
	"strings"
	"net/http" // For http status codes

	"github.com/gin-gonic/gin"
	"douyin/pkg/utils/response" // Adjust path if necessary
	"douyin/repository/db/model" // Adjust path if necessary
	"gorm.io/gorm"
	// "douyin/pkg/utils/jwt" // Assuming GetUserIDFromContext exists here or similar
	// We'll rely on c.Get("userID") as per the task description.
)

// RBAC provides role-based access control.
// It checks if the user associated with the current request has the required permission.
func RBAC(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get userID from context (set by AuthMiddleware)
		userIDAny, exists := c.Get("userID") // As per issue's AuthMiddleware hint
		if !exists {
			response.Fail(c, http.StatusUnauthorized, "用户未登录 (User not logged in)")
			c.Abort()
			return
		}

		userID, ok := userIDAny.(uint) // Assuming userID is uint
		if !ok {
			// Attempt to convert from other numeric types, e.g., float64 if JSON unmarshalling decoded it that way
			if userIDFloat, okFloat := userIDAny.(float64); okFloat {
				userID = uint(userIDFloat)
				ok = true
			} else if userIDInt, okInt := userIDAny.(int); okInt {
				userID = uint(userIDInt)
				ok = true
			}
			if !ok {
				response.Fail(c, http.StatusForbidden, "用户ID格式无效 (Invalid user ID format)")
				c.Abort()
				return
			}
		}


		// 2. Get *gorm.DB from context (set by DBInjectorMiddleware)
		dbAny, exists := c.Get("db")
		if !exists {
			response.Fail(c, http.StatusInternalServerError, "数据库连接未设置 (DB connection not set in context)")
			c.Abort()
			return
		}
		db, ok := dbAny.(*gorm.DB)
		if !ok || db == nil {
			response.Fail(c, http.StatusInternalServerError, "无效的数据库连接 (Invalid DB connection in context)")
			c.Abort()
			return
		}

		// 3. Query user_roles
		var userRoles []model.UserRole
		if err := db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询用户角色失败 (Failed to query user roles): "+err.Error())
			c.Abort()
			return
		}
		if len(userRoles) == 0 {
			response.Fail(c, http.StatusForbidden, "用户未分配任何角色 (User has no assigned roles)")
			c.Abort()
			return
		}

		roleIDs := make([]uint, len(userRoles))
		for i, ur := range userRoles {
			roleIDs[i] = ur.RoleID
		}

		// 4. Query role_permissions
		var rolePermissions []model.RolePermission
		if err := db.Where("role_id IN ?", roleIDs).Find(&rolePermissions).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询角色权限失败 (Failed to query role permissions): "+err.Error())
			c.Abort()
			return
		}
		if len(rolePermissions) == 0 {
			response.Fail(c, http.StatusForbidden, "角色没有任何权限 (Roles have no permissions)")
			c.Abort()
			return
		}

		permIDs := make([]uint, len(rolePermissions))
		for i, rp := range rolePermissions {
			permIDs[i] = rp.PermissionID
		}

		// 5. Query permissions table and check against requiredPerm
		var userPermissions []model.Permission
		if err := db.Where("id IN ?", permIDs).Find(&userPermissions).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询权限详情失败 (Failed to query permission details): "+err.Error())
			c.Abort()
			return
		}

		for _, p := range userPermissions {
			if strings.EqualFold(p.Name, requiredPerm) {
				c.Next() // Permission granted
				return
			}
		}

		// 6. No matching permission found
		response.Fail(c, http.StatusForbidden, "无权限访问此资源 (Forbidden: You don't have the required permission)")
		c.Abort()
	}
}
