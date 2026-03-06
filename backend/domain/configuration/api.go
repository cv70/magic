package configuration

import (
	"backend/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func (d *ConfigurationDomain) ApiGetSystemConfig(c *gin.Context) {
	var req GetSystemConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	cfg, err := d.GetSystemConfig(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get system config: %v", err)
		utils.RespError(c, 500, "failed to get system config")
		return
	}

	utils.RespSuccess(c, GetSystemConfigRes{SystemConfig: cfg})
}

func (d *ConfigurationDomain) ApiSearchSystemConfigs(c *gin.Context) {
	var req SearchSystemConfigsReq
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

	cfgs, total, err := d.SearchSystemConfigs(c, req.Key, req.Category, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search system configs: %v", err)
		utils.RespError(c, 500, "failed to search system configs")
		return
	}

	utils.RespSuccess(c, SearchSystemConfigsRes{SystemConfigs: cfgs, Total: total})
}

func (d *ConfigurationDomain) ApiAddSystemConfig(c *gin.Context) {
	var req AddSystemConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	cfg, err := d.AddSystemConfig(c, req.Key, req.Value, req.Description, req.Category)
	if err != nil {
		klog.Errorf("failed to add system config: %v", err)
		utils.RespError(c, 500, "failed to add system config")
		return
	}

	utils.RespSuccess(c, AddSystemConfigRes{ID: cfg.ID})
}

func (d *ConfigurationDomain) ApiUpdateSystemConfig(c *gin.Context) {
	var req UpdateSystemConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	cfg, err := d.UpdateSystemConfig(c, req.ID, req.Key, req.Value, req.Description, req.Category)
	if err != nil {
		klog.Errorf("failed to update system config: %v", err)
		utils.RespError(c, 500, "failed to update system config")
		return
	}

	utils.RespSuccess(c, UpdateSystemConfigRes{ID: cfg.ID})
}

func (d *ConfigurationDomain) ApiGetProviderConfig(c *gin.Context) {
	var req GetProviderConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	cfg, err := d.GetProviderConfig(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get provider config: %v", err)
		utils.RespError(c, 500, "failed to get provider config")
		return
	}

	utils.RespSuccess(c, GetProviderConfigRes{ProviderConfig: cfg})
}

func (d *ConfigurationDomain) ApiSearchProviderConfigs(c *gin.Context) {
	var req SearchProviderConfigsReq
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

	cfgs, total, err := d.SearchProviderConfigs(c, req.ProviderName, req.ConfigKey, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search provider configs: %v", err)
		utils.RespError(c, 500, "failed to search provider configs")
		return
	}

	utils.RespSuccess(c, SearchProviderConfigsRes{ProviderConfigs: cfgs, Total: total})
}

func (d *ConfigurationDomain) ApiAddProviderConfig(c *gin.Context) {
	var req AddProviderConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	cfg, err := d.AddProviderConfig(c, req.ProviderName, req.ConfigKey, req.ConfigValue, req.Description)
	if err != nil {
		klog.Errorf("failed to add provider config: %v", err)
		utils.RespError(c, 500, "failed to add provider config")
		return
	}

	utils.RespSuccess(c, AddProviderConfigRes{ID: cfg.ID})
}

func (d *ConfigurationDomain) ApiUpdateProviderConfig(c *gin.Context) {
	var req UpdateProviderConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	cfg, err := d.UpdateProviderConfig(c, req.ID, req.ProviderName, req.ConfigKey, req.ConfigValue, req.Description)
	if err != nil {
		klog.Errorf("failed to update provider config: %v", err)
		utils.RespError(c, 500, "failed to update provider config")
		return
	}

	utils.RespSuccess(c, UpdateProviderConfigRes{ID: cfg.ID})
}

