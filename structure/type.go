package structure

import "encoding/json"

type KeyValues struct {
    *KeyInfo	`json:"key_info"`
    Value   interface{} `json:"value"`
}

type ValueOf interface {
    Value() interface{}
    String() string
}

func (values *KeyValues) String() string {
    bs,err := json.Marshal(values)
    if err != nil {
        return ""
    }
    return string(bs[:])
}