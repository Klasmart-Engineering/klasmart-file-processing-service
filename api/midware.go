package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/crypto"
	"net/http"
	"strconv"
	"strings"
	"time"
)
const (
	tokenValidDuration = time.Minute * 20
)
func (s *Server) checkSecret(c *gin.Context) {
	secret := c.Query("s")
	if secret == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "s is empty",
		})
		c.Abort()
		return
	}

	secretParam := "s=" + secret
	paramsParts := strings.Split(c.Request.RequestURI, secretParam)
	if len(paramsParts) != 2 || paramsParts[1] != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "secret must be the last param",
		})
		c.Abort()
		return
	}
	url := paramsParts[0]
	if strings.HasSuffix(url, "&") {
		url = url[:len(url) - 1]
	}

	token := crypto.HMac(url)
	if token != secret {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "secret invalid",
		})
		c.Abort()
		return
	}

}
func (s *Server) checkT(c *gin.Context) bool {
	ts := c.Query("t")
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid param t",
		})
		c.Abort()
		return false
	}
	t := time.Unix(tsInt, 0)
	now := time.Now()
	dur := now.Sub(t)
	if dur < 0 ||
		dur > tokenValidDuration {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "t is out of range",
		})
		c.Abort()
		return false
	}
	return true
}

func (s *Server) mustToken(c *gin.Context){
	//If no secret key is configured, ignore check token
	if config.Get().API.SecretKey == "" {
		return
	}

	ctn := s.checkT(c)
	if !ctn{
		return
	}
	s.checkSecret(c)
}