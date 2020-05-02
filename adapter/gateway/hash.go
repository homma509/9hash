package gateway

import (
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

// NewHashResource HashResourceの生成
func NewHashResource(h *domain.HashModel, m *db.DynamoModelMapper) *HashResource {
	return &HashResource{
		HashModel: *h,
		Mapper:    m,
	}
}

// EntityName Entity名の取得
func (h *HashResource) EntityName() string {
	return h.Mapper.GetEntityNameFromStruct(*h)
}

// PK PKの取得
func (h *HashResource) PK() string {
	return h.Mapper.GetPK(h)
}

// SetPK PKの設定
func (h *HashResource) SetPK() {
	h.DynamoResourceSchema.PK = h.PK()
}

// SK SKの取得
func (h *HashResource) SK() string {
	return h.Mapper.GetSK(h)
}

// SetSK SKの設定
func (h *HashResource) SetSK() {
	h.DynamoResourceSchema.SK = h.SK()
}

// ID IDの取得
func (h *HashResource) ID() uint64 {
	return h.HashModel.ID
}

// SetID IDの設定
func (h *HashResource) SetID(id uint64) {
	h.HashModel.ID = id
}

// Version Versionの取得
func (h *HashResource) Version() uint64 {
	return h.DynamoResourceBase.Version
}

// SetVersion Versionの設定
func (h *HashResource) SetVersion(v uint64) {
	h.DynamoResourceBase.Version = v
}

// CreatedAt 作成日時の取得
func (h *HashResource) CreatedAt() time.Time {
	return h.DynamoResourceBase.CreatedAt
}

// SetCreatedAt 作成日時の設定
func (h *HashResource) SetCreatedAt(t time.Time) {
	h.DynamoResourceBase.CreatedAt = t
}

// UpdatedAt 更新日時の取得
func (h *HashResource) UpdatedAt() time.Time {
	return h.DynamoResourceBase.UpdatedAt
}

// SetUpdatedAt 更新日時の設定
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

// GetHashByID IDよりHashを取得
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

// GetHashByKey KeyよりHashを取得
func (h *HashOperator) GetHashByKey(key string) (*domain.HashModel, error) {
	table, err := h.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hashs []HashResource
	err = table.Scan().
		Filter("'Key' = ?", key).
		Filter("begins_with('PK', ?)", h.Mapper.GetEntityNameFromStruct(HashResource{})).
		All(&hashs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(hashs) == 0 {
		return nil, errors.WithStack(domain.ErrNotFound)
	}

	return &hashs[0].HashModel, nil
}

// GetHashs 全てのHashを取得
func (h *HashOperator) GetHashs() ([]*domain.HashModel, error) {
	table, err := h.Client.ConnectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var hashResources []HashResource
	err = table.Scan().
		Filter("begins_with('PK', ?)", h.Mapper.GetEntityNameFromStruct(HashResource{})).
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

// CreateHash Hashの作成
func (h *HashOperator) CreateHash(m *domain.HashModel) (*domain.HashModel, error) {
	r := NewHashResource(m, h.Mapper)

	err := h.Mapper.CreateResource(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &r.HashModel, nil
}

// CreateHashs 複数のHashの作成
func (h *HashOperator) CreateHashs(ms []*domain.HashModel) ([]*domain.HashModel, error) {
	rs := make([]*HashResource, len(ms))
	for i, m := range ms {
		rs[i] = NewHashResource(m, h.Mapper)
	}

	err := h.Mapper.CreateResources(ToDynamoResources(rs))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hashs := make([]*domain.HashModel, len(rs))
	for i, r := range rs {
		hashs[i] = &r.HashModel
	}

	return hashs, nil
}

// ToDynamoResources HashResourceをDynamoResourceに変換
func ToDynamoResources(rs []*HashResource) []db.DynamoResource {
	ds := make([]db.DynamoResource, len(rs))
	for i, r := range rs {
		ds[i] = r
	}
	return ds
}

// UpdateHash Hashの更新
func (h *HashOperator) UpdateHash(m *domain.HashModel) (*domain.HashModel, error) {
	r, err := h.getUserResourceByID(m.ID)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return nil, errors.WithStack(domain.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	r.Value = m.Value

	err = h.Mapper.UpdateResource(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &r.HashModel, nil
}

// DeleteHash Hashの削除
func (h *HashOperator) DeleteHash(m *domain.HashModel) error {
	r := NewHashResource(m, h.Mapper)

	err := h.Mapper.DeleteResource(r)
	if err != nil {
		if err.Error() == dynamo.ErrNotFound.Error() {
			return errors.WithStack(domain.ErrNotFound)
		}
		return errors.WithStack(err)
	}

	return nil
}
