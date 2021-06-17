package permission

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
	Name       string       `json:"name"`
	SourceType string       `json:"source_type"`
	Actions    []Action     `json:"actions"`
	Children   []Permission `json:"children"`
}

func LoadPermissions() {

}
