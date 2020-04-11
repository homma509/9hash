package domain

// HashModel Haasのモデル
type HashModel struct {
	Key   string
	Value string
}

// NewHashModel Hashモデルのインスタンスを生成します
func NewHashModel(key, value string) *HashModel {
	return &HashModel{
		Key:   key,
		Value: value,
	}
}

// HashRepository Hashモデルのリポジトリ
type HashRepository interface {
	CreateHash(h *HashModel) (*HashModel, error)
	DeleteHash(id uint64) error
	GetHashByKey(key string) (*HashModel, error)
}
