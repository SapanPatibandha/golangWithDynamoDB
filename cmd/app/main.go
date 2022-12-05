package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"honnef.co/go/tools/config"

	"github.com/SapanPatibandha/golangWithDynamoDB/config"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/reposotery/adaptor"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/reposotery/instance"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/routes"
	RulesProduct "github.com/SapanPatibandha/golangWithDynamoDB/internal/routes/product"
	"github.com/SapanPatibandha/golangWithDynamoDB/internal/rules"
	"github.com/SapanPatibandha/golangWithDynamoDB/utils/logger"
)

func main() {
	fmt.Println("Welcome")

	configs := config.GetConfig()

	connection := instance.GetConnection()
	repository := adaptor.NewAdaptor(connection)

	logger.INFO("waiting for service to start", nil)

	errors := Migrate(connection)

	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("error on migration.. ", err)
		}
	}

	logger.PANIC("", checkTables(connection))

	port := fmt.Sprintf(":%v", configs.port)
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
	response, err := conneciton.ListTables(&dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("tables not found", nil)
		}

		for _, tableName := range response.TableNames {
			logger.INFO("table found:", tableName)
		}
	}

	return err
}
