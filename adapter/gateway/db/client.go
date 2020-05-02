package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
