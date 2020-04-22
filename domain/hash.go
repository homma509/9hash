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
	GetHashByID(uint64) (*HashModel, error)
	GetHashByKey(string) (*HashModel, error)
	CreateHash(*HashModel) (*HashModel, error)
	CreateHashs([]*HashModel) ([]*HashModel, error)
	UpdateHash(*HashModel) (*HashModel, error)
	DeleteHash(*HashModel) error
}
