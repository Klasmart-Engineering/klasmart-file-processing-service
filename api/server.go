package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"sync"
)

type Server struct {
	engine *gin.Engine
}

func (s *Server) Start(){
	//gin.DefaultWriter = log.LoggerOut()
	//gin.DefaultErrorWriter = log.LoggerOut()
	s.engine = gin.Default()
	s.route()
	s.engine.Run(fmt.Sprintf(":%d", config.Get().API.Port))
}
var (
	_apiServer *Server
	_apiServerOnce sync.Once
)

func GetServer() *Server{
	_apiServerOnce.Do(func() {
		_apiServer = new(Server)
	})
	return _apiServer
}