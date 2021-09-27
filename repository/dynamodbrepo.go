package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const keyName = "key"
const valueName = "value"

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

// DynamoDBRepo stores all entities dynamo db
type DynamoDBRepo struct {
	mutex            sync.RWMutex
	tableName        string
	connection       *dynamodb.DynamoDB
	toStructFunction func(jsonString string) (interface{}, error)
}

// GetConnection takes a configuration, creates a session and returns a connection
// the assoicated dynamodb
func GetConnection(config *aws.Config) *dynamodb.DynamoDB {
	session := session.Must(session.NewSession(config))
	return dynamodb.New(session)
}

// NewStoreItemDynamoDBRepo creates a DynamoDBRepo and checks if the table exists. If not it will be created.
func NewStoreItemDynamoDBRepo(config *aws.Config, tableName string, itemTemplate StoreItem) *DynamoDBRepo {
	connection := GetConnection(config)
	exists, _ := doesTableExist(connection, tableName)
	if !exists {
		err := createTable(connection, tableName)
		if err != nil {
			log.Fatalf("Table %s could not be created.", tableName)
		}
	}
	return &DynamoDBRepo{
		tableName:        tableName,
		connection:       connection,
		toStructFunction: itemTemplate.ToStruct,
	}
}

// NewDynamoDBRepo creates a DynamoDBRepo and checks if the table exists. If not it will be created.
// the toStruct function allowed the internal unmarshalling
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
func NewDynamoDBRepo(config *aws.Config, tableName string, toStruct func(jsonString string) (interface{}, error)) *DynamoDBRepo {
	connection := GetConnection(config)
	exists, _ := doesTableExist(connection, tableName)
	if !exists {
		err := createTable(connection, tableName)
		if err != nil {
			log.Fatalf("Table %s could not be created.", tableName)
		}
	}
	return &DynamoDBRepo{
		tableName:        tableName,
		connection:       connection,
		toStructFunction: toStruct,
	}
}

// Save one item
func (repo *DynamoDBRepo) Save(key string, in interface{}) (KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	// converting item to storeable item
	keyValuePair := KeyValuePair{
		Key:   key,
		Value: repo.toJSON(in),
	}

	av, err := dynamodbattribute.MarshalMap(keyValuePair)
	if err != nil {
		return getEmptyKeyValuePair(), err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(repo.tableName),
		ExpressionAttributeNames: map[string]*string{
			"#" + keyName: aws.String(keyName),
		},
		ConditionExpression: aws.String("attribute_not_exists(#" + keyName + ")"),
	}
	_, err = repo.connection.PutItem(input)

	if err != nil {
		return getEmptyKeyValuePair(), err
	}
	keyValuePair.Value = in
	return keyValuePair, nil
}

// FindAll items
func (repo *DynamoDBRepo) FindAll() ([]KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	params := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	result, err := repo.connection.Scan(params)
	if err != nil {
		return []KeyValuePair{}, err
	}

	items := []KeyValuePair{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return []KeyValuePair{}, err
	}

	// convert string into struct
	for i, item := range items {
		result, _ := repo.toStructFunction(item.Value.(string))
		items[i].Value = result
	}

	return items, nil
}

// Delete an item from the repository
func (repo *DynamoDBRepo) Delete(key string) error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
		TableName:    aws.String(repo.tableName),
		ReturnValues: aws.String("ALL_OLD"),
	}
	item, err := repo.connection.DeleteItem(input)
	if err != nil {
		return err
	}
	if item.Attributes == nil {
		return fmt.Errorf("could not find item with key %s", key)
	}
	return nil
}

// Find retrieves and item from the repository
func (repo *DynamoDBRepo) Find(key string) (KeyValuePair, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	result, err := repo.connection.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
	})

	if result.Item == nil && err == nil {
		err = fmt.Errorf("could not find item with key %s", key)
	}
	if err != nil {
		return getEmptyKeyValuePair(), err
	}

	keyValuePair := KeyValuePair{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &keyValuePair)
	if err != nil {
		return getEmptyKeyValuePair(), err
	}
	stringValue := ""
	err = dynamodbattribute.Unmarshal(result.Item[valueName], &stringValue)
	if err != nil {
		return getEmptyKeyValuePair(), err
	}
	storeItem, _ := repo.toStructFunction(stringValue)
	keyValuePair.Value = storeItem
	return keyValuePair, err
}

func doesTableExist(connection *dynamodb.DynamoDB, tableName string) (bool, error) {
	input := &dynamodb.ListTablesInput{}
	result, err := connection.ListTables(input)
	if err != nil {
		return false, err
	}

	for _, name := range result.TableNames {
		if *name == tableName {
			return true, nil
		}
	}

	return false, nil
}

func createTable(connection *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(keyName),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(keyName),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := connection.CreateTable(input)
	return err
}

func dropTable(connection *dynamodb.DynamoDB, tableName string) error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}
	_, err := connection.DeleteTable(input)
	return err
}

func isConnected(connection *dynamodb.DynamoDB) bool {
	timeoutChannel := make(chan bool, 1)
	go func() {
		_, err := doesTableExist(connection, "")
		timeoutChannel <- err == nil
	}()

	select {
	case result := <-timeoutChannel:
		return result
	case <-time.After(8 * time.Second):
		return false
	}
}

func getEmptyKeyValuePair() KeyValuePair {
	return KeyValuePair{
		Key:   "",
		Value: nil,
	}
}

func (repo *DynamoDBRepo) toJSON(in interface{}) string {
	byteArray, _ := json.Marshal(in)
	return string(byteArray)
}
