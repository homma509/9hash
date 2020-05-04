package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

// DynamoClient DynamoDBのClientの構造体
type DynamoClient struct {
	Client *dynamo.DB
	Config *aws.Config
}

// NewClient DynamoClientの生成
func NewClient(config *aws.Config) *DynamoClient {
	return &DynamoClient{Config: config}
}

// Connect DynamoDBの接続
func (c *DynamoClient) Connect() (*dynamo.DB, error) {
	if c.Client == nil {
		sess, err := session.NewSession(c.Config)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		c.Client = dynamo.New(sess)
	}
	return c.Client, nil
}

// ConnectTable DyanamoDBのテーブルへの接続
func (c *DynamoClient) ConnectTable(tableName string) (*dynamo.Table, error) {
	db, err := c.Connect()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	table := db.Table(tableName)

	return &table, nil
}

// CreateTestTable テスト用テーブルの作成
func (c *DynamoClient) CreateTestTable(tableName string, table interface{}) error {
	db, err := c.Connect()
	if err != nil {
		return errors.WithStack(err)
	}

	err = db.CreateTable(tableName, table).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DropTable テーブルの削除
func (c *DynamoClient) DropTable(tableName string) error {
	db, err := c.Connect()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.Client().DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
