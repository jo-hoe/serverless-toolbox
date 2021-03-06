package repository

import "encoding/json"

// KeyValuePair has the stored entity in addition to an autogenerated id
type KeyValuePair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// KeyValueRepo is a generic repository which accepts values of type interface
type KeyValueRepo interface {
	FindAll() ([]KeyValuePair, error)
	// saves an item, if the key of the item already exists an error is returned
	Save(key string, in interface{}) (KeyValuePair, error)
	// saves an item, if the key of the item already exists it is overwritten
	Overwrite(key string, in interface{}) (KeyValuePair, error)
	Delete(key string) error
	Find(key string) (KeyValuePair, error)
}

// ToStruct converts json string to struct
func (item KeyValuePair) ToStruct(jsonString string) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonString), &item)
	return item, err
}
