package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/registry"
	"github.com/homma509/9hash/usecase"
)

// ValidatorPostSetting バリデーション設定
var ValidatorPostSetting = []*ValidatorSetting{
	{ArgName: "value", ValidateTags: "required"},
}

// RequestPostHash PostHashのリクエスト
type RequestPostHash struct {
	Value string `json:"value"`
}

// HashResponse レスポンス用のJSON形式を表した構造体
type HashResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostHash Hashの新規作成
func PostHashs(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	errs := ValidateBody(req.Body, ValidatorPostSetting)
	if errs != nil {
		return Response400(errs)
	}

	// JSON形式から構造体に変換
	var h RequestPostHash
	err := json.Unmarshal([]byte(req.Body), &h)
	if err != nil {
		return Response500(err)
	}

	// 新規作成処理
	creator := registry.GetFactory().BuildCreateHash()
	res, err := creator.Execute(&usecase.CreateHashRequest{
		Value: h.Value,
	})
	if err != nil {
		return Response500(err)
	}

	// 201レスポンス
	return Response201(res.HashID())
}
