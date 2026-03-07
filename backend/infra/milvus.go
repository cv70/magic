package infra

import (
	"backend/config"
	"backend/datasource/vectordao"
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func NewMilvus(ctx context.Context, c *config.MilvusConfig) (*vectordao.VectorDB, error) {
	cfg := client.Config{
		Address:  c.Host,
		Username: c.User,
		Password: c.Password,
		DBName:   c.DBName,
	}

	milvusClient, err := client.NewClient(ctx, cfg)
	return vectordao.NewVectorDB(milvusClient), err
}
