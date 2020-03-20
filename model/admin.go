package model

type Login struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Staff struct {
	Id       int    `db:"id" json:"id"`
	User     string `db:"user" json:"user"`
	RealName string `db:"realname" json:"realname"`
	Phone    string `db:"phone" json:"phone"`
}
