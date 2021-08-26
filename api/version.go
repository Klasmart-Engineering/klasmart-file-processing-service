package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/constant"
	"net/http"
)

func (s *Server) version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"git_hash":        constant.GitHash,
		"build_timestamp": constant.BuildTimestamp,
	})
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
	})
}
