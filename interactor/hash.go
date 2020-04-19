package interactor

import (
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/usecase"
	"github.com/pkg/errors"
)

// HashCreator Hash新規作成
type HashCreator struct {
	HashRepository domain.HashRepository
}

// NewHashCreator Hash新規作成をインスタンス生成します
func NewHashCreator(rep domain.HashRepository) *HashCreator {
	return &HashCreator{
		HashRepository: rep,
	}
}

// Execute Hashを新規作成
func (h *HashCreator) Execute(req *usecase.CreateHashRequest) (*usecase.CreateHashResponse, error) {
	hash, err := h.HashRepository.CreateHash(req.ToHashModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.CreateHashResponse{Hash: hash}, nil
}

// HashUpdator Hash更新
type HashUpdator struct {
	HashRepository domain.HashRepository
}

// NewHashUpdator Hash更新をインスタンス生成します
func NewHashUpdator(rep domain.HashRepository) *HashUpdator {
	return &HashUpdator{
		HashRepository: rep,
	}
}

func (h *HashUpdator) Execute(req *usecase.UpdateHashRequest) (*usecase.UpdateHashResponse, error) {
	err := h.HashRepository.UpdateHash(req.ToHashModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.UpdateHashResponse{}, nil
}
