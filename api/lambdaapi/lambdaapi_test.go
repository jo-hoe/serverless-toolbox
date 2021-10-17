package lambdaapi

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jo-hoe/gocommon/repository"
)

var mockedItem = MockItem{
	MockString: "mock",
}

// MockItem for testing
type MockItem struct {
	MockString string
}

// ToStruct converts json string to struct
func (mockItem *MockItem) ToStruct(jsonString string) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonString), &mockItem)
	return mockItem, err
}

func TestNewLambdaCrdAPI(t *testing.T) {
	service := NewLambdaCrdAPI(repository.NewInMemoryRepo(), mockedItem.ToStruct)

	if service == nil {
		t.Errorf("Could not create service")
	}
}

func TestLambdaCrdAPI_HTTPMethodProxy_GET(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	_, err := repo.Save("", mockedItem)
	checkError(err, t)
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	request := generateMockedRequest("GET", "", "")

	response, _ := service.HTTPMethodProxy(request)

	if !strings.Contains(response.Body, mockedItem.MockString) {
		t.Errorf("Expected body to contain %+v. Body was actually %+v", mockedItem, response.Body)
	}
}

func TestLambdaCrdAPI_HTTPMethodProxy_405(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	request := generateMockedRequest("NONSENSE", "", "")

	response, _ := service.HTTPMethodProxy(request)

	if response.StatusCode != 405 {
		t.Errorf("Expected status code to be 405 but was %d", response.StatusCode)
	}
}

func TestLambdaCrdAPI_Get(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	_, err := repo.Save("", mockedItem)
	checkError(err, t)
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	request := generateMockedRequest("GET", "", "")

	response, _ := service.Get(request)

	if !strings.Contains(response.Body, mockedItem.MockString) {
		t.Errorf("Expected body to contain %+v. Body was actually %+v", mockedItem, response.Body)
	}
}

func TestLambdaCrdAPI_PostResponse(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	body := "{\"MockString\":\"mock\"}"
	request := generateMockedRequest("POST", "", body)

	response, _ := service.Post(request)

	if !strings.Contains(response.Body, "mock") {
		t.Errorf("Expected body to contain %+v. Body was actually %+v", body, response.Body)
	}
}

func TestLambdaCrdAPI_Post(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	body := "{\"MockString\":\"mock\"}"
	request := generateMockedRequest("POST", "", body)

	_, err := service.Post(request)
	checkError(err, t)

	items, _ := repo.FindAll()
	actual := items[0]
	if reflect.DeepEqual(actual, mockedItem) {
		t.Errorf("Expected %+v but found %+v", mockedItem, actual)
	}
}

func TestLambdaCrdAPI_Delete(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	keyValuePair, err := repo.Save("myKey", mockedItem)
	checkError(err, t)
	request := generateMockedRequest("DELETE", "/somepath/"+keyValuePair.Key, "")

	_, err = service.Delete(request)
	checkError(err, t)

	items, _ := repo.FindAll()
	if len(items) != 0 {
		t.Errorf("No items were expected to be in the repo. Instead found %v", items)
	}
}

func TestLambdaCrdAPI_DeleteInvalid(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	request := generateMockedRequest("DELETE", "invalid", "")

	response, _ := service.Delete(request)

	if response.StatusCode != 400 {
		t.Errorf("Expected response to deliver 400. But received %v", response.StatusCode)
	}
}

func TestLambdaCrdAPI_DeleteNonExisting(t *testing.T) {
	repo := repository.NewInMemoryRepo()
	service := NewLambdaCrdAPI(repo, mockedItem.ToStruct)
	request := generateMockedRequest("DELETE", "/somepath/nonexisting", "")

	response, _ := service.Delete(request)

	// response shall always contain 204 even if item not avaiable
	if response.StatusCode != 204 {
		t.Errorf("Expected response to deliver 204. But received %v", response.StatusCode)
	}
}

func generateMockedRequest(httpMethod string, path string, body string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		HTTPMethod: httpMethod,
		Body:       body,
		Path:       path,
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
