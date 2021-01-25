package json

import jsoniter "github.com/json-iterator/go"

// 定义JSON操作
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

func MarshalToJson(v interface{}) (string, error) {
	json, err := jsoniter.MarshalToString(v)
	if err != nil {
		return json, err
	}
	return json, nil
}

func JsonToMarshal(json string, v interface{}) error {
	return Unmarshal([]byte(json), v)
}
