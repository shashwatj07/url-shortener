package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type dynamoDBRepo struct {
	tableName string
}

type PostRepository interface {
	Save(post *Urlpair) (*Urlpair, error)
	FindAll() ([]Urlpair, error)
	FindByID(id string) (*Urlpair, error)
	Delete(post *Urlpair) error
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

func (repo *dynamoDBRepo) Save(post *Urlpair) (*Urlpair, error) {

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

func (repo *dynamoDBRepo) FindAll() ([]Urlpair, error) {
	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	// Make the DynamoDB Query API call
	result, err := dynamoDBClient.Scan(params)
	if err != nil {
		return nil, err
	}
	var posts []Urlpair = []Urlpair{}
	for _, i := range result.Items {
		post := Urlpair{}

		err = dynamodbattribute.UnmarshalMap(i, &post)

		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *dynamoDBRepo) FindByID(id string) (*Urlpair, error) {
	// Get a new DynamoDB client
	dynamoDBClient := createDynamoDBClient()

	result, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	post := Urlpair{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		panic(err)
	}
	return &post, nil
}

// Delete: TODO
func (repo *dynamoDBRepo) Delete(post *Urlpair) error {
	return nil
}
