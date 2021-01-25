package str

import "encoding/json"

type String string

func (v String) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

func (v *String) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, v)
	return err
}
