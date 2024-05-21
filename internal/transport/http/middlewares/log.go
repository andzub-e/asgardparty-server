package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

// Log is a gin middleware for logging requests and responses.
func Log(logger *zap.SugaredLogger, timeFormat string, utc bool) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()

		var requestBodyBytes []byte
		if c.Request.Body != nil {
			requestBodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}

		responseBodyWriter := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBodyWriter

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if utc {
			end = end.UTC()
		}

		requestHeaders := ""
		for name, values := range c.Request.Header {
			requestHeaders += name + ": " + values[0] + "; "
		}

		responseHeaders := ""
		for name, values := range responseBodyWriter.Header() {
			responseHeaders += name + ": " + values[0] + "; "
		}

		// Log the request and response details
		logger.Infow("request and response logged",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"time", end.Format(timeFormat),
			"latency", latency,
			"request_headers", requestHeaders,
			"request_body", string(requestBodyBytes),
			"response_headers", responseHeaders,
			"response_body", responseBodyWriter.body.String(),
		)
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
