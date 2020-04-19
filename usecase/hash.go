package usecase

import (
	"github.com/homma509/9hash/domain"
	"github.com/teris-io/shortid"
)

// IGetHash Hash取得UseCase
type IGetHash interface {
	Execute(req *GetHashRequest) (*GetHashResponse, error)
}

// GetHashRequest Hash取得Request
type GetHashRequest struct {
	ID uint64
}

// GetHashResponse Hash取得Response
type GetHashResponse struct {
	Hash *domain.HashModel
}

// IGetHashs Hash一覧取得UseCase
type IGetHashs interface {
	Execute(req *GetHashsRequest) (*GetHashsResponse, error)
}

// GetHashsRequest Hash一覧取得Request
type GetHashsRequest struct {
}

// GetHashsResponse Hash一覧取得Response
type GetHashsResponse struct {
	Hashs []*domain.HashModel
}

// ICreateHash Hash新規作成UseCase
type ICreateHash interface {
	Execute(req *CreateHashRequest) (*CreateHashResponse, error)
}

// CreateHashRequest Hash新規作成Request
type CreateHashRequest struct {
	Value string
}

// CreateHashResponse Hash新規作成Response
type CreateHashResponse struct {
	Hash *domain.HashModel
}

// ToHashModel Hash新規作成RequestをHashモデルに変換します
func (h *CreateHashRequest) ToHashModel() *domain.HashModel {
	return domain.NewHashModel(hashKey(), h.Value)
}

// IUpdateHash Hash更新UseCase
type IUpdateHash interface {
	Execute(req *UpdateHashRequest) (*UpdateHashResponse, error)
}

// UpdateHashRequest Hash更新Request
type UpdateHashRequest struct {
	ID    uint64
	Value string
}

// UpdateHashResponse Hash更新Response
type UpdateHashResponse struct {
	Hash *domain.HashModel
}

func (h *UpdateHashRequest) ToHashModel() *domain.HashModel {
	return &domain.HashModel{
		ID:    h.ID,
		Value: h.Value,
	}
}

func hashKey() string {
	key := shortid.MustGenerate()
	for key == "shorten" {
		key = shortid.MustGenerate()
	}
	return key
}
