package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"golangdynamocrud/model"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"sync"
	"os"
	// "strings"
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

func assembleItemForGetById(inputIdentifier string) (*dynamodb.GetItemInput){
	item := dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String("users:"+inputIdentifier),
			},
		},
	}
	return &item
}

func assembleItemForDeleteById(inputIdentifier string) (*dynamodb.DeleteItemInput) {
	item := dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String("users:"+inputIdentifier),
			},
		},
		TableName: aws.String(os.Getenv("TABLE")),
	}
	return &item
}

func GetItemById(identificador string) (map[string]*dynamodb.AttributeValue, error) {
	result, err := GetDynamoInstance().GetItem(assembleItemForGetById(identificador))
	if err != nil {
		return nil, err
	}
	return result.Item, nil
}
func AssembleUserItem(result map[string]*dynamodb.AttributeValue) (model.Usuario, error) {
	var item model.Usuario

	err := dynamodbattribute.UnmarshalMap(result, &item)
	if err != nil {
		return item, err
	}
	return item, nil
}
func GenerateFilterForQueryUsers()(*dynamodb.ScanInput , error){
	filt := expression.Name("identifier").BeginsWith("users:")
	proj := expression.NamesList(expression.Name("identifier"), expression.Name("idade"), expression.Name("nome"), expression.Name("profissao"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(os.Getenv("TABLE")),
	}

	return params, nil
}

func AssembleUsersList()([]model.Usuario, error){
	listUsers := make([]model.Usuario,0)
	filterForTable, err := GenerateFilterForQueryUsers()
	if err != nil {
		return listUsers, err
	}
	var itemUser model.Usuario
	dynamoItems, err := GetDynamoInstance().Scan(filterForTable)
	for _, i := range dynamoItems.Items {
		itemUser = model.Usuario{}
	 	
		err = dynamodbattribute.UnmarshalMap(i, &itemUser)

		if err != nil {
			return listUsers, err
		}
		listUsers = append(listUsers, itemUser)
	}
	return listUsers, nil
}

func DeleteItemById(identifier string) (error) {
	_, err  := GetDynamoInstance().DeleteItem(assembleItemForDeleteById(identifier))
	return err
}