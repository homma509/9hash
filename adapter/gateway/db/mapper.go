package db

import (
	"fmt"
	"reflect"
	"time"

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

	query := table.Put(r).If("Version = ?", oldVersion)

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
	n, err := d.atomicCount(tableName)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

func (d *DynamoModelMapper) atomicCount(tableName string) (uint64, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	type count struct {
		CurrentNumber uint64
	}

	var result count
	err = table.
		Update("PK", fmt.Sprintf("AtomicCounter-%s", tableName)).
		Range("SK", "AtomicCounter").
		Add("CurrentNumber", 1).
		Value(&result)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return result.CurrentNumber, nil
}
