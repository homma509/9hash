package db

import (
	"fmt"
	"reflect"
	"time"

	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

// DynamoResource DynamoDBのリソースを表した構造体
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

// DynamoModelMapper ModelをDynamoDbにマップする構造体
type DynamoModelMapper struct {
	Client    *ResourceTableOperator
	TableName string
	PKName    string
	SKName    string
}

// GetEntityNameFromStruct StructからEntity名を取得
func (d *DynamoModelMapper) GetEntityNameFromStruct(s interface{}) string {
	r := reflect.TypeOf(s)
	return r.Name()
}

// GetPK PKの取得
func (d *DynamoModelMapper) GetPK(r DynamoResource) string {
	return fmt.Sprintf("%s-%011d", r.EntityName(), r.ID())
}

// GetSK SKの取得
func (d *DynamoModelMapper) GetSK(r DynamoResource) string {
	return fmt.Sprintf("%011d", r.ID())
}

// GetEintityByID IDからEntityを取得
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

func (d *DynamoModelMapper) buildQueryCreates(rs []DynamoResource) (*dynamo.WriteTx, error) {
	table, err := d.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tx := d.Client.Client.Client.WriteTx()
	tx.Idempotent(true)

	for _, r := range rs {
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

		tx.Put(table.Put(r).If("attribute_not_exists(ID)"))
	}

	return tx, nil
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

// CreateResource Resourceの作成
func (d *DynamoModelMapper) CreateResource(r DynamoResource) error {
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

// CreateResources 複数のResourceの作成
func (d *DynamoModelMapper) CreateResources(rs []DynamoResource) error {
	tx, err := d.buildQueryCreates(rs)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// UpdateResource Resourceの更新
func (d *DynamoModelMapper) UpdateResource(r DynamoResource) error {
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

// DeleteResource Resourceの削除
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
