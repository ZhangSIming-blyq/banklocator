package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for authentication information in the request headers or cookies
		authToken, exist := c.GetQuery("auth")

		if exist && authToken != "zhangsiminghahahaha" {
			// If authentication information is missing or invalid,
			// redirect the user to the login page
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else if !exist {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		c.Next()
	}
}
