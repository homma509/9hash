package db

import "github.com/guregu/dynamo"

type TableOperator struct {
	Client    *DynamoClient
	TableName string
}

func NewTableOperator(client *DynamoClient, tableName string) *TableOperator {
	return &TableOperator{
		Client:    client,
		TableName: tableName,
	}
}

func (o *TableOperator) ConnectDB() (*dynamo.DB, error) {
	return o.Client.Connect()
}

func (o *TableOperator) ConnectTable() (*dynamo.Table, error) {
	return o.Client.ConnectTable(o.TableName)
}
