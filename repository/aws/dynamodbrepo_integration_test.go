package aws

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jo-hoe/gocommon/repository"
	"github.com/jo-hoe/gocommon/serialization"
)

const testTableName = "testTable"

var defaultConfig = aws.NewConfig().WithRegion("us-west-2").WithEndpoint("http://localhost:8000")

var mockedItem = serialization.MockItem{
	MockString: "mock",
}

var nestedMockedItem = serialization.NestedMockItem{
	NestedItem: mockedItem,
}

func TestMain(m *testing.M) {
	m.Run()
	cleanup()
	os.Exit(0)
}

func cleanup() {
	// clean up created tables if they were created#
	repo, success := connect()
	if !success {
		return
	}
	exists, _ := doesTableExist(repo, testTableName)
	if exists {
		err := dropTable(repo, testTableName)
		if err != nil {
			fmt.Printf("Could not delete table. Error: %v", err)
		}
	}
}

func TestNewDynamoDBRepo(t *testing.T) {
	defer cleanup()
	skipTestIfNoConnectionAvaiable(t)
	mockItem := serialization.MockItem{}
	repo := NewDynamoDBRepo(defaultConfig, testTableName, mockItem.ToStruct)

	if repo == nil {
		t.Errorf("Repo was not initialzed")
	}
}

func TestNewStoreItemDynamoDBRepo(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	if repo == nil {
		t.Errorf("Repo was not initialzed")
	}
}

func TestDynamoDBSave(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)
	testKey := getRandomKey()

	_, err := repo.Save(testKey, mockedItem)
	checkError(err, t)
	actual, err := repo.Find(testKey)

	expected := repository.KeyValuePair{
		Key:   testKey,
		Value: mockedItem,
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v but found %+v. Error: %v", expected, actual, err)
	}
}

func TestDynamoDBFindAllCount(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)
	beforeItems, _ := repo.FindAll()
	beforeCount := len(beforeItems)

	_, err := repo.Save(getRandomKey(), mockedItem)
	checkError(err, t)

	allItems, err := repo.FindAll()
	count := len(allItems)
	if count != beforeCount+1 {
		t.Errorf("Expected %d items but found %d items. Error: %v", beforeCount+1, count, err)
	}
}

func TestDynamoDBFindAll(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	storedItem, _ := repo.Save(getRandomKey(), mockedItem)

	allItems, err := repo.FindAll()
	for _, item := range allItems {
		if item.Key == storedItem.Key {
			if !reflect.DeepEqual(item.Value, mockedItem) {
				t.Errorf("Expected %v but retrieved %v. Error: %v", mockedItem, item.Value, err)
			}
			return
		}
	}
	t.Errorf("Did not find stored item %+v. Error: %v", storedItem, err)
}

func TestDynamoDBSaveValue(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	result, err := repo.Save(getRandomKey(), mockedItem)

	if !reflect.DeepEqual(result.Value, mockedItem) {
		t.Errorf("Expected %+v item but found %+v. Error: %v", mockedItem, result, err)
	}
}

func TestDynamoDBSaveTwiceError(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	_, err := repo.Save(t.Name(), mockedItem)
	checkError(err, t)
	_, err = repo.Save(t.Name(), mockedItem)

	if err == nil {
		t.Error("Expected error was nil although same key was inserted twice")
	}
}

func TestDynamoDBOverwrite(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	_, err := repo.Overwrite(t.Name(), mockedItem)
	checkError(err, t)
	item, err := repo.Overwrite(t.Name(), mockedItem)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if item.Value != mockedItem {
		t.Errorf("Expected %+v but received %+v", mockedItem, item.Value)
	}
}

func TestDynamoDBSaveTwiceLength(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)
	itemsBefore, _ := repo.FindAll()
	countBefore := len(itemsBefore)

	_, err := repo.Save(t.Name(), mockedItem)
	checkError(err, t)
	_, err = repo.Save(t.Name(), mockedItem)
	checkFailure(err, t)

	items, err := repo.FindAll()
	count := len(items)
	if count != countBefore+1 && count > 0 {
		t.Errorf("Expected length to be %d but was %d. Error: %v", countBefore+1, count, err)
	}
}

func TestDynamoDBSaveMultiple(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)
	itemsBefore, _ := repo.FindAll()
	countBefore := len(itemsBefore)

	_, err := repo.Save(getRandomKey(), mockedItem)
	checkError(err, t)
	_, err = repo.Save(getRandomKey(), mockedItem)
	checkError(err, t)

	items, _ := repo.FindAll()
	count := len(items)
	if count != countBefore+2 && count > 0 {
		t.Errorf("Expected length to be %d but was %d. Error: %v", countBefore+2, count, err)
	}
}

func TestDynamoDBDelete(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)
	itemsBefore, err := repo.FindAll()
	checkError(err, t)
	countBefore := len(itemsBefore)

	_, err = repo.Save(getRandomKey(), mockedItem)
	checkError(err, t)
	result, err := repo.Save(getRandomKey(), mockedItem)
	checkError(err, t)
	err = repo.Delete(result.Key)
	checkError(err, t)

	allItems, _ := repo.FindAll()
	count := len(allItems)
	if count != countBefore+1 {
		t.Errorf("Expected %d item but found %d items", countBefore+1, count)
	}
}

func TestDynamoDBDeleteInvalid(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	err := repo.Delete("invalid")

	if err == nil {
		t.Errorf("Error is nil although Delete was called with an invalid value")
	}
}

func TestDynamoDBFind(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	storedItem, _ := repo.Save(getRandomKey(), mockedItem)
	result, err := repo.Find(storedItem.Key)

	actual := result.Value
	if !reflect.DeepEqual(actual, mockedItem) {
		t.Errorf("Expected %+v but retrieved %+v. Error: %v", mockedItem, actual, err)
	}
}

func TestDynamoDBFindNested(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionNestedMockItems(t)

	storedItem, _ := repo.Save(getRandomKey(), nestedMockedItem)
	result, err := repo.Find(storedItem.Key)

	if !reflect.DeepEqual(result.Value, nestedMockedItem) {
		t.Errorf("Expected %v but retrieved %v. Error: %v", mockedItem, result.Value, err)
	}
}

func TestDynamoDBFindInvalid(t *testing.T) {
	defer cleanup()
	repo := createLocalConnectionMockItems(t)

	_, err := repo.Find(getRandomKey())

	if err == nil {
		t.Errorf("Error is nil although Find was called with an invalid value")
	}
}

func createLocalConnectionMockItems(t *testing.T) *DynamoDBRepo {
	skipTestIfNoConnectionAvaiable(t)
	return NewStoreItemDynamoDBRepo(defaultConfig, testTableName, serialization.MockItem{})
}

func createLocalConnectionNestedMockItems(t *testing.T) *DynamoDBRepo {
	skipTestIfNoConnectionAvaiable(t)
	return NewStoreItemDynamoDBRepo(defaultConfig, testTableName, serialization.NestedMockItem{})
}

func skipTestIfNoConnectionAvaiable(t *testing.T) {
	_, success := connect()
	if !success {
		t.Skip("Skipping test because no local db was found.")
	}
}

func connect() (connection *dynamodb.DynamoDB, success bool) {
	connection = GetConnection(defaultConfig)
	if connectionSuccessCache {
		success = isConnected(connection)
		connectionSuccessCache = success
	}
	return connection, connectionSuccessCache
}

var connectionSuccessCache = true

func getRandomKey() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdef" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func checkFailure(err error, t *testing.T) {
	if err == nil {
		t.Error(err)
	}
}
