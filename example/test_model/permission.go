package test_model

type TestPermission struct {
	Id         string `json:"id"`
	Url        string `json:"url"`
	Method     string `json:"method"`
	Remarks    string `json:"remarks"`
	Permission string `json:"permission"`
	Anonymous  bool   `json:"anonymous"`
}
