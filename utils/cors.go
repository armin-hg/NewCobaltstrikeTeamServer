package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc { //跨域问题
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		//c.Header("Content-Security-Policy", "default-src 'none'; connect-src 'self' ws://127.0.0.1:88;")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		if c.Request.TLS != nil {
			//fmt.Println(c.Request.Host)
			if strings.Contains(c.Request.Host, ":8443") {
				//fmt.Println("https:8443")
				secureMiddleware := secure.New(secure.Options{
					SSLRedirect: false,
					//SSLHost:     config.Config.TeamServer.Domain + ":" + "8443", //客户端请求ui
				})
				err := secureMiddleware.Process(c.Writer, c.Request)
				if err != nil {
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
