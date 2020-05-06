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

// ConnectTable テーブルへの接続
func (o *TableOperator) ConnectTable() (*dynamo.Table, error) {
	return o.Client.ConnectTable(o.TableName)
}

// CreateTestTable テスト用テーブルの作成
func (o *TableOperator) CreateTestTable(schema interface{}) error {
	return o.Client.CreateTestTable(o.TableName, schema)
}

// DropTable テーブルの削除
func (o *TableOperator) DropTable() error {
	return o.Client.DropTable(o.TableName)
}
