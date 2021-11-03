package authorization

import "github.com/casbin/casbin"

// Rbac ...
type Rbac struct {
	*casbin.Enforcer
}

// AddPermissionForUser ...
func (r *Rbac) AddPermissionForUser(user string, permission ...string) bool {
	return r.Enforcer.AddPermissionForUser(user, permission...)
}

// AddPermissionForRole ...
func (r *Rbac) AddPermissionForRole(role string, permission ...string) bool {
	return r.Enforcer.AddPermissionForUser("role::"+role, permission...)
}

// AddRoleForUser ...
func (r *Rbac) AddRoleForUser(user string, role string) bool {
	return r.Enforcer.AddRoleForUser(user, "role::"+role)
}

// =========================================================================== //

// Action ...
type Action string

var (
	// List ...
	List           Action = "list"
	// Get ...
	Get            Action = "get"
	// Create ...
	Create         Action = "create"
	// Update ...
	Update         Action = "update"
	// Delete ...
	Delete         Action = "delete"
	// Upload ...
	Upload         Action = "upload"
	// UploadMultiple ...
	UploadMultiple Action = "upload_multiple"
	// Cancel ...
	Cancel         Action = "cancel"    // 取消一个未完成的操作
	// BatchGet ...
	BatchGet       Action = "batch_get" // 批量获取多个资源
	// Move ...
	Move           Action = "move"      // 将资源从一个父级移动到另一个父级
	// Search ...
	Search         Action = "search"    // List 的替代方法，用于获取不符合 List 语义的数据
	// UnDelete ...
	UnDelete       Action = "undelete"  // 恢复之前删除的资源
)

// User ...
type User struct {
	Name        string
	Roles       []Role
	Permissions []Permission
}

// Role ...
type Role struct {
	Name        string
	Permissions []Permission
}

// Permission ...
type Permission struct {
	Select bool `json:"select"`

	Name    string   `json:"name"`
	Source  string   `json:"source"`
	Actions []Action `json:"actions"`

	MultipleChoice bool `json:"multiple_choice"` //是否是单选，因为有的功能是互斥的
	Priority       int  `json:"priority"`        // 数字越小的，优先级越高，编号不能相同

	Children []Permission `json:"children"`
}

// LoadPermissions ...
func LoadPermissions() {

}
