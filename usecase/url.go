package usecase

// IGetURL URL取得UseCase
type IGetURL interface {
	Execute(req *GetURLRequest) (*GetURLResponse, error)
}

// GetURLRequest URL取得Request
type GetURLRequest struct {
	Key string
}

type GetURLResponse struct {
	URL string
}
