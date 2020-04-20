package interactor

import (
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/usecase"
	"github.com/pkg/errors"
)

// URLGetter URL取得
type URLGetter struct {
	HashRepository domain.HashRepository
}

// NewURLGetter URL取得をインスタンス生成します
func NewURLGetter(rep domain.HashRepository) *URLGetter {
	return &URLGetter{
		HashRepository: rep,
	}
}

// Execute URLを取得
func (u *URLGetter) Execute(req *usecase.GetURLRequest) (*usecase.GetURLResponse, error) {
	hash, err := u.HashRepository.GetHashByKey(req.Key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &usecase.GetURLResponse{URL: hash.Value}, nil
}
