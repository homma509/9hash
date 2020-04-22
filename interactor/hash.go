package interactor

import (
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/usecase"
	"github.com/pkg/errors"
)

// HashGetter Hash取得
type HashGetter struct {
	HashRepository domain.HashRepository
}

// NewHashGetter Hash取得をインスタンス生成します
func NewHashGetter(rep domain.HashRepository) *HashGetter {
	return &HashGetter{
		HashRepository: rep,
	}
}

// Execute Hashを取得
func (h *HashGetter) Execute(req *usecase.GetHashRequest) (*usecase.GetHashResponse, error) {
	hash, err := h.HashRepository.GetHashByID(req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.GetHashResponse{Hash: hash}, nil
}

// NewHashsGetter Hashs取得
type HashsGetter struct {
	HashRepository domain.HashRepository
}

// NewHashsGetter Hashs取得をインスタンス生成します
func NewHashsGetter(rep domain.HashRepository) *HashsGetter {
	return &HashsGetter{
		HashRepository: rep,
	}
}

// Execute Hashsを取得
func (h *HashsGetter) Execute(req *usecase.GetHashsRequest) (*usecase.GetHashsResponse, error) {
	hashs, err := h.HashRepository.GetHashs()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.GetHashsResponse{Hashs: hashs}, nil
}

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

// HashsCreator Hashs新規作成
type HashsCreator struct {
	HashRepository domain.HashRepository
}

// NewHashsCreator Hashs新規作成をインスタンス生成します
func NewHashsCreator(rep domain.HashRepository) *HashsCreator {
	return &HashsCreator{
		HashRepository: rep,
	}
}

// Execute Hashsを新規作成
func (h *HashsCreator) Execute(req *usecase.CreateHashsRequest) (*usecase.CreateHashsResponse, error) {
	hashs, err := h.HashRepository.CreateHashs(req.ToHashsModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.CreateHashsResponse{Hashs: hashs}, nil
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
	hash, err := h.HashRepository.UpdateHash(req.ToHashModel())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.UpdateHashResponse{Hash: hash}, nil
}

// HashDeleter Hash削除
type HashDeleter struct {
	HashRepository domain.HashRepository
}

// NewHashDeleter Hash削除をインスタンス生成します
func NewHashDeleter(rep domain.HashRepository) *HashDeleter {
	return &HashDeleter{
		HashRepository: rep,
	}
}

func (h *HashDeleter) Execute(req *usecase.DeleteHashRequest) (*usecase.DeleteHashResponse, error) {
	hash, err := h.HashRepository.GetHashByID(req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = h.HashRepository.DeleteHash(hash)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &usecase.DeleteHashResponse{}, nil
}
