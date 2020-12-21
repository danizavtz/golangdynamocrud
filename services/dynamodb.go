package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"golangdynamocrud/model"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"sync"
	"os"
)
var singleton sync.Once
var instance *dynamodb.DynamoDB

func GetDynamoInstance() *dynamodb.DynamoDB {
	singleton.Do(func () {
		sess := GetAwsSession()
		instance = dynamodb.New(sess)
	})
	return instance
}

func MarshalMapForAttributes(item model.Usuario) (map[string]*dynamodb.AttributeValue, error) {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, err
	}
	return av, nil
}

func AssembleDynamoItem(itemMarshalled map[string]*dynamodb.AttributeValue) *dynamodb.PutItemInput {
	input := &dynamodb.PutItemInput{
		Item:      itemMarshalled,
		TableName: aws.String(os.Getenv("TABLE")),
	}
	return input
}