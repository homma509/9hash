package db

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

type DynamoResource interface {
	EntityName() string
	PK() string
	SetPK()
	SK() string
	SetSK()
	ID() uint64
	SetID(id uint64)
	Version() uint64
	SetVersion(v uint64)
	CreatedAt() time.Time
	SetCreatedAt(t time.Time)
	UpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

type DynamoModelMapper struct {
	Client    *ResourceTableOperator
	TableName string
	PKName    string
	SKName    string
}

func (d *DynamoModelMapper) GetEntityNameFromStruct(s interface{}) string {
	r := reflect.TypeOf(s)
	return r.Name()
}

func (d *DynamoModelMapper) PutResource(r DynamoResource) error {
	if d.isNewEntity(r) {
		return d.createResource(r)
	}
	return d.updateResource(r)
}

func (d *DynamoModelMapper) GetPK(r DynamoResource) string {
	return fmt.Sprintf("%s-%011d", r.EntityName(), r.ID())
}

func (d *DynamoModelMapper) GetSK(r DynamoResource) string {
	return fmt.Sprintf("%011d", r.ID())
}
func (d *DynamoModelMapper) GetEintityByID(id uint64, r DynamoResource, ret interface{}) (interface{}, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.SetID(id)
	err = table.
		Get(d.PKName, r.PK()).
		Range(d.SKName, dynamo.Equal, r.SK()).
		One(ret)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ret, nil
}

func (d *DynamoModelMapper) buildQueryCreate(r DynamoResource) (*dynamo.Put, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := d.generateID(r.EntityName())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.SetCreatedAt(time.Now())
	r.SetUpdatedAt(time.Now())
	r.SetID(id)
	r.SetVersion(1)
	r.SetPK()
	r.SetSK()

	query := table.Put(r).If("attribute_not_exists(ID)")

	return query, nil
}

func (d *DynamoModelMapper) buildQueryUpdate(r DynamoResource) (*dynamo.Put, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	oldVersion := r.Version()

	r.SetUpdatedAt(time.Now())
	r.SetVersion(oldVersion + 1)

	query := table.Put(r).If("Version", dynamo.Equal, oldVersion)

	return query, nil
}

func (d *DynamoModelMapper) buildQueryDelete(r DynamoResource) (*dynamo.Delete, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query := table.
		Delete(d.PKName, r.PK()).
		Range(d.SKName, r.SK())

	return query, nil
}

func (d *DynamoModelMapper) createResource(r DynamoResource) error {
	query, err := d.buildQueryCreate(r)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *DynamoModelMapper) updateResource(r DynamoResource) error {
	query, err := d.buildQueryUpdate(r)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *DynamoModelMapper) DeleteResource(r DynamoResource) error {
	query, err := d.buildQueryDelete(r)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *DynamoModelMapper) isNewEntity(r DynamoResource) bool {
	return r.Version() == 0
}

func (d *DynamoModelMapper) generateID(tableName string) (uint64, error) {
	attr, err := d.atomicCount(fmt.Sprintf("AtomicCounter-%s", tableName), "AtomicCounter", "CurrentNumber", 1)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	s := aws.StringValue(attr.N)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

func (d *DynamoModelMapper) atomicCount(pk, sk, counterName string, value uint64) (*dynamodb.AttributeValue, error) {
	db, err := d.Client.ConnectDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Printf("テーブル名：%s", d.TableName)

	// TODO dynamo.DBを使って簡易化する
	output, err := db.Client().UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(pk),
			},
			"SK": {
				S: aws.String(sk),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			counterName: {
				Action: aws.String("ADD"),
				Value: &dynamodb.AttributeValue{
					N: aws.String(fmt.Sprintf("%d", value)),
				},
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueUpdatedNew),
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return output.Attributes[counterName], nil
}
