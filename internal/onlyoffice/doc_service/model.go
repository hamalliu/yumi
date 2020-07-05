package doc_service

type Callback struct {
	Actions        []Action  `json:"actions"`
	Changeshistory []Changes `json:"changeshistory"`
	ChangesUrl     string    `json:"changesurl"`
	Forcesavetype  int       `json:"forcesavetype"`
	History        History   `json:"history"`
	Key            string    `json:"key"`
	Status         int       `json:"status"`
	Url            string    `json:"url"`
	UserData       string    `json:"userdata"`
	Users          []string  `json:"users"`
}

type History struct {
}

type Action struct {
}

type Changes struct {
}
