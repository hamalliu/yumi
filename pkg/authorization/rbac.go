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
