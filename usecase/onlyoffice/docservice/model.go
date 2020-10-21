package docservice

//Callback ...
type Callback struct {
	Actions        []Action  `json:"actions"`
	Changeshistory []Changes `json:"changeshistory"`
	ChangesURL     string    `json:"changesurl"`
	Forcesavetype  int       `json:"forcesavetype"`
	History        History   `json:"history"`
	Key            string    `json:"key"`
	Status         int       `json:"status"`
	URL            string    `json:"url"`
	UserData       string    `json:"userdata"`
	Users          []string  `json:"users"`
}

//History ...
type History struct {
}

//Action ...
type Action struct {
}

//Changes ...
type Changes struct {
}
