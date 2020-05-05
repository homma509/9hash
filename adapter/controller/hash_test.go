package controller

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/mocks"
)

// TestPostHashs_201 複数Hashの新規作成の正常時
func TestPostHashs_201(t *testing.T) {
	// テスト用のDynamoDBを設定
	table := mocks.SetupDB(t)
	defer table.Cleanup()

	// リクエストパラーメーター設定
	body := map[string]interface{}{
		"values": []string{
			"http://test1.example.com",
			"http://test2.example.com",
		},
	}
	req, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	// 新規作成処理
	res := PostHashs(events.APIGatewayProxyRequest{
		Body: string(req),
	})

	// ステータスコードの確認
	if res.StatusCode != 201 {
		t.Errorf("PostHashs of StatusCode is wrong(expected=%d, actual=%d)", 201, res.StatusCode)
	}

	// 作成データを構造体に変換
	var hs []*domain.HashModel
	err = json.Unmarshal([]byte(res.Body), &hs)
	if err != nil {
		t.Fatal(err)
	}

	// データの確認
	for _, h := range hs {
		hash, err := table.HashOperator.GetHashByID(h.ID)
		if err != nil {
			t.Errorf("PostHash of Data is not found(ID = %d)", h.ID)
		}
		if h.Key != hash.Key {
			t.Errorf("PostHash of Key is wrong(expected=%s, actual=%s)", h.Key, hash.Key)
		}
		if h.Value != hash.Value {
			t.Errorf("PostHash of Value is wrong(expected=%s, actual=%s)", h.Value, hash.Value)
		}
	}
}
