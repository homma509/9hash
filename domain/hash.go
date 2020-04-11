package domain

type HashModel struct {
	Key   string
	Value string
}

func NewHashModel(key, value string) *HashModel {
	return &HashModel{
		Key:   key,
		Value: value,
	}
}

type HashRepository interface {
	CreateHash(h *HashModel) (*HashModel, error)
	DeleteHash(id uint64) error
	GetHashByKey(key string) (*HashModel, error)
}
