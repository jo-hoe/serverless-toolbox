package aws

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/jo-hoe/gocommon/repository"
	"github.com/jo-hoe/gocommon/serialization"
)

// SSMParameterStoreRepo stores entries in AWS Parameter Store.
// Values will always be stored encrypted
type SSMParameterStoreRepo struct {
	mutex            sync.RWMutex
	path             string
	ssmClient        ssmiface.SSMAPI
	toStructFunction func(jsonString string) (interface{}, error)
}

// NewSSMParameterStoreRepo creates a new instance of the repository
// The repo can take structs and store them in serialized form.
func NewSSMParameterStoreRepo(path string, ssmClient ssmiface.SSMAPI, itemTemplate serialization.Serializable) *SSMParameterStoreRepo {
	return &SSMParameterStoreRepo{
		path:             path,
		ssmClient:        ssmClient,
		toStructFunction: itemTemplate.ToStruct,
	}
}

// NewStringSSMParameterStoreRepo creates a new instance of the repository
// The repo stores the string without conversion.
func NewStringSSMParameterStoreRepo(path string, ssmClient ssmiface.SSMAPI) *SSMParameterStoreRepo {
	return &SSMParameterStoreRepo{
		path:      path,
		ssmClient: ssmClient,
		toStructFunction: func(jsonString string) (interface{}, error) {
			return jsonString, nil
		},
	}

}

func (repo *SSMParameterStoreRepo) FindAll() ([]repository.KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	results := []repository.KeyValuePair{}

	getParametersByPathInput := &ssm.GetParametersByPathInput{
		Path:           aws.String(repo.path),
		WithDecryption: aws.Bool(true),
	}

	err := repo.ssmClient.GetParametersByPathPages(getParametersByPathInput, func(resp *ssm.GetParametersByPathOutput, lastPage bool) bool {
		for _, param := range resp.Parameters {
			fullKey := *param.Name
			item := repository.KeyValuePair{
				Key:   fullKey[:len(repo.path)-1], // remove path from key
				Value: *param.Value,
			}

			results = append(results, item)
		}
		return true
	})

	if err != nil {
		results = nil
	}
	return results, err
}

func (repo *SSMParameterStoreRepo) Save(key string, in interface{}) (repository.KeyValuePair, error) {
	return repo.save(key, in, false)
}

func (repo *SSMParameterStoreRepo) Overwrite(key string, in interface{}) (repository.KeyValuePair, error) {
	return repo.save(key, in, true)
}

func (repo *SSMParameterStoreRepo) save(key string, in interface{}, overwrite bool) (repository.KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result := repository.KeyValuePair{}

	serialized, err := serialization.ToJSON(in)
	if err != nil {
		return result, err
	}

	_, err = repo.ssmClient.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(repo.path + key),
		Value:     aws.String(serialized),
		Type:      aws.String("SecureString"),
		Overwrite: aws.Bool(overwrite),
	})

	if err == nil {
		result.Key = key
		result.Value = in
	}

	return result, err
}

// Only put a variable with the same name >=30 sec after deletion
// Not sure why, but this hint is documented in the AWS docu
// see https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_DeleteParameters.html
func (repo *SSMParameterStoreRepo) Delete(key string) error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	input := &ssm.DeleteParameterInput{
		Name: aws.String(repo.path + key),
	}
	_, err := repo.ssmClient.DeleteParameter(input)
	return err
}

func (repo *SSMParameterStoreRepo) Find(key string) (repository.KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	input := &ssm.GetParameterInput{
		Name:           aws.String(repo.path + key),
		WithDecryption: aws.Bool(true),
	}
	param, err := repo.ssmClient.GetParameter(input)

	if err != nil {
		return repository.KeyValuePair{}, err
	}

	value, err := repo.toStructFunction(*param.Parameter.Value)
	if err != nil {
		return repository.KeyValuePair{}, err
	}

	result := repository.KeyValuePair{
		Key:   key,
		Value: value,
	}

	return result, err
}

func NewSSMSession(region string) ssmiface.SSMAPI {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatalf("Could not initials SSM session %+s", err)
	}

	return ssm.New(sess, aws.NewConfig().WithRegion(region))
}
