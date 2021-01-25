package model

type IRole interface {
	GetId() string

	GetName() string

	GetNickname() string

	GetPermissions() []string
}

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Nickname    string   `json:"nickname"`
	Permissions []string `json:"permissions"`
}
