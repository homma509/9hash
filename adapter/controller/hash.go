package controller

import (
	"encoding/json"
	"strconv"

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

// PostHashs Hashの新規作成
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

func PutHash(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
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

	// パスパラメータからHashIDを取得する
	id, err := strconv.ParseUint(req.PathParameters["hash_id"], 10, 64)
	if err != nil {
		return Response500(err)
	}

	// 更新処理
	updator := registry.GetFactory().BuildUpdateHash()
	_, err = updator.Execute(&usecase.UpdateHashRequest{
		ID:    id,
		Value: h.Value,
	})
	if err != nil {
		return Response500(err)
	}

	return Response200()
}
