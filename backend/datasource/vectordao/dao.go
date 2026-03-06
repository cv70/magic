package vectordao

import (
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

type VectorDB struct {
	client.Client
}

func NewVectorDB(cli client.Client) *VectorDB {
	return &VectorDB{
		Client: cli,
	}
}
