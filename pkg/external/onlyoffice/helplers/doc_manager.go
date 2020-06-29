package helplers

type History struct {
	ServerVersion string   `json:"serverVersion"`
	Changes       []Change `json:"changes"`

	User    User   `json:"user"`
	Created string `json:"created"`

	Key     string `json:"key"`
	Version int
}

type Change struct {
	Created string `json:"created"`
	User    User   `json:"user"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type HistoryData struct {
	ChangesUrl string      `json:"changesUrl"`
	Key        string      `json:"key"`
	Pervious   Pervious    `json:"pervious"`
	Url        string      `json:"url"`
	Version    interface{} `json:"version"`
}

type Pervious struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

//======================================================================================================================
type DocManager struct {
}

func (dm DocManager) CreateDemo() string {
	return ""
}

func (dm DocManager) SaveFileData() string {
	return ""
}

func (dm DocManager) GetFileData() string {
	return ""
}

func (dm DocManager) GetFileUri() string {
	return ""
}

func (dm DocManager) GetLocalFileUri() string {
	return ""
}

func (dm DocManager) GetServerUrl() string {
	return ""
}

func (dm DocManager) GetCallback() string {
	return ""
}

func (dm DocManager) GetKey() string {
	return ""
}

func (dm DocManager) GetDate() string {
	return ""
}

func (dm DocManager) GetChanges() string {
	return ""
}

func (dm DocManager) CountVersion() string {
	return ""
}

func (dm DocManager) GetHistory() string {
	return ""
}
