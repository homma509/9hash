package controller

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
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
	if res.StatusCode != 201 {
		t.Errorf("PostHashs of StatusCode is wrong(expected=%d, actual=%d)", 201, res.StatusCode)
	}
}
