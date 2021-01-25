package test_model

type TestRole struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Nickname    string           `json:"nickname"`
	Permissions []TestPermission `json:"permissions"`
}

func (r TestRole) GetId() string {
	return r.Id
}
func (r TestRole) GetName() string {
	return r.Name
}
func (r TestRole) GetNickname() string {
	return r.Nickname
}

func (r TestRole) GetPermissions() []string {
	var permissions []string
	for _, permission := range r.Permissions {
		permissions = append(permissions, permission.Permission)
	}
	return permissions
}
