package adapter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Database struct {
	connection *dynamodb.DynamoDB
	logMode    bool
}

type Interface interface{}

func NewAdapter() Interface {}

func (db *Database) Health() bool {}

func (db *Database) FindAll() {}

func (db *Database) FindOne(condition map[string]interface{}, tableName string) (response *dynamodb.GetItemOutput, err error) {

	conditionParsed, err := dynamodbattribute.MarshalMap(condition)

	if err != nil {
		return nil, err
	}

	inptu := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       conditionParsed,
	}

	return db.connection.GetItem(inptu)
}

func (db *Database) CreateOrUpdate() {}

func (db *Database) Delete() {}
