package registry

import "os"

var envs *Envs

// Envs 環境変数を表した構造体
type Envs struct {
	Cache map[string]string
}

// GetEnvs Envsのインスタンスを取得
func GetEnvs() *Envs {
	if envs == nil {
		envs = &Envs{
			Cache: map[string]string{},
		}
	}
	return envs
}

func (e *Envs) env(key string) string {
	return os.Getenv(key)
}

// DynamoLocalEndpoint DynamoDBのローカルEndpoint
func (e *Envs) DynamoLocalEndpoint() string {
	return e.env("DYNAMO_LOCAL_ENDPOINT")
}

// DynamoTableName DynamoDBのテーブル名
func (e *Envs) DynamoTableName() string {
	return e.env("DYNAMO_TABLE_NAME")
}

// DynamoPKName DynamoDBテーブルのPK列名
func (e *Envs) DynamoPKName() string {
	return e.env("DYNAMO_PK_NAME")
}

// DynamoSKName DynamoDBテーブルのSK列名
func (e *Envs) DynamoSKName() string {
	return e.env("DYNAMO_SK_NAME")
}

// RegionName AWSのRegion名
func (e *Envs) RegionName() string {
	return e.env("REGION_NAME")
}
