package controller

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/registry"
	"github.com/homma509/9hash/usecase"
)

// GetURL URLの取得
func GetURL(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからKeyを取得する
	key := req.PathParameters["key"]
	if key == "" {
		return Response500(domain.ErrBadRequest)
	}

	// 取得処理
	getter := registry.GetFactory().BuildGetURL()
	res, err := getter.Execute(&usecase.GetURLRequest{
		Key: key,
	})
	if err != nil {
		if err.Error() == domain.ErrNotFound.Error() {
			return Response404()
		}
		return Response500(err)
	}

	// 308レスポンス
	return Response308(res.URL)
}
