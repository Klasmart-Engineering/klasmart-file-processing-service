package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/runtime"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/service"
	"io/ioutil"
	"net/http"
)

//file processor api
func (s *Server) processFile(c *gin.Context) {
	req := new(entity.PublishFileRequest)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	ctx := c.Request.Context()

	err = service.GetFileProcessingService().MQ().Publish(ctx, req.Topic(), req.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, "")
}

func (s *Server) workers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"workers": runtime.GetWorkersInfo().Num(),
	})
}

func (s *Server) pendingList(c *gin.Context) {
	messages, err := service.GetFileProcessingService().PendingMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (s *Server) failedList(c *gin.Context) {
	//contains handle failed list & publish failed list
	publishFailed, _ := ioutil.ReadFile(config.Get().MQ.RedisFailedPersistence)
	handleFailed, _ := ioutil.ReadFile(config.Get().Log.FailedFile)
	c.JSON(http.StatusOK, gin.H{
		"publish_failed": string(publishFailed),
		"handle_failed": string(handleFailed),
	})
}


func (s *Server) supportClassifies(c *gin.Context) {
	extensions := service.GetFileProcessingService().SupportExtensions()
	c.JSON(http.StatusOK, gin.H{
		"supports": extensions,
	})
}