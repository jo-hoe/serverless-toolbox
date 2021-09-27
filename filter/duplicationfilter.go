package filter

import (
	// alternative is
	"log"

	"github.com/jo-hoe/gocommon/repository"
)

// DuplicationFilter stores items returns all items which
// were passed to the function.
type DuplicationFilter interface {
	Filter(structList []interface{}) []interface{}
}

// PersistantDuplicationFilter uses persistency to check if items have been seen before.
type PersistantDuplicationFilter struct {
	itemLog repository.HashKeyValueRepo
}

// NewPersistantDuplicationFilter creates an instance of the DefaultDuplicationFilter
func NewPersistantDuplicationFilter(repo repository.KeyValueRepo) *PersistantDuplicationFilter {
	hashKeyValueRepo := *repository.NewHashKeyValueRepo(repo)
	return &PersistantDuplicationFilter{
		itemLog: hashKeyValueRepo,
	}
}

// Filter a slice of structs for duplicates
func (duplicationfilter *PersistantDuplicationFilter) Filter(structSlice []interface{}) []interface{} {
	// check if items was sent previously
	distinctItemSlice := make([]interface{}, 0)
	for _, item := range structSlice {
		if !duplicationfilter.itemLog.ContainsValue(item) {
			distinctItemSlice = append(distinctItemSlice, item)
			_, err := duplicationfilter.itemLog.Save(item)
			checkError(err)
		}
	}

	return distinctItemSlice
}

func checkError(err error) {
	if err != nil {
		log.Print(err)
	}
}
