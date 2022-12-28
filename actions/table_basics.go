// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// snippet-start:[gov2.dynamodb.TableBasics.complete]
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
func (basics TableBasics) TableExists() (bool, error) {
	exists := true
	_, err := basics.DynamoDbClient.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(basics.TableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", basics.TableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", basics.TableName, err)
		}
		exists = false
	}
	return exists, err
}

// snippet-end:[gov2.dynamodb.DescribeTable]

// snippet-start:[gov2.dynamodb.CreateTable]

// CreateMovieTable creates a DynamoDB table with a composite primary key defined as
// a string sort key named `title`, and a numeric partition key named `year`.
// This function uses NewTableExistsWaiter to wait for the table to be created by
// DynamoDB before it returns.
// func (basics TableBasics) CreateMovieTable() (*types.TableDescription, error) {
// 	var tableDesc *types.TableDescription
// 	table, err := basics.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
// 		AttributeDefinitions: []types.AttributeDefinition{{
// 			AttributeName: aws.String("year"),
// 			AttributeType: types.ScalarAttributeTypeN,
// 		}, {
// 			AttributeName: aws.String("title"),
// 			AttributeType: types.ScalarAttributeTypeS,
// 		}},
// 		KeySchema: []types.KeySchemaElement{{
// 			AttributeName: aws.String("year"),
// 			KeyType:       types.KeyTypeHash,
// 		}, {
// 			AttributeName: aws.String("title"),
// 			KeyType:       types.KeyTypeRange,
// 		}},
// 		TableName: aws.String(basics.TableName),
// 		ProvisionedThroughput: &types.ProvisionedThroughput{
// 			ReadCapacityUnits:  aws.Int64(10),
// 			WriteCapacityUnits: aws.Int64(10),
// 		},
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't create table %v. Here's why: %v\n", basics.TableName, err)
// 	} else {
// 		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
// 		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
// 			TableName: aws.String(basics.TableName)}, 5*time.Minute)
// 		if err != nil {
// 			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
// 		}
// 		tableDesc = table.TableDescription
// 	}
// 	return tableDesc, err
// }

// CreateQueryTable creates a DynamoDB table with a composite primary key defined as
// a string partition key named `query` and a numeric sort key named `datetime`.
// This function uses NewTableExistsWaiter to wait for the table to be created by
// DynamoDB before it returns.
// This table uses on-demand pay-per-request billing.
func (basics TableBasics) CreateQueryTable() (*types.TableDescription, error) {
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
		log.Printf("Couldn't create table %v. Here's why: %v\n", basics.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(basics.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

// snippet-end:[gov2.dynamodb.CreateTable]

// snippet-start:[gov2.dynamodb.ListTables]

// ListTables lists the DynamoDB table names for the current account.
func (basics TableBasics) ListTables() ([]string, error) {
	var tableNames []string
	tables, err := basics.DynamoDbClient.ListTables(
		context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Couldn't list tables. Here's why: %v\n", err)
	} else {
		tableNames = tables.TableNames
	}
	return tableNames, err
}

// snippet-end:[gov2.dynamodb.ListTables]

// snippet-start:[gov2.dynamodb.PutItem]

func (basics TableBasics) AddQuery(queryResponse QueryResponse) error {
	item, err := attributevalue.MarshalMap(queryResponse)
	if err != nil {
		panic(err)
	}
	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

// // AddMovie adds a movie the DynamoDB table.
// func (basics TableBasics) AddMovie(movie Movie) error {
// 	item, err := attributevalue.MarshalMap(movie)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
// 		TableName: aws.String(basics.TableName), Item: item,
// 	})
// 	if err != nil {
// 		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
// 	}
// 	return err
// }

// snippet-end:[gov2.dynamodb.PutItem]

// snippet-start:[gov2.dynamodb.UpdateItem]

// UpdateMovie updates the rating and plot of a movie that already exists in the
// DynamoDB table. This function uses the `expression` package to build the update
// expression.
func (basics TableBasics) UpdateMovie(movie Movie) (map[string]map[string]interface{}, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]map[string]interface{}
	update := expression.Set(expression.Name("info.rating"), expression.Value(movie.Info["rating"]))
	update.Set(expression.Name("info.plot"), expression.Value(movie.Info["plot"]))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err = basics.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(basics.TableName),
			Key:                       movie.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update movie %v. Here's why: %v\n", movie.Title, err)
		} else {
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
			}
		}
	}
	return attributeMap, err
}

// snippet-end:[gov2.dynamodb.UpdateItem]
