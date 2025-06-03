package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"douyin/consts"
	"douyin/pkg/utils/track"
	"github.com/opentracing/opentracing-go/ext" // Added for opentracing ext tags
)

// Jaeger 中间件：用于集成 Jaeger 分布式追踪
func Jaeger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取追踪ID
		traceId := c.GetHeader("uber-trace-id")
		var span opentracing.Span
		if traceId != "" {
			var err error
			// 如果请求中携带traceId，则根据该ID创建父span
			span, err = track.GetParentSpan(c.FullPath(), traceId, c.Request.Header)
			if err != nil {
				// 输出中文错误提示到控制台，并返回
				// Consider using c.Error() and letting the ErrorHandler middleware handle the response
				// For now, keeping original behavior for this specific error path.
				c.String(500, "获取父追踪失败：%s", err.Error())
				return
			}
		} else {
			// 否则，启动一个新的span
			span = track.StartSpan(opentracing.GlobalTracer(), c.FullPath())
		}
		// 在处理结束后关闭span
		defer span.Finish()

		// Set standard OpenTracing tags
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.RequestURI) // Using RequestURI for full path + query
		ext.Component.Set(span, "gin-http-server")
		// Optionally, set client IP
		// ext.PeerHostIPv4.SetString(span, c.ClientIP())


		// 将span存入上下文，便于后续使用
		c.Set(consts.SpanCTX, opentracing.ContextWithSpan(c, span))
		c.Next()
	}
}
