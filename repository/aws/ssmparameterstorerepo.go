package aws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/jo-hoe/gocommon/repository"
)

// SSMParameterStoreRepo stores entries in AWS Parameter Store.
// Values will always be stored encrypted
type SSMParameterStoreRepo struct {
	path             string
	ssmClient        ssmiface.SSMAPI
	toStructFunction func(jsonString string) (interface{}, error)
}

// NewSSMParameterStoreRepo creates a new instance of the repository
// The repo can take structs and store them in serialized form.
func NewSSMParameterStoreRepo(path string, ssmClient ssmiface.SSMAPI, itemTemplate StoreItem) *SSMParameterStoreRepo {
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
	result := repository.KeyValuePair{}

	_, err := repo.ssmClient.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(repo.path + key),
		Value:     aws.String(repo.toJSON(in)),
		Type:      aws.String("SecureString"),
		Overwrite: aws.Bool(false),
	})

	if err == nil {
		result.Key = key
		result.Value = in
	}

	return result, err
}

func (repo *SSMParameterStoreRepo) Delete(key string) error {
	input := &ssm.DeleteParameterInput{
		Name: aws.String(repo.path + key),
	}
	_, err := repo.ssmClient.DeleteParameter(input)
	return err
}

func (repo *SSMParameterStoreRepo) Find(key string) (repository.KeyValuePair, error) {
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

// Converts a struct into a json string. If input is a string and not a struct
// no conversion will take place and string is returned instead
func (repo *SSMParameterStoreRepo) toJSON(in interface{}) string {
	// check if input is pure string
	if _, ok := in.(string); ok {
		return fmt.Sprintf("%v", in) // do no conversion if input is pure string
	} else {
		byteArray, _ := json.Marshal(in)
		return string(byteArray)
	}
}
