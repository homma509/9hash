package registry

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/homma509/9hash/adapter/gateway"
	"github.com/homma509/9hash/adapter/gateway/db"
	"github.com/homma509/9hash/domain"
	"github.com/homma509/9hash/interactor"
	"github.com/homma509/9hash/usecase"
)

var factory *Factory

type Factory struct {
	Envs  *Envs
	cache map[string]interface{}
}

// ClearFactory インスタンスを削除します
func ClearFactory() {
	factory = nil
}

// GetFactory Factoryのインスタンスを取得します
func GetFactory() *Factory {
	if factory == nil {
		factory = &Factory{
			Envs: GetEnvs(),
		}
	}
	return factory
}

func (f *Factory) container(key string, builder func() interface{}) interface{} {
	if f.cache == nil {
		f.cache = map[string]interface{}{}
	}
	if _, ok := f.cache[key]; !ok {
		f.cache[key] = builder()
	}
	return f.cache[key]
}

func (f *Factory) BuildDynamoClient() *db.DynamoClient {
	return f.container("DynamoClient", func() interface{} {
		config := &aws.Config{
			Region: f.Envs.RegionName(),
		}
		if f.Envs.DynamoLocalEndpoint() != "" {
			config.Credentials = credentials.NewStaticCredentials("dummy_id", "dummy_secret", "dymmy_token")
			config.Endpoint = aws.String(f.Envs.DynamoLocalEndpoint())
		}
		return db.NewClient(config)
	}).(*db.DynamoClient)
}

func (f *Factory) BuildResourceTableOperator() *db.ResourceTableOperator {
	return f.container("ResourceTableOperator", func() interface{} {
		return db.NewResourceTableOperator(
			f.BuildDynamoClient(),
			f.Envs.DynamoTableName(),
		)
	}).(*db.ResourceTableOperator)
}

func (f *Factory) BuildDynamoModelMapper() *db.DynamoModelMapper {
	return f.container("DynamoModelMapper", func() interface{} {
		return &db.DynamoModelMapper{
			Client:    f.BuildDynamoClient(),
			TableName: f.Envs.DynamoTableName(),
			PKName:    f.Envs.DynamoPKName(),
			SKName:    f.Envs.DynamoSKName(),
		}
	}).(*db.DynamoModelMapper)
}

func (f *Factory) BuildHashOperator() domain.HashRepository {
	return f.container("HashOperator", func() interface{} {
		return &gateway.HashOperator{
			Client: f.BuildResourceTableOperator(),
			Mapper: f.BuildDynamoModelMapper(),
		}
	}).(domain.HashRepository)
}

func (f *Factory) BuildCreateHash() usecase.ICreateHash {
	return f.container("CreateHash", func() interface{} {
		return interactor.NewHashCreator(
			f.BuildHashOperator()
		)
	}).(usecase.ICreateHash)
}
