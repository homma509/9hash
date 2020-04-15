package db

import "time"

type DynamoResourceTime struct {
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

type DynamoResourceBase struct {
	Version int `dynamo:"Version"`
	DynamoResourceTime
}

// type Resource interface {
// 	SetPK() string
// 	SetPK() string
// }

type DynamoResourceSchema struct {
	PK string `dynamo:"PK,hash"`
	SK string `dynamo:"SK,range"`
}

type ResourceTableOperator struct {
	TableOperator
}

func NewResourceTableOperator(client *DynamoClient, tableName string) *ResourceTableOperator {
	return &ResourceTableOperator{
		TableOperator: *NewTableOperator(client, tableName),
	}
}
