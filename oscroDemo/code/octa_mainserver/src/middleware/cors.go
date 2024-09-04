package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	c := cors.Config{
		//AllowOrigins: []string{"http://192.168.1.178:9099"},
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders: []string{
			"Content-Type",
			"Access-Token",
			"Authorization",
			"serverName",
			"token",
		},
		MaxAge: 6 * time.Hour,
	}

	return cors.New(c)
}
