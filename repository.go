package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var repo = NewDynamoDBRepository()
var dynamoDBClient = createDynamoDBClient()

// Struct to store table name
type dynamoDBRepo struct {
	tableName string
}

// Interface to specify the operations to be done on the database
type PostRepository interface {
	Save(post *urlStruct) (*urlStruct, error)
	FindByID(id string) (*urlStruct, error)
	Delete(id string) error
}

// Instantiatethe dynamoDBRepo struct with the urlKV table
// return: the struct instance
func NewDynamoDBRepository() PostRepository {
	return &dynamoDBRepo{
		tableName: "urlKV",
	}
}

// Create a DynamoDB client session
// return: The new DynamoDB client session
func createDynamoDBClient() *dynamodb.DynamoDB {
	// Create AWS Session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{Region: aws.String("ap-south-1")},
	}))

	// Return DynamoDB client
	return dynamodb.New(sess)
}

// parameter post: The long-short url pair and the expiry date
// The method is written for dynamoDBRepo struct. Save the new url entry pair and its expiry date in the database
// return: url pair, expiry date and error if any
func (repo *dynamoDBRepo) Save(post *urlStruct) (*urlStruct, error) {

	// Transforms the post to map[string]*dynamodb.AttributeValue
	attributeValue, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		return nil, err
	}

	// Create the Item Input
	item := &dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(repo.tableName),
	}

	// Save the Item into DynamoDB
	_, err = dynamoDBClient.PutItem(item)
	if err != nil {
		return nil, err
	}

	return post, err
}

// parameter short_url: The short url to search for in database
// The method is written for dynamoDBRepo. Find item by id which is the short url
// return: The url pair if found and error if any
func (repo *dynamoDBRepo) FindByID(short_url string) (*urlStruct, error) {

	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(short_url),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	post := urlStruct{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		return nil, err
	}
	println(post.LongURL)
	return &post, nil
}

// parameter short_url: The short url to be deleted
// The method is written for dynamoDBRepo. Delete item by id which is the short url
// return: error if any
func (repo *dynamoDBRepo) Delete(short_url string) error {
	item := &dynamodb.DeleteItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(short_url),
			},
		},
	}

	_, err := dynamoDBClient.DeleteItem(item)
	return err
}
