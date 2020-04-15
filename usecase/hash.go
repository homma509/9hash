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

// ToHashModel Hash新規作成RequestをHashモデルに変換します
func (h *CreateHashRequest) ToHashModel() *domain.HashModel {
	return domain.NewHashModel(hashKey(), h.Value)
}

// CreateHashResponse Hash新規作成Response
type CreateHashResponse struct {
	Hash *domain.HashModel
}

func (h *CreateHashResponse) HashID() string {
	return h.Hash.ID
}

func hashKey() string {
	key := shortid.MustGenerate()
	for key == "shorten" {
		key = shortid.MustGenerate()
	}
	return key
}
