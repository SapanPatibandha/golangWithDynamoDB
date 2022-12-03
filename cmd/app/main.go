package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/SapanPatibandha/golangWithDynamoDB/internal/reposotery/adaptor"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/reposotery/instance"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/rules"
	"github.com/SapanPatibandha/golangWithDynamoDB/utils/logger"
)

func main() {
	fmt.Println("Welcome")

	configs := Config.GetConfig()

	connection := instance.GetConnection()

	repository := adaptor.NewAdaptor(connection)

	logger.INFO("waiting for service to start")

	errors := Migrate(connection)

	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("error on migration.. ", err)
		}
	}

	logger.PANIC("", checkTables(connection))

	port := fmp.Springf(":%v", configs.port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("service is running on port", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error
	callMigrateAndAppendError(&errors, connection, &RulesProduct.Rules{})
	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, rule rules.Interface) {
	err := rule.Migrate(connection)

	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(conneciton *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTables(&dynamodb.ListTablesInput{}))

	if response != nil {
		if len(response.tables) == 0 {
			logger.INFO("tables not found", nil)
		}

		for _, tableName := range response.tableName {
			logger.INFO("table found:", tableName)
		}
	}

	return err
}
