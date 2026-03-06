package main

import (
	"backend/config"
	"backend/domain/financing"
	"backend/domain/news"
	"backend/infra"
	"context"

	"github.com/cv70/pkgo/mistake"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	mistake.Unwrap(err)

	// Initialize infrastructure with configuration
	registry, err := infra.NewRegistry(ctx, cfg)
	mistake.Unwrap(err)

	r := gin.Default()
	v1 := r.Group("/api/v1")

	// 融资模块
	{
		financingDomain := financing.NewFinancingDomain(registry.DB)
		v1.POST("/financing/business-plans/get", financingDomain.ApiGetBusinessPlan)
		v1.POST("/financing/business-plans/batch-get", financingDomain.ApiGetBusinessPlans)
		v1.POST("/financing/business-plans/create", financingDomain.ApiCreateBusinessPlan)
		v1.POST("/financing/business-plans/update", financingDomain.ApiUpdateBusinessPlan)
		v1.POST("/financing/business-plans/delete", financingDomain.ApiDeleteBusinessPlan)
	}

	// 新闻中心模块
	{
		newsDomain := news.NewsDomain{
			DB:       registry.DB,
			VectorDB: registry.VectorDB,
		}
		v1.POST("/news/list", newsDomain.ApiGetNews)
		v1.POST("/news/search", newsDomain.ApiSearchNews)
		v1.POST("/news/add", newsDomain.ApiAddNews)
	}

	err = r.Run(":8888")
	mistake.Unwrap(err)
}
