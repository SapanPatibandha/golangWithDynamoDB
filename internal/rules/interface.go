package rules

import (
	"io"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Interface interface {
	ConvertIoReturnToStruct(data io.Reader, mode interface{}) (body interface{}, err error)
	GetMock() interface{}
	Migrate(connection *dynamodb.DynamoDB) error
	Validate(model interface{}) error
}
