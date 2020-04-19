package domain

// HashModel Haasのモデル
type HashModel struct {
	ID    uint64
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
	GetHashs() ([]*HashModel, error)
	GetHashByID(id uint64) (*HashModel, error)
	GetHashByKey(key string) (*HashModel, error)
	CreateHash(h *HashModel) (*HashModel, error)
	UpdateHash(h *HashModel) error
	DeleteHash(h *HashModel) error
}
