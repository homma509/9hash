package controller

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/mocks"
)

// TestGetURL GetURLAPI OKテストケース
func TestGetURL(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
		req    map[string]interface{}
		want   map[string]interface{}
	}{
		{
			"正常ケース: 308",
			GetURL,
			308,
			map[string]interface{}{
				"key":   "dummy_key",
				"value": "http://test.example.com",
			},
			map[string]interface{}{
				"location":      "http://test.example.com",
				"cache-control": "no-store",
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// モックデータを作成
			hash, err := table.HashOperator.CreateHash(&domain.HashModel{
				Key:   fmt.Sprintf("%v", tc.req["key"]),
				Value: fmt.Sprintf("%v", tc.req["value"]),
			})
			if err != nil {
				t.Fatal(err)
			}

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"key": hash.Key,
				},
			})

			// レスポンスStatusCodeの確認
			if tc.status != res.StatusCode {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// レスポンスHeaderの確認
			if reflect.DeepEqual(tc.want, res.Headers) {
				t.Errorf("Header is wrong(want=%v, actual=%v)", tc.want, res.Headers)
			}
		})
	}
}

// TestGetURLErr GetURLAPI 異常テストケース
func TestGetURLErr(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
		req    map[string]interface{}
		want   map[string]interface{}
	}{
		{
			"異常ケース: 404(存在しないKey)",
			GetURL,
			404,
			map[string]interface{}{
				"key": "dummy_key",
			},
			map[string]interface{}{
				"message": "結果が見つかりません。",
			},
		},
		{
			"異常ケース: 500(Keyの未入力)",
			GetURL,
			500,
			map[string]interface{}{
				"key": "",
			},
			map[string]interface{}{
				"message": "サーバエラーが発生しました。",
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"key": fmt.Sprintf("%v", tc.req["key"]),
				},
			})

			// レスポンスStatusCodeの確認
			if tc.status != res.StatusCode {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// レスポンスBodyをMapへ変換
			var resBody map[string]interface{}
			err := json.Unmarshal([]byte(res.Body), &resBody)
			if err != nil {
				t.Fatal(err)
			}

			// レスポンスメッセージの確認
			if tc.want["message"] != resBody["message"] {
				t.Errorf("Response Body message is not equal(want=%s, actual=%v)", tc.want["message"], resBody["message"])
			}
		})
	}
}
