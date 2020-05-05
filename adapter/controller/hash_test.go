package controller

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/mocks"
)

// TestGetHash GetHashAPI OKテストケース
func TestGetHash(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
		req    *domain.HashModel
		want   map[string]interface{}
	}{
		{
			"正常ケース: 200",
			GetHash,
			200,
			&domain.HashModel{
				Value: "http://test.example.com",
			},
			map[string]interface{}{
				"value": "http://test.example.com",
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// モックデータを作成
			hash, err := table.HashOperator.CreateHash(tc.req)

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"hash_id": fmt.Sprintf("%d", hash.ID),
				},
			})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// レスポンスBodyをモデルへ変換
			var h *domain.HashModel
			err = json.Unmarshal([]byte(res.Body), &h)
			if err != nil {
				t.Fatal(err)
			}

			// レスポンスデータの確認
			if hash.ID != h.ID {
				t.Errorf("ID is wrong(want=%d, actual=%d)", hash.ID, h.ID)
			}
			if hash.Key != h.Key {
				t.Errorf("Key is wrong(want=%s, actual=%s)", hash.Key, h.Key)
			}
			if hash.Value != h.Value {
				t.Errorf("Value is wrong(want=%s, actual=%s)", hash.Value, h.Value)
			}
		})
	}
}

// TestGetHashErr GetHashAPI BadRequestテストケース
func TestGetHashErr(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
		req    string
		want   map[string]interface{}
	}{
		{
			"異常ケース: 404(存在しないIDの取得)",
			GetHash,
			404,
			"1",
			map[string]interface{}{
				"message": "結果が見つかりません。",
			},
		},
		{
			"異常ケース: 500(IDの未入力)",
			GetHash,
			500,
			"",
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
					"hash_id": tc.req,
				},
			})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
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

// TestGetHashs GetHashsAPI OKテストケース
func TestGetHashs(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
		req    []*domain.HashModel
	}{
		{
			"正常ケース: 200",
			GetHashs,
			200,
			[]*domain.HashModel{
				{
					Value: "http://test1.example.com",
				},
				{
					Value: "http://test2.example.com",
				},
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// モックデータを作成
			hashs, err := table.HashOperator.CreateHashs(tc.req)

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// レスポンスBodyをモデルへ変換
			var hs []*domain.HashModel
			err = json.Unmarshal([]byte(res.Body), &hs)
			if err != nil {
				t.Fatal(err)
			}

			// モデルをソート
			sort.Slice(hashs, func(i, j int) bool { return hashs[i].ID < hashs[j].ID })
			sort.Slice(hs, func(i, j int) bool { return hs[i].ID < hs[j].ID })

			// レスポンスデータの確認
			for i, hash := range hashs {
				if hash.ID != hs[i].ID {
					t.Errorf("ID is wrong(want=%d, actual=%d)", hash.ID, hs[i].ID)
				}
				if hash.Key != hs[i].Key {
					t.Errorf("Key is wrong(want=%s, actual=%s)", hash.Key, hs[i].Key)
				}
				if hash.Value != hs[i].Value {
					t.Errorf("Value is wrong(want=%s, actual=%s)", hash.Value, hs[i].Value)
				}
			}
		})
	}
}

// TestPostHashs PostHashsAPI OKテストケース
func TestPostHashs(t *testing.T) {
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
			"正常ケース: 201",
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

// TestPostHashsErr PostHashs BadRequestテストケース
func TestPostHashsErr(t *testing.T) {
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
			"異常ケース: 400(valuesの未入力)",
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

// TestPutHash PutHashAPI OKテストケース
func TestPutHash(t *testing.T) {
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
			"正常ケース: 200",
			PutHash,
			200,
			map[string]interface{}{
				"value": "http://test_update.example.com",
			},
			map[string]interface{}{
				"value": "http://test_update.example.com",
			},
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// モックデータを作成
			hash, err := table.HashOperator.CreateHash(&domain.HashModel{
				Value: "test.example.com",
			})

			// リクエストMapをJSONに変換
			req, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatal(err)
			}

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"hash_id": fmt.Sprintf("%d", hash.ID),
				},
				Body: string(req),
			})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// DynamoDBのデータを取得
			h, err := table.HashOperator.GetHashByID(hash.ID)
			if err != nil {
				t.Fatal(err)
			}

			// データの確認
			if hash.ID != h.ID {
				t.Errorf("ID is wrong(want=%d, actual=%d)", hash.ID, h.ID)
			}
			if hash.Key != h.Key {
				t.Errorf("Key is wrong(want=%s, actual=%s)", hash.Key, h.Key)
			}
			if tc.want["value"] != h.Value {
				t.Errorf("Value is wrong(want=%s, actual=%s)", tc.want["value"], h.Value)
			}
		})
	}
}

// TestPutHashErr PutHash BadRequestテストケース
func TestPutHashErr(t *testing.T) {
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
			"異常ケース: 400(valueの未入力)",
			PutHash,
			400,
			map[string]interface{}{
				"hash_id": "1",
				"value":   "",
			},
			map[string]interface{}{
				"message": "入力値を確認してください。",
				"errors": map[string]interface{}{
					"value": "値を入力してください。",
				},
			},
		},
		{
			"異常ケース: 404(存在しないIDを更新)",
			PutHash,
			404,
			map[string]interface{}{
				"hash_id": "999",
				"value":   "test_update.example.com",
			},
			map[string]interface{}{
				"message": "結果が見つかりません。",
			},
		},
		{
			"異常ケース: 500(不正なID)",
			PutHash,
			500,
			map[string]interface{}{
				"hash_id": "dummy",
				"value":   "test_update.example.com",
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
			// モックデータを作成
			_, err := table.HashOperator.CreateHash(&domain.HashModel{
				Value: "test.example.com",
			})

			// リクエストMapをJSONに変換
			req, err := json.Marshal(tc.req)
			if err != nil {
				t.Fatal(err)
			}

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"hash_id": fmt.Sprintf("%v", tc.req["hash_id"]),
				},
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
			if _, ok := tc.want["errors"]; ok {
				if !reflect.DeepEqual(tc.want["errors"], resBody["errors"]) {
					t.Errorf("Response Body errors is not equal(want=%v, actual=%v)", tc.want["errors"], resBody["errors"])
				}
			}
		})
	}
}

// TestDeleteHash DeleteHashAPI OKテストケース
func TestDeleteHash(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// テストケース
	tests := []struct {
		name   string
		api    func(events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		status int
	}{
		{
			"正常ケース: 200",
			DeleteHash,
			200,
		},
	}

	// テストケースの実行
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// モックデータを作成
			hash, err := table.HashOperator.CreateHash(&domain.HashModel{
				Value: "test.example.com",
			})
			if err != nil {
				t.Fatal(err)
			}

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"hash_id": fmt.Sprintf("%d", hash.ID),
				},
			})

			// レスポンスStatusCodeの確認
			if res.StatusCode != tc.status {
				t.Errorf("StatusCode is wrong(want=%d, actual=%d)", tc.status, res.StatusCode)
			}

			// DynamoDBのデータを取得
			h, _ := table.HashOperator.GetHashByID(hash.ID)
			if h != nil {
				t.Errorf("Hash Data not delted(ID=%d)", h.ID)
			}
		})
	}
}

// TestDeleteHashErr DeleteHash BadRequestテストケース
func TestDeleteHashErr(t *testing.T) {
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
			"異常ケース: 404(存在しないIDを削除)",
			DeleteHash,
			404,
			map[string]interface{}{
				"hash_id": "999",
			},
			map[string]interface{}{
				"message": "結果が見つかりません。",
			},
		},
		{
			"異常ケース: 500(不正なID)",
			DeleteHash,
			500,
			map[string]interface{}{
				"hash_id": "dummy",
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
			// モックデータを作成
			_, err := table.HashOperator.CreateHash(&domain.HashModel{
				Value: "test.example.com",
			})

			// APIの実行
			res := tc.api(events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"hash_id": fmt.Sprintf("%v", tc.req["hash_id"]),
				},
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
		})
	}
}
