package controller

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/registry"
	"github.com/homma509/9hash/usecase"
)

// ValidatorPostSetting バリデーション設定
var ValidatorPostSetting = []*ValidatorSetting{
	{ArgName: "values", ValidateTags: "required"},
}

// RequestPostHashs PostHashsのリクエスト
type RequestPostHashs struct {
	Values []string `json:"values"`
}

// ValidatorPutSetting バリデーション設定
var ValidatorPutSetting = []*ValidatorSetting{
	{ArgName: "value", ValidateTags: "required"},
}

// RequestPutHash PutHashのリクエスト
type RequestPutHash struct {
	Value string `json:"value"`
}

// GetHash Hashの取得
func GetHash(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからHashIDを取得する
	id, err := strconv.ParseUint(req.PathParameters["hash_id"], 10, 64)
	if err != nil {
		return Response500(err)
	}

	// 取得処理
	getter := registry.GetFactory().BuildGetHash()
	res, err := getter.Execute(&usecase.GetHashRequest{
		ID: id,
	})
	if err != nil {
		if err.Error() == domain.ErrNotFound.Error() {
			return Response404()
		}
		return Response500(err)
	}

	// 200レスポンス
	return Response200(res.Hash)
}

// GetHashs Hashsの取得
func GetHashs(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// 取得処理
	getter := registry.GetFactory().BuildGetHashs()
	res, err := getter.Execute(&usecase.GetHashsRequest{})
	if err != nil {
		return Response500(err)
	}

	// 200レスポンス
	return Response200(res.Hashs)
}

// PostHashs 複数Hashの新規作成
func PostHashs(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	errs := ValidateBody(req.Body, ValidatorPostSetting)
	if errs != nil {
		return Response400(errs)
	}

	// JSON形式から構造体に変換
	var h RequestPostHashs
	err := json.Unmarshal([]byte(req.Body), &h)
	if err != nil {
		return Response500(err)
	}

	// 新規作成処理
	creator := registry.GetFactory().BuildCreateHash()
	res, err := creator.Execute(&usecase.CreateHashsRequest{
		Values: h.Values,
	})
	if err != nil {
		return Response500(err)
	}

	// 201レスポンス
	return Response201(res.Hashs)
}

// PutHash Hashの更新
func PutHash(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// バリデーション処理
	errs := ValidateBody(req.Body, ValidatorPutSetting)
	if errs != nil {
		return Response400(errs)
	}

	// JSON形式から構造体に変換
	var h RequestPutHash
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
	res, err := updator.Execute(&usecase.UpdateHashRequest{
		ID:    id,
		Value: h.Value,
	})
	if err != nil {
		if err.Error() == domain.ErrNotFound.Error() {
			return Response404()
		}
		return Response500(err)
	}

	return Response200(res.Hash)
}

// DeleteHash Hashの削除
func DeleteHash(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// パスパラメータからHashIDを取得する
	id, err := strconv.ParseUint(req.PathParameters["hash_id"], 10, 64)
	if err != nil {
		return Response500(err)
	}

	// 更新処理
	deleter := registry.GetFactory().BuildDeleteHash()
	_, err = deleter.Execute(&usecase.DeleteHashRequest{
		ID: id,
	})
	if err != nil {
		if err.Error() == domain.ErrNotFound.Error() {
			return Response404()
		}
		return Response500(err)
	}

	return Response200(nil)
}
