package mocks

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/homma509/9hash/adapter/gateway/db"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/registry"
)

// TableOperator テーブル操作を表した構造体
type TableOperator struct {
	Operator     *db.ResourceTableOperator
	HashOperator domain.HashRepository
}

// SetupDB テスト用のテーブルの作成
func SetupDB(t *testing.T) *TableOperator {
	t.Helper()

	os.Setenv("DYNAMO_TABLE_NAME", generateRandomTableName(t))

	registry.ClearFactory()
	f := registry.GetFactory()
	o := &TableOperator{}
	o.Operator = f.BuildResourceTableOperator()
	o.HashOperator = f.BuildHashOperator()

	o.Operator.CreateTestTable()

	return o
}

// Cleanup テスト用テーブルの削除
func (o *TableOperator) Cleanup() {
	o.Operator.DropTable()
}

func generateRandomTableName(t *testing.T) string {
	length := 60
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		t.Fatal(err)
	}
	l := length / 2
	if length%2 == 1 {
		l++
	}
	return fmt.Sprintf("%x", buf[:l])
}
