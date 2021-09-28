package aws

import "encoding/json"

// StoreItem converts a json string into the struct
//
// Example:
// For the struct
// type Person struct {
// 	Name string
// }
//
// The ToStruct function looks like:
// func (person *Person) ToStruct(jsonString string) (interface{}, error) {
//	err := json.Unmarshal([]byte(jsonString), &person)
//	return person, err
// }
type StoreItem interface {
	ToStruct(jsonString string) (interface{}, error)
}

type MockItem struct {
	MockString string
}

// ToStruct converts json string to struct
func (mockItem MockItem) ToStruct(jsonString string) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonString), &mockItem)
	return mockItem, err
}

type NestedMockItem struct {
	NestedItem MockItem
}

// ToStruct converts json string to struct
func (mockItem NestedMockItem) ToStruct(jsonString string) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonString), &mockItem)
	return mockItem, err
}
