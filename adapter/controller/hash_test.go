package controller

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/mocks"
)

// TestPostHashs_200 PostHashsAPI OKテストケース
func TestPostHashs_200(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
		req    map[string]interface{}
		want   []map[string]interface{}
	}{
		{
			"正常ケース",
			PostHashs,
			201,
			map[string]interface{}{
				"values": []string{
					"http://test1.example.com",
					"http://test2.example.com",
					"http://test3.example.com",
				},
			},
			[]map[string]interface{}{
				{"value": "http://test1.example.com"},
				{"value": "http://test2.example.com"},
				{"value": "http://test3.example.com"},
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// リクエストMapをJSONに変換
			req, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatal(err)
			}

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				Body: string(req),
			})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// レスポンスBodyをモデルへ変換
			var hashs []*domain.HashModel
			err = json.Unmarshal([]byte(res.Body), &hashs)
			if err != nil {
				t.Fatal(err)
			}

			// レスポンスデータの確認
			for _, hash := range hashs {
				// DynamoDBに保存されたデータを取得
				h, err := table.HashOperator.GetHashByID(hash.ID)
				if err != nil {
					t.Errorf("Data is not found(ID = %d)", hash.ID)
				}
				// Valueのチェック
				if hash.Value != h.Value {
					t.Errorf("Value is wrong(want=%s, actual=%s)", hash.Value, h.Value)
				}
			}

		})
	}
}

// TestPostHashs_400 PostHashs BadRequestテストケース
func TestPostHashs_400(t *testing.T) {
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
			"異常ケース: values 未入力",
			PostHashs,
			400,
			map[string]interface{}{
				"values": []string{},
			},
			map[string]interface{}{
				"message": "入力値を確認してください。",
				"errors": map[string]interface{}{
					"values": "値を入力してください。",
				},
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// リクエストMapをJSONに変換
			req, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatal(err)
			}

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				Body: string(req),
			})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// レスポンスBodyをMapへ変換
			var resBody map[string]interface{}
			err = json.Unmarshal([]byte(res.Body), &resBody)
			if err != nil {
				t.Fatal(err)
			}

			// レスポンスメッセージの確認
			if tc.want["message"] != resBody["message"] {
				t.Errorf("Response Body message is not equal(want=%s, actual=%v)", tc.want["message"], resBody["message"])
			}

			// レスポンスエラーの確認
			if !reflect.DeepEqual(tc.want["errors"], resBody["errors"]) {
				t.Errorf("Response Body errors is not equal(want=%v, actual=%v)", tc.want["errors"], resBody["errors"])
			}
		})
	}
}
