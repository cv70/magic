package scheduling

import (
	"backend/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func (d *SchedulingDomain) ApiGetTask(c *gin.Context) {
	var req GetTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	task, err := d.GetTask(c, req.ID)
	if err != nil {
		klog.Errorf("failed to get task: %v", err)
		utils.RespError(c, 500, "failed to get task")
		return
	}

	utils.RespSuccess(c, GetTaskRes{Task: task})
}

func (d *SchedulingDomain) ApiSearchTasks(c *gin.Context) {
	var req SearchTasksReq
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

	tasks, total, err := d.SearchTasks(c, req.Name, req.TaskType, req.SchedulerID, req.Enabled, req.Page, req.Limit)
	if err != nil {
		klog.Errorf("failed to search tasks: %v", err)
		utils.RespError(c, 500, "failed to search tasks")
		return
	}

	utils.RespSuccess(c, SearchTasksRes{Tasks: tasks, Total: total})
}

func (d *SchedulingDomain) ApiAddTask(c *gin.Context) {
	var req AddTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	task, err := d.AddTask(c, req.Name, req.TaskType, req.SchedulerID, req.CronExpression, enabled)
	if err != nil {
		klog.Errorf("failed to add task: %v", err)
		utils.RespError(c, 500, "failed to add task")
		return
	}

	utils.RespSuccess(c, AddTaskRes{ID: task.ID})
}

func (d *SchedulingDomain) ApiUpdateTask(c *gin.Context) {
	var req UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	_, err := d.UpdateTask(c, req.ID, req.Name, req.TaskType, req.CronExpression, req.SchedulerID, req.Enabled)
	if err != nil {
		klog.Errorf("failed to update task: %v", err)
		utils.RespError(c, 500, "failed to update task")
		return
	}

	utils.RespSuccess(c, UpdateTaskRes{ID: req.ID})
}

func (d *SchedulingDomain) ApiRunTask(c *gin.Context) {
	var req RunTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	id, err := d.RunTask(c, req.ID)
	if err != nil {
		klog.Errorf("failed to run task: %v", err)
		utils.RespError(c, 500, "failed to run task")
		return
	}

	utils.RespSuccess(c, TaskActionRes{ID: id})
}

func (d *SchedulingDomain) ApiStopTask(c *gin.Context) {
	var req StopTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	id, err := d.StopTask(c, req.ID)
	if err != nil {
		klog.Errorf("failed to stop task: %v", err)
		utils.RespError(c, 500, "failed to stop task")
		return
	}

	utils.RespSuccess(c, TaskActionRes{ID: id})
}

func (d *SchedulingDomain) ApiRestartTask(c *gin.Context) {
	var req RestartTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		klog.Errorf("failed to parse body: %v", err)
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	id, err := d.RestartTask(c, req.ID)
	if err != nil {
		klog.Errorf("failed to restart task: %v", err)
		utils.RespError(c, 500, "failed to restart task")
		return
	}

	utils.RespSuccess(c, TaskActionRes{ID: id})
}
