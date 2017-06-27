package ClaremontMenuGo

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "HELLO WORLD!")
	})

	router.Run(":8080")
}
