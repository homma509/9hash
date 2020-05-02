package db

import "time"

// DynamoResourceTime Resourceの時間要素を表した構造体
type DynamoResourceTime struct {
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

// DynamoResourceBase Resourceの基本要素を表した構造体
type DynamoResourceBase struct {
	Version uint64 `dynamo:"Version"`
	DynamoResourceTime
}

// type Resource interface {
// 	SetPK() string
// 	SetPK() string
// }

// DynamoResourceSchema Resourceのスキーマ要素を表した構造体
type DynamoResourceSchema struct {
	PK string `dynamo:"PK,hash"`
	SK string `dynamo:"SK,range"`
}

// ResourceTableOperator リソース操作を表した構造体
type ResourceTableOperator struct {
	TableOperator
}

// NewResourceTableOperator ResourceTableOperatorの生成
func NewResourceTableOperator(client *DynamoClient, tableName string) *ResourceTableOperator {
	return &ResourceTableOperator{
		TableOperator: *NewTableOperator(client, tableName),
	}
}
