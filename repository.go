package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var repo = NewDynamoDBRepository()
var dynamoDBClient = createDynamoDBClient()

type dynamoDBRepo struct {
	tableName string
}

type PostRepository interface {
	Save(post *urlStruct) (*urlStruct, error)
	FindByID(id string) (*urlStruct, error)
	Delete(post *urlStruct) error
}

func NewDynamoDBRepository() PostRepository {
	return &dynamoDBRepo{
		tableName: "urlKV",
	}
}

func createDynamoDBClient() *dynamodb.DynamoDB {
	// Create AWS Session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Return DynamoDB client
	return dynamodb.New(sess)
}

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

// Delete: TODO
func (repo *dynamoDBRepo) Delete(post *urlStruct) error {
	return nil
}
