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

// Factory インターフェースの詳細を生成する構造体
type Factory struct {
	Envs  *Envs
	cache map[string]interface{}
}

// ClearFactory Factoryのクリア
func ClearFactory() {
	factory = nil
}

// GetFactory Factoryのインスタンスを取得
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

// BuildDynamoClient DyamoClientの生成
func (f *Factory) BuildDynamoClient() *db.DynamoClient {
	return f.container("DynamoClient", func() interface{} {
		config := &aws.Config{
			Region: aws.String(f.Envs.RegionName()),
		}
		if f.Envs.DynamoDBLocalEndpoint() != "" {
			config.Credentials = credentials.NewStaticCredentials("dummy_id", "dummy_secret", "dymmy_token")
			config.Endpoint = aws.String(f.Envs.DynamoDBLocalEndpoint())
		}
		return db.NewClient(config)
	}).(*db.DynamoClient)
}

// BuildResourceTableOperator ResourceTableOperatorの生成
func (f *Factory) BuildResourceTableOperator() *db.ResourceTableOperator {
	return f.container("ResourceTableOperator", func() interface{} {
		return db.NewResourceTableOperator(
			f.BuildDynamoClient(),
			f.Envs.DynamoTableName(),
		)
	}).(*db.ResourceTableOperator)
}

// BuildDynamoModelMapper DynamoModelMapperの生成
func (f *Factory) BuildDynamoModelMapper() *db.DynamoModelMapper {
	return f.container("DynamoModelMapper", func() interface{} {
		return &db.DynamoModelMapper{
			Client:    f.BuildResourceTableOperator(),
			TableName: f.Envs.DynamoTableName(),
			PKName:    f.Envs.DynamoPKName(),
			SKName:    f.Envs.DynamoSKName(),
		}
	}).(*db.DynamoModelMapper)
}

// BuildHashOperator HashOperatorの生成
func (f *Factory) BuildHashOperator() domain.HashRepository {
	return f.container("HashOperator", func() interface{} {
		return &gateway.HashOperator{
			Client: f.BuildResourceTableOperator(),
			Mapper: f.BuildDynamoModelMapper(),
		}
	}).(domain.HashRepository)
}

// BuildGetHash GetHashの生成
func (f *Factory) BuildGetHash() usecase.IGetHash {
	return f.container("GetHash", func() interface{} {
		return interactor.NewHashGetter(
			f.BuildHashOperator())
	}).(usecase.IGetHash)
}

// BuildGetHashs GetHashsの生成
func (f *Factory) BuildGetHashs() usecase.IGetHashs {
	return f.container("GetHashs", func() interface{} {
		return interactor.NewHashsGetter(
			f.BuildHashOperator())
	}).(usecase.IGetHashs)
}

// BuildCreateHash CreateHashの生成
func (f *Factory) BuildCreateHash() usecase.ICreateHashs {
	return f.container("CreateHash", func() interface{} {
		return interactor.NewHashsCreator(
			f.BuildHashOperator())
	}).(usecase.ICreateHashs)
}

// BuildUpdateHash UpdateHashの生成
func (f *Factory) BuildUpdateHash() usecase.IUpdateHash {
	return f.container("UpdateHash", func() interface{} {
		return interactor.NewHashUpdator(
			f.BuildHashOperator())
	}).(usecase.IUpdateHash)
}

// BuildDeleteHash DeleteHashの生成
func (f *Factory) BuildDeleteHash() usecase.IDeleteHash {
	return f.container("DeleteHash", func() interface{} {
		return interactor.NewHashDeleter(
			f.BuildHashOperator())
	}).(usecase.IDeleteHash)
}

// BuildGetURL GetURLの生成
func (f *Factory) BuildGetURL() usecase.IGetURL {
	return f.container("GetURL", func() interface{} {
		return interactor.NewURLGetter(
			f.BuildHashOperator())
	}).(usecase.IGetURL)
}
