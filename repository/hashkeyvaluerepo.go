package repository

import (
	"crypto/md5"
	"fmt"

	"github.com/google/uuid"
)

// HashKeyValueRepo extends functionalities of KeyValueStore. Implementations are based on
// methods provided in the inner repository. The performance is not optimal since it is not
// based upon a native persistence layer, but instead on an additional layer of abstraction.
type HashKeyValueRepo struct {
	wrappedRepo KeyValueRepo
}

// NewHashKeyValueRepo creates a new instance and uses an initialized KeyValueRepo
func NewHashKeyValueRepo(repo KeyValueRepo) *HashKeyValueRepo {
	return &HashKeyValueRepo{
		wrappedRepo: repo,
	}
}

// FindAll calls function of wrapped repository
func (repo *HashKeyValueRepo) FindAll() ([]KeyValuePair, error) {
	return repo.wrappedRepo.FindAll()
}

// Save calls function of wrapped repository
func (repo *HashKeyValueRepo) Save(in interface{}) (KeyValuePair, error) {
	return repo.wrappedRepo.Save(ToKey(in), in)
}

// Overwrite calls function of wrapped repository
func (repo *HashKeyValueRepo) Overwrite(in interface{}) (KeyValuePair, error) {
	return repo.wrappedRepo.Overwrite(ToKey(in), in)
}

// Delete calls function of wrapped repository
func (repo *HashKeyValueRepo) Delete(key string) error {
	return repo.wrappedRepo.Delete(key)
}

// Find calls function of wrapped repository
func (repo *HashKeyValueRepo) Find(key string) (KeyValuePair, error) {
	return repo.wrappedRepo.Find(key)
}

// Count calls FindAll() and calculates the length
func (repo *HashKeyValueRepo) Count() (int, error) {
	items, _ := repo.FindAll()
	return len(items), nil
}

// Contains checks if a given key is in the repository
func (repo *HashKeyValueRepo) Contains(key string) bool {
	_, err := repo.Find(key)
	return err == nil
}

// ContainsValue checks if a value is in the repository
func (repo *HashKeyValueRepo) ContainsValue(in interface{}) bool {
	return repo.Contains(ToKey(in))
}

// ToKey transforms a struct into a hash key for this repository
func ToKey(in interface{}) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%v", in)))

	id, _ := uuid.FromBytes(h.Sum(nil))
	return id.String()
}
