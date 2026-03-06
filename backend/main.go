package main

import (
	"backend/config"
	"backend/domain/ai_generation"
	"backend/domain/configuration"
	"backend/domain/content"
	"backend/domain/identity"
	"backend/domain/publishing"
	"backend/domain/scheduling"
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

	// 内容模块
	{
		contentDomain := content.NewContentDomain(registry.DB)
		v1.POST("/content/get", contentDomain.ApiGetContent)
		v1.POST("/content/search", contentDomain.ApiSearchContents)
		v1.POST("/content/add", contentDomain.ApiAddContent)
		v1.POST("/content/update", contentDomain.ApiUpdateContent)
		v1.POST("/content/delete", contentDomain.ApiDeleteContent)
	}

	// AI 生成模块
	{
		aiDomain := ai_generation.NewAIGenerationDomain(registry.DB)
		v1.POST("/ai/generator/get", aiDomain.ApiGetGenerator)
		v1.POST("/ai/generator/search", aiDomain.ApiSearchGenerators)
		v1.POST("/ai/generator/add", aiDomain.ApiAddGenerator)
		v1.POST("/ai/generator/update", aiDomain.ApiUpdateGenerator)
		v1.POST("/ai/generate", aiDomain.ApiGenerateContent)
		v1.POST("/ai/prompt-template/get", aiDomain.ApiGetPromptTemplate)
		v1.POST("/ai/prompt-template/search", aiDomain.ApiSearchPromptTemplates)
		v1.POST("/ai/prompt-template/add", aiDomain.ApiAddPromptTemplate)
		v1.POST("/ai/prompt-template/update", aiDomain.ApiUpdatePromptTemplate)
	}

	// 配置模块
	{
		configDomain := configuration.NewConfigurationDomain(registry.DB)
		v1.POST("/config/system/get", configDomain.ApiGetSystemConfig)
		v1.POST("/config/system/search", configDomain.ApiSearchSystemConfigs)
		v1.POST("/config/system/add", configDomain.ApiAddSystemConfig)
		v1.POST("/config/system/update", configDomain.ApiUpdateSystemConfig)
		v1.POST("/config/provider/get", configDomain.ApiGetProviderConfig)
		v1.POST("/config/provider/search", configDomain.ApiSearchProviderConfigs)
		v1.POST("/config/provider/add", configDomain.ApiAddProviderConfig)
		v1.POST("/config/provider/update", configDomain.ApiUpdateProviderConfig)
	}

	// 身份认证模块
	{
		identityDomain := identity.NewIdentityDomain(registry.DB)
		v1.POST("/identity/user/get", identityDomain.ApiGetUser)
		v1.POST("/identity/user/search", identityDomain.ApiSearchUsers)
		v1.POST("/identity/user/add", identityDomain.ApiAddUser)
		v1.POST("/identity/user/update", identityDomain.ApiUpdateUser)
		v1.POST("/identity/role/get", identityDomain.ApiGetRole)
		v1.POST("/identity/role/search", identityDomain.ApiSearchRoles)
		v1.POST("/identity/role/add", identityDomain.ApiAddRole)
		v1.POST("/identity/role/update", identityDomain.ApiUpdateRole)
		v1.POST("/identity/permission/get", identityDomain.ApiGetPermission)
		v1.POST("/identity/permission/search", identityDomain.ApiSearchPermissions)
		v1.POST("/identity/permission/add", identityDomain.ApiAddPermission)
		v1.POST("/identity/permission/update", identityDomain.ApiUpdatePermission)
	}

	// 发布模块
	{
		publishingDomain := publishing.NewPublishingDomain(registry.DB)
		v1.POST("/publish/publisher/get", publishingDomain.ApiGetPublisher)
		v1.POST("/publish/publisher/search", publishingDomain.ApiSearchPublishers)
		v1.POST("/publish/publisher/add", publishingDomain.ApiAddPublisher)
		v1.POST("/publish/publisher/update", publishingDomain.ApiUpdatePublisher)
		v1.POST("/publish/content", publishingDomain.ApiPublishContent)
		v1.POST("/publish/task/get", publishingDomain.ApiGetPublishTask)
		v1.POST("/publish/task/search", publishingDomain.ApiSearchPublishTasks)
		v1.POST("/publish/log/get", publishingDomain.ApiGetPublishLog)
		v1.POST("/publish/log/search", publishingDomain.ApiSearchPublishLogs)
	}

	// 调度模块
	{
		schedulingDomain := scheduling.NewSchedulingDomain(registry.DB)
		v1.POST("/scheduling/task/get", schedulingDomain.ApiGetTask)
		v1.POST("/scheduling/task/search", schedulingDomain.ApiSearchTasks)
		v1.POST("/scheduling/task/add", schedulingDomain.ApiAddTask)
		v1.POST("/scheduling/task/update", schedulingDomain.ApiUpdateTask)
		v1.POST("/scheduling/task/run", schedulingDomain.ApiRunTask)
		v1.POST("/scheduling/task/stop", schedulingDomain.ApiStopTask)
		v1.POST("/scheduling/task/restart", schedulingDomain.ApiRestartTask)
	}

	err = r.Run(":8888")
	mistake.Unwrap(err)
}
