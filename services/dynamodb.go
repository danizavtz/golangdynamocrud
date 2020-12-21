package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"golangdynamocrud/model"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
	var itemUser = model.Usuario{}
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