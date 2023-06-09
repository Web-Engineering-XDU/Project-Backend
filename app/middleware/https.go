package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func Https(c *gin.Context) {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     ":8080",
	})
	err := secureMiddleware.Process(c.Writer, c.Request)
	// If there was an error, do not continue.
	if err != nil {
		return
	}
	c.Next()
}
