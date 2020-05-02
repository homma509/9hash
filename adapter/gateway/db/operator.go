package db

import "github.com/guregu/dynamo"

// TableOperator テーブル操作を表した構造体
type TableOperator struct {
	Client    *DynamoClient
	TableName string
}

// NewTableOperator TableOperatorの生成
func NewTableOperator(client *DynamoClient, tableName string) *TableOperator {
	return &TableOperator{
		Client:    client,
		TableName: tableName,
	}
}

// ConnectDB DBの接続
func (o *TableOperator) ConnectDB() (*dynamo.DB, error) {
	return o.Client.Connect()
}

// ConnectTable テーブルへの接続
func (o *TableOperator) ConnectTable() (*dynamo.Table, error) {
	return o.Client.ConnectTable(o.TableName)
}
