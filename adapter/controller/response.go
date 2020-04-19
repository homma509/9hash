package controller

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/glog"
)

// Response400Body バリデーションエラーメッセージを含めた400レスポンス
type Response400Body struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func headers() map[string]string {
	return map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
}

// Response200 OKメッセージを含めたレスポンス
func Response200() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers(),
		Body:       `{"message":"OK"}`,
	}
}

// Response201 IDを含めたレスポンス
func Response201(id uint64) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    headers(),
		Body:       fmt.Sprintf(`{"message":"OK","id":%d}`, id),
	}
}

// Response400 400レスポンス
func Response400(errs map[string]error) events.APIGatewayProxyResponse {
	glog.Warningf("%+v", errs)
	res := &Response400Body{
		Message: "入力値を確認してください。",
		Errors:  ToMessages(errs),
	}

	b, err := json.Marshal(res)
	if err != nil {
		return Response500(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers:    headers(),
		Body:       string(b),
	}
}

// Response500 500レスポンス
func Response500(err error) events.APIGatewayProxyResponse {
	glog.Errorf("%+v", err)
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    headers(),
		Body:       `{"message":"サーバエラーが発生しました。"}`,
	}
}
