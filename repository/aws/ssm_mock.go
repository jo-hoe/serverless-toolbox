package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// for mocking suggestion, refer to https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/ssmiface/
// Define a mock struct to be used in your unit tests.
type mockSSM struct {
	ssmiface.SSMAPI
	mapItem map[string]interface{}
	path    string
}

func NewMockSSM(path string, mapItem map[string]interface{}) *mockSSM {
	return &mockSSM{
		mapItem: mapItem,
		path:    path,
	}
}

func (mock *mockSSM) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	mock.mapItem[*input.Name] = *input.Value
	return nil, nil
}

func (mock *mockSSM) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	result := new(ssm.GetParameterOutput)
	result.Parameter = new(ssm.Parameter)
	err := errors.New("error")

	if val, ok := mock.mapItem[*input.Name]; ok {
		value := fmt.Sprintf("%v", val)
		result.Parameter.Value = &value
		err = nil
	} else {
		result = nil
	}

	return result, err
}

func (mock *mockSSM) DeleteParameter(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	result := new(ssm.DeleteParameterOutput)
	err := errors.New("error")

	if _, ok := mock.mapItem[*input.Name]; ok {
		err = nil
	} else {
		result = nil
	}

	return result, err
}

func (mock *mockSSM) GetParametersByPathPages(input *ssm.GetParametersByPathInput, fn func(*ssm.GetParametersByPathOutput, bool) bool) error {
	err := errors.New("error")

	if *input.Path == mock.path {
		allitems := make([]*ssm.Parameter, 0)
		for key, element := range mock.mapItem {
			param := new(ssm.Parameter)
			param.Name = &key
			value := fmt.Sprintf("%v", element)
			param.Value = &value
			allitems = append(allitems, param)
		}
		fn(&ssm.GetParametersByPathOutput{
			Parameters: allitems,
		}, true)
		err = nil
	}
	return err
}
