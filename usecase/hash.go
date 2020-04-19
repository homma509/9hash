package usecase

import (
	"github.com/homma509/9hash/domain"
	"github.com/teris-io/shortid"
)

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

func (h *CreateHashResponse) HashID() uint64 {
	return h.Hash.ID
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
