// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	// "log"
	"time"

	"github.com/moov-io/base/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"	
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// snippet-start:[gov2.dynamodb.TableBasics.struct]

// TableBasics encapsulates the Amazon DynamoDB service actions used in the examples.
// It contains a DynamoDB service client that is used to act on the specified table.
type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

// snippet-end:[gov2.dynamodb.TableBasics.struct]

// snippet-start:[gov2.dynamodb.DescribeTable]

// TableExists determines whether a DynamoDB table exists.
func (basics TableBasics) TableExists(logger log.Logger) (bool, error) {
	exists := true
	_, err := basics.DynamoDbClient.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(basics.TableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			logger.Logf("Table %v does not exist.\n", basics.TableName)
			err = nil
		} else {
			logger.Logf("Couldn't determine existence of table %v. Here's why: %v\n", basics.TableName, err)
		}
		exists = false
	}
	return exists, err
}

// snippet-end:[gov2.dynamodb.DescribeTable]

// CreateQueryTable creates a DynamoDB table with a composite primary key defined as
// a string partition key named `query` and a numeric sort key named `datetime`.
// This function uses NewTableExistsWaiter to wait for the table to be created by
// DynamoDB before it returns.
// This table uses on-demand pay-per-request billing.
func (basics TableBasics) CreateQueryTable(logger log.Logger) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := basics.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("query"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("datetime"),
			AttributeType: types.ScalarAttributeTypeN,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("query"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("datetime"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName: aws.String(basics.TableName),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		logger.Logf("Couldn't create table %v. Here's why: %v\n", basics.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(basics.TableName)}, 5*time.Minute)
		if err != nil {
			logger.Logf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

// snippet-end:[gov2.dynamodb.CreateTable]

// snippet-start:[gov2.dynamodb.PutItem]

func (basics TableBasics) AddQuery(logger log.Logger, inSearchResponse searchResponse) error {
	item, err := attributevalue.MarshalMap(inSearchResponse)
	if err != nil {
		panic(err)
	}
	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		logger.Logf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

// snippet-end:[gov2.dynamodb.PutItem]

// Create a function that will establish a connection to DynamoDB, create a table if it does not already exist, and add an item to the table.
func WriteQueryToDB(logger log.Logger, inSearchResponse searchResponse) error {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.Logf("unable to load SDK config, %v", err)
	}
	tableName := "watchman-queries"
	tableBasics := TableBasics{TableName: tableName,
		DynamoDbClient: dynamodb.NewFromConfig(sdkConfig)}

	exists, err := tableBasics.TableExists(logger)
	if err != nil {
		panic(err)
	}
	if !exists {
		logger.Logf("Creating table %v... \n", tableName)
		_, err = tableBasics.CreateQueryTable(logger)
		if err != nil {
			panic(err)
		} else {
			logger.Logf("Created table %v. \n", tableName)
		}
	} else {
		logger.Logf("Table %v already exists.\n", tableName)
	}

    logger.Logf("Attempting to write search response to table %v... \n", tableName)

    err = tableBasics.AddQuery(logger, inSearchResponse)

    if err != nil {
		panic(err)
    }

	logger.Logf("Successfully wrote search response to table %v... \n", tableName)

	return err
}