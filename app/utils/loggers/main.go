package loggers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var JSON = NewLogger("json")
var Text = NewLogger("text")

func New() gin.HandlerFunc {
	log := NewLogger("json")

	return func(c *gin.Context) {
		start := time.Now()

		method := c.Request.Method
		uri := c.Request.RequestURI

		ip := c.ClientIP()

		entry := map[string]interface{}{
			"user_ip":    ip,
			"user_agent": c.Request.UserAgent(),
			"referrer":   c.Request.Referer(),
			"method":     method,
			"path":       uri,
			"start_time": start.Format(time.RFC3339),
		}
		if len(c.Errors) > 0 {
			log.WithFields(entry).Error(c.Errors.String())
		} else {
			title := fmt.Sprintf(`%s | %s "%s"`, ip, method, uri)
			log.WithFields(entry).Info(title)
		}
		c.Next()
	}
}
