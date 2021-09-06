package authorization

import "github.com/casbin/casbin"

type Rbac struct {
	*casbin.Enforcer
}

func (r *Rbac) AddPermissionForUser(user string, permission ...string) bool {
	return r.Enforcer.AddPermissionForUser(user, permission...)
}

func (r *Rbac) AddPermissionForRole(role string, permission ...string) bool {
	return r.Enforcer.AddPermissionForUser("role::"+role, permission...)
}

func (r *Rbac) AddRoleForUser(user string, role string) bool {
	return r.Enforcer.AddRoleForUser(user, "role::"+role)
}

// =========================================================================== //

type Action string

var (
	List           Action = "list"
	Get            Action = "get"
	Create         Action = "create"
	Update         Action = "update"
	Delete         Action = "delete"
	Upload         Action = "upload"
	UploadMultiple Action = "upload_multiple"
	Cancel         Action = "cancel"    // 取消一个未完成的操作
	BatchGet       Action = "batch_get" // 批量获取多个资源
	Move           Action = "move"      // 将资源从一个父级移动到另一个父级
	Search         Action = "search"    // List 的替代方法，用于获取不符合 List 语义的数据
	UnDelete       Action = "undelete"  // 恢复之前删除的资源
)

type User struct {
	Name        string
	Roles       []Role
	Permissions []Permission
}

type Role struct {
	Name        string
	Permissions []Permission
}

type Permission struct {
	Select bool `json:"select"`

	Name    string   `json:"name"`
	Source  string   `json:"source"`
	Actions []Action `json:"actions"`

	MultipleChoice bool `json:"multiple_choice"` //是否是单选，因为有的功能是互斥的
	Priority       int  `json:"priority"`        // 数字越小的，优先级越高，编号不能相同

	Children []Permission `json:"children"`
}

func LoadPermissions() {

}
