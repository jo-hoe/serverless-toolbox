package lambdaapi

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jo-hoe/serverless-toolbox/repository"
)

// LambdaCrdAPI Lambda implementation of CrdAPI
type LambdaCrdAPI struct {
	repo             repository.HashKeyValueRepo
	toStructFunction func(jsonString string) (interface{}, error)
}

// NewLambdaCrdAPI generating a struct to allow CRD actions
func NewLambdaCrdAPI(repo repository.KeyValueRepo, toStructFunction func(jsonString string) (interface{}, error)) *LambdaCrdAPI {
	hashKeyValueRepo := *repository.NewHashKeyValueRepo(repo)
	return &LambdaCrdAPI{
		repo:             hashKeyValueRepo,
		toStructFunction: toStructFunction,
	}
}

// HTTPMethodProxy proxies requests to CRD method based on HTTP method
func (lambdaCrdAPI *LambdaCrdAPI) HTTPMethodProxy(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		return lambdaCrdAPI.Get(request)
	case "POST":
		return lambdaCrdAPI.Post(request)
	case "DELETE":
		return lambdaCrdAPI.Delete(request)
	// return 405 - Method Not Allowed
	default:
		return &events.APIGatewayProxyResponse{
			StatusCode: 405,
		}, nil
	}
}

// Get method returns all stored entities
func (lambdaCrdAPI *LambdaCrdAPI) Get(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	items, err := lambdaCrdAPI.repo.FindAll()
	jsonString, _ := toJSON(items)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonString,
	}, err
}

// Post method create a new entity
func (lambdaCrdAPI *LambdaCrdAPI) Post(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	requestBodyItem, _ := lambdaCrdAPI.toStructFunction(request.Body)
	item, err := lambdaCrdAPI.repo.Save(requestBodyItem)

	statusCode := 400
	jsonString := ""

	if err == nil {
		jsonString, err = toJSON(item)
	}

	if err == nil {
		statusCode = 200
	} else {
		log.Printf("Error during post request processing %+v", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       jsonString,
	}, err
}

// Delete removes an entity
func (lambdaCrdAPI *LambdaCrdAPI) Delete(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	regex := regexp.MustCompile(`(?:\/)?(?:.*\/)(?P<id>.+)`)
	match := regex.FindStringSubmatch(request.Path)

	statusCode := 400
	var err error
	if len(match) == 2 {
		err = lambdaCrdAPI.repo.Delete(match[1])
		statusCode = 204
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
	}, err
}

func toJSON(item interface{}) (string, error) {
	byteArray, err := json.MarshalIndent(item, "", "    ")
	jsonString := string(byteArray)
	return jsonString, err
}
