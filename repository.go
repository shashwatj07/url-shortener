package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/shashwatj07/url-shortener/blob/dynamodb/entity.go"
)

type dynamoDBRepo struct {
	tableName string
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

func (repo *dynamoDBRepo) Save(post *entity.Urlpair) (*entity.Urlpair, error) {

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

func (repo *dynamoDBRepo) FindAll() ([]entity.Urlpair, error) {
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
	var posts []entity.Urlpair = []entity.Urlpair{}
	for _, i := range result.Items {
		post := entity.Urlpair{}

		err = dynamodbattribute.UnmarshalMap(i, &post)

		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *dynamoDBRepo) FindByID(id string) (*entity.Urlpair, error) {
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
	post := entity.Urlpair{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &post)
	if err != nil {
		panic(err)
	}
	return &post, nil
}

// Delete: TODO
func (repo *dynamoDBRepo) Delete(post *entity.Urlpair) error {
	return nil
}
