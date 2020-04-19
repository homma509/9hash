package gateway

import (
	"fmt"
	"time"

	"github.com/guregu/dynamo"
	"github.com/homma509/9hash/adapter/gateway/db"
	"github.com/homma509/9hash/domain"
	"github.com/pkg/errors"
)

// HashResource DynamoDb上のデータ構造を表した構造体
type HashResource struct {
	db.DynamoResourceSchema
	db.DynamoResourceBase
	domain.HashModel
	Mapper *db.DynamoModelMapper `dynamo:"-"`
}

func NewHashResource(h *domain.HashModel, m *db.DynamoModelMapper) *HashResource {
	return &HashResource{
		HashModel: *h,
		Mapper:    m,
	}
}

func (h *HashResource) EntityName() string {
	return h.Mapper.GetEntityNameFromStruct(*h)
}

func (h *HashResource) PK() string {
	return h.Mapper.GetPK(h)
}

func (h *HashResource) SetPK() {
	h.DynamoResourceSchema.PK = h.PK()
}

func (h *HashResource) SK() string {
	return h.Mapper.GetSK(h)
}

func (h *HashResource) SetSK() {
	h.DynamoResourceSchema.SK = h.SK()
}

func (h *HashResource) ID() uint64 {
	return h.HashModel.ID
}

func (h *HashResource) SetID(id uint64) {
	h.HashModel.ID = id
}

func (h *HashResource) Version() uint64 {
	return h.DynamoResourceBase.Version
}

func (h *HashResource) SetVersion(v uint64) {
	h.DynamoResourceBase.Version = v
}

func (h *HashResource) CreatedAt() time.Time {
	return h.DynamoResourceBase.CreatedAt
}

func (h *HashResource) SetCreatedAt(t time.Time) {
	h.DynamoResourceBase.CreatedAt = t
}

func (h *HashResource) UpdatedAt() time.Time {
	return h.DynamoResourceBase.UpdatedAt
}

func (h *HashResource) SetUpdatedAt(t time.Time) {
	h.DynamoResourceBase.UpdatedAt = t
}

// HashOperator Hashを操作する構造体
type HashOperator struct {
	Client *db.ResourceTableOperator
	Mapper *db.DynamoModelMapper
}

func (h *HashOperator) getUserResourceByID(id uint64) (*HashResource, error) {
	var r HashResource
	_, err := h.Mapper.GetEintityByID(id, &HashResource{}, &r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &r, nil
}

func (h *HashOperator) GetHashByID(id uint64) (*domain.HashModel, error) {
	r, err := h.getUserResourceByID(id)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return nil, errors.WithStack(domain.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return &r.HashModel, nil
}

func (h *HashOperator) GetHashByKey(key string) (*domain.HashModel, error) {
	table, err := h.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hashs []HashResource
	err = table.Scan().
		Filter("Key", dynamo.Equal, key).
		Filter(fmt.Sprintf("begins_with('%s', ?)", h.Mapper.GetEntityNameFromStruct(HashResource{}))).
		All(&hashs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(hashs) == 0 {
		return nil, errors.WithStack(domain.ErrNotFound)
	}

	return &hashs[0].HashModel, nil
}

func (h *HashOperator) GetHashs() ([]*domain.HashModel, error) {
	table, err := h.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hashResources []HashResource
	err = table.Scan().
		Filter(fmt.Sprintf("begins_with('%s', ?)", h.Mapper.GetEntityNameFromStruct(HashResource{}))).
		All(&hashResources)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hashs = make([]*domain.HashModel, len(hashResources))
	for i := range hashResources {
		hashs[i] = &hashResources[i].HashModel
	}

	return hashs, nil
}

func (h *HashOperator) CreateHash(m *domain.HashModel) (*domain.HashModel, error) {
	r := NewHashResource(m, h.Mapper)

	err := h.Mapper.PutResource(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &r.HashModel, nil
}

func (h *HashOperator) UpdateHash(m *domain.HashModel) (*domain.HashModel, error) {
	r, err := h.getUserResourceByID(m.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.Value = m.Value

	err = h.Mapper.PutResource(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &r.HashModel, nil
}

func (h *HashOperator) DeleteHash(m *domain.HashModel) error {
	r := NewHashResource(m, h.Mapper)

	err := h.Mapper.DeleteResource(r)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
