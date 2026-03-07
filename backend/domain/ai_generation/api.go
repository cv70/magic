package ai_generation

import (
	"backend/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func (d *AIGenerationDomain) ApiGetGenerator(c *gin.Context) {
	var req GetGeneratorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	gen, err := d.GetGenerator(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get generator: %v", err)
		utils.RespError(c, 500, "failed to get generator")
		return
	}

	utils.RespSuccess(c, GetGeneratorRes{Generator: gen})
}

func (d *AIGenerationDomain) ApiSearchGenerators(c *gin.Context) {
	var req SearchGeneratorsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	gens, total, err := d.SearchGenerators(c, req.Provider, req.Model, req.Enabled, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search generators: %v", err)
		utils.RespError(c, 500, "failed to search generators")
		return
	}

	utils.RespSuccess(c, SearchGeneratorsRes{Generators: gens, Total: total})
}

func (d *AIGenerationDomain) ApiAddGenerator(c *gin.Context) {
	var req AddGeneratorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	gen, err := d.AddGenerator(c, req.Name, req.Provider, req.Model, req.APIKey, req.APIEndpoint, enabled)
	if err != nil {
		klog.Errorf("failed to add generator: %v", err)
		utils.RespError(c, 500, "failed to add generator")
		return
	}

	utils.RespSuccess(c, AddGeneratorRes{ID: gen.ID})
}

func (d *AIGenerationDomain) ApiUpdateGenerator(c *gin.Context) {
	var req UpdateGeneratorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	gen, err := d.UpdateGenerator(c, req.ID, req.Name, req.Provider, req.Model, req.APIKey, req.APIEndpoint, req.Enabled)
	if err != nil {
		klog.Errorf("failed to update generator: %v", err)
		utils.RespError(c, 500, "failed to update generator")
		return
	}

	utils.RespSuccess(c, UpdateGeneratorRes{ID: gen.ID})
}

func (d *AIGenerationDomain) ApiGenerateContent(c *gin.Context) {
	var req GenerateContentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	id, err := d.GenerateContent(c, req.GeneratorID, req.Input)
	if err != nil {
		klog.Errorf("failed to generate content: %v", err)
		utils.RespError(c, 500, "failed to generate content")
		return
	}

	utils.RespSuccess(c, GenerateContentRes{ID: id})
}

func (d *AIGenerationDomain) ApiGetPromptTemplate(c *gin.Context) {
	var req GetPromptTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	pt, err := d.GetPromptTemplate(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get prompt template: %v", err)
		utils.RespError(c, 500, "failed to get prompt template")
		return
	}

	utils.RespSuccess(c, GetPromptTemplateRes{PromptTemplate: pt})
}

func (d *AIGenerationDomain) ApiSearchPromptTemplates(c *gin.Context) {
	var req SearchPromptTemplatesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	pts, total, err := d.SearchPromptTemplates(c, req.Name, req.Provider, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search prompt templates: %v", err)
		utils.RespError(c, 500, "failed to search prompt templates")
		return
	}

	utils.RespSuccess(c, SearchPromptTemplatesRes{PromptTemplates: pts, Total: total})
}

func (d *AIGenerationDomain) ApiAddPromptTemplate(c *gin.Context) {
	var req AddPromptTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	pt, err := d.AddPromptTemplate(c, req.Name, req.Description, req.Template, req.InputVariables)
	if err != nil {
		klog.Errorf("failed to add prompt template: %v", err)
		utils.RespError(c, 500, "failed to add prompt template")
		return
	}

	utils.RespSuccess(c, AddPromptTemplateRes{ID: pt.ID})
}

func (d *AIGenerationDomain) ApiUpdatePromptTemplate(c *gin.Context) {
	var req UpdatePromptTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	pt, err := d.UpdatePromptTemplate(c, req.ID, req.Name, req.Description, req.Template, req.InputVariables)
	if err != nil {
		klog.Errorf("failed to update prompt template: %v", err)
		utils.RespError(c, 500, "failed to update prompt template")
		return
	}

	utils.RespSuccess(c, UpdatePromptTemplateRes{ID: pt.ID})
}

// ===================== Generate Task API Handlers =====================

func (d *AIGenerationDomain) ApiGetGenerateTask(c *gin.Context) {
	var req GetGenerateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	task, err := d.GetGenerateTask(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get generate task: %v", err)
		utils.RespError(c, 500, "failed to get generate task")
		return
	}

	utils.RespSuccess(c, GetGenerateTaskRes{GenerateTask: task})
}

