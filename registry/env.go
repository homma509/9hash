package registry

import "os"

var envs *Envs

type Envs struct {
	Cache map[string]string
}

// GetEnvs Envsのインスタンスを取得します
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

func (e *Envs) DynamoLocalEndpoint() string {
	return e.env("DYNAMO_LOCAL_ENDPOINT")
}

func (e *Envs) DynamoTableName() string {
	return f.env("DYNAMO_TABLE_NAME")
}

func (e *Envs) DynamoPKName() string {
	return e.env("DYNAMO_PK_NAME")
}

func (e *Envs) DynamoSKName() string {
	return e.env("DYNAMO_SK_NAME")
}

func (e *Envs) RegionName() string {
	return e.env("REGION_NAME")
}
