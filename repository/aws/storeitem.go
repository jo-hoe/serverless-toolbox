package aws

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
