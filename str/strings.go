package str

import "encoding/json"

type Strings []string

func (v Strings) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

func (v *Strings) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, v)
	return err
}
