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

	// ===================== Authentication Routes (无需认证) =====================
	{
		identityDomain := identity.NewIdentityDomain(registry.DB)
		v1.POST("/auth/login", identityDomain.ApiLogin)
		v1.POST("/auth/register", identityDomain.ApiRegister)
	}

	// ===================== Protected Routes (需要认证) =====================
	protected := v1.Group("/")
	protected.Use(infra.AuthMiddleware())

	// 认证用户信息
	{
		identityDomain := identity.NewIdentityDomain(registry.DB)
		protected.POST("/auth/me", identityDomain.ApiMe)
	}

	// 内容和草稿模块
	{
		contentDomain := content.NewContentDomainWithDB(registry.DB)

		// 草稿 API
		protected.POST("/drafts/create", contentDomain.ApiCreateDraft)
		protected.POST("/drafts/get", contentDomain.ApiGetDraft)
		protected.POST("/drafts/search", contentDomain.ApiSearchDrafts)
		protected.POST("/drafts/update", contentDomain.ApiUpdateDraft)
		protected.POST("/drafts/delete", contentDomain.ApiDeleteDraft)
		protected.POST("/drafts/:id/versions", contentDomain.ApiGetDraftVersions)
		protected.POST("/drafts/:id/revert/:version", contentDomain.ApiRevertDraftVersion)
		protected.POST("/drafts/:id/publish", contentDomain.ApiPublishDraft)

		// 内容 API
		protected.POST("/content/create", contentDomain.ApiCreateContent)
		protected.POST("/content/get", contentDomain.ApiGetContent)
		protected.POST("/content/search", contentDomain.ApiSearchContents)
		protected.POST("/content/update", contentDomain.ApiUpdateContent)
		protected.POST("/content/delete", contentDomain.ApiDeleteContent)
	}

	// AI 生成模块
	{
		aiDomain := ai_generation.NewAIGenerationDomain(registry.DB)
		protected.POST("/ai/generator/get", aiDomain.ApiGetGenerator)
		protected.POST("/ai/generator/search", aiDomain.ApiSearchGenerators)
		protected.POST("/ai/generator/add", aiDomain.ApiAddGenerator)
		protected.POST("/ai/generator/update", aiDomain.ApiUpdateGenerator)
		protected.POST("/ai/generate", aiDomain.ApiGenerateContent)
		protected.POST("/ai/task/get", aiDomain.ApiGetGenerateTask)
		protected.POST("/ai/prompt-template/get", aiDomain.ApiGetPromptTemplate)
		protected.POST("/ai/prompt-template/search", aiDomain.ApiSearchPromptTemplates)
		protected.POST("/ai/prompt-template/add", aiDomain.ApiAddPromptTemplate)
		protected.POST("/ai/prompt-template/update", aiDomain.ApiUpdatePromptTemplate)
	}

	// 配置模块
	{
		configDomain := configuration.NewConfigurationDomain(registry.DB)
		protected.POST("/config/system/get", configDomain.ApiGetSystemConfig)
		protected.POST("/config/system/search", configDomain.ApiSearchSystemConfigs)
		protected.POST("/config/system/add", configDomain.ApiAddSystemConfig)
		protected.POST("/config/system/update", configDomain.ApiUpdateSystemConfig)
		protected.POST("/config/provider/get", configDomain.ApiGetProviderConfig)
		protected.POST("/config/provider/search", configDomain.ApiSearchProviderConfigs)
		protected.POST("/config/provider/add", configDomain.ApiAddProviderConfig)
		protected.POST("/config/provider/update", configDomain.ApiUpdateProviderConfig)
	}

	// 身份认证模块 (Admin)
	{
		identityDomain := identity.NewIdentityDomain(registry.DB)
		protected.POST("/identity/user/get", identityDomain.ApiGetUser)
		protected.POST("/identity/user/search", identityDomain.ApiSearchUsers)
		protected.POST("/identity/user/add", identityDomain.ApiAddUser)
		protected.POST("/identity/user/update", identityDomain.ApiUpdateUser)
		protected.POST("/identity/role/get", identityDomain.ApiGetRole)
		protected.POST("/identity/role/search", identityDomain.ApiSearchRoles)
		protected.POST("/identity/role/add", identityDomain.ApiAddRole)
		protected.POST("/identity/role/update", identityDomain.ApiUpdateRole)
		protected.POST("/identity/permission/get", identityDomain.ApiGetPermission)
		protected.POST("/identity/permission/search", identityDomain.ApiSearchPermissions)
		protected.POST("/identity/permission/add", identityDomain.ApiAddPermission)
		protected.POST("/identity/permission/update", identityDomain.ApiUpdatePermission)
	}

	// 发布模块
	{
		publishingDomain := publishing.NewPublishingDomain(registry.DB)
		protected.POST("/publish/publisher/get", publishingDomain.ApiGetPublisher)
		protected.POST("/publish/publisher/search", publishingDomain.ApiSearchPublishers)
		protected.POST("/publish/publisher/add", publishingDomain.ApiAddPublisher)
		protected.POST("/publish/publisher/update", publishingDomain.ApiUpdatePublisher)
		protected.POST("/publish/content", publishingDomain.ApiPublishContent)
		protected.POST("/publish/task/get", publishingDomain.ApiGetPublishTask)
		protected.POST("/publish/task/search", publishingDomain.ApiSearchPublishTasks)
		protected.POST("/publish/log/get", publishingDomain.ApiGetPublishLog)
		protected.POST("/publish/log/search", publishingDomain.ApiSearchPublishLogs)
		protected.POST("/analytics/summary", publishingDomain.ApiAnalyticsSummary)
		protected.POST("/analytics/ranking", publishingDomain.ApiAnalyticsRanking)
		protected.POST("/analytics/platform-comparison", publishingDomain.ApiAnalyticsPlatformComparison)
	}

	// 调度模块
	{
		schedulingDomain := scheduling.NewSchedulingDomain(registry.DB)
		protected.POST("/scheduling/task/get", schedulingDomain.ApiGetTask)
		protected.POST("/scheduling/task/search", schedulingDomain.ApiSearchTasks)
		protected.POST("/scheduling/task/add", schedulingDomain.ApiAddTask)
		protected.POST("/scheduling/task/update", schedulingDomain.ApiUpdateTask)
		protected.POST("/scheduling/task/run", schedulingDomain.ApiRunTask)
		protected.POST("/scheduling/task/stop", schedulingDomain.ApiStopTask)
		protected.POST("/scheduling/task/restart", schedulingDomain.ApiRestartTask)
	}

	err = r.Run(":8888")
	mistake.Unwrap(err)
}
