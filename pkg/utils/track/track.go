package track

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"douyin/pkg/utils/log"                      // 本地日志工具包，用于中文日志输出
	"github.com/opentracing/opentracing-go"     // OpenTracing 标准接口
	"github.com/opentracing/opentracing-go/ext" // OpenTracing 扩展，用于设置span属性
	"github.com/uber/jaeger-client-go"          // Jaeger 客户端
	"github.com/uber/jaeger-client-go/config"   // Jaeger 配置包
)

// GetDefaultConfig 返回默认的 Jaeger 配置
// 该配置中使用 const 采样策略，即所有的 span 都被采样；Reporter 会将采样到的 span 发送到本地 Jaeger Agent
func GetDefaultConfig() *config.Configuration {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const", // 固定采样，所有请求都采样
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,             // 输出 span 到日志
			LocalAgentHostPort: "127.0.0.1:6831", // 本地 Jaeger Agent 地址
		},
	}
	return cfg
}

// InitJaeger 初始化 Jaeger 分布式追踪，返回全局 tracer 和一个 io.Closer，用于关闭 tracer
func InitJaeger() (opentracing.Tracer, io.Closer) {
	cfg := GetDefaultConfig()
	// 修改服务名称为 douyin商城，以便在 Jaeger UI 中区分不同服务
	cfg.ServiceName = "douyin商城"
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("初始化 Jaeger 失败: %v", err))
	}
	// 设置全局 tracer
	opentracing.SetGlobalTracer(tracer)
	log.Infof("Jaeger 初始化成功，服务名称：douyin商城")
	return tracer, closer
}

// StartSpan 启动一个顶级 span，用于记录整个操作的追踪信息
func StartSpan(tracer opentracing.Tracer, name string) opentracing.Span {
	span := tracer.StartSpan(name)
	log.Infof("启动顶级 Span: %s", name)
	return span
}

// WithSpan 从上下文中启动一个新的 span，并返回更新后的上下文
func WithSpan(ctx context.Context, name string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	log.Infof("从上下文启动新的 Span: %s", name)
	return span, ctx
}

// GetCarrier 将 span 上下文注入到 HTTP 请求头中，便于跨进程传递追踪信息
func GetCarrier(span opentracing.Span) (opentracing.HTTPHeadersCarrier, error) {
	carrier := opentracing.HTTPHeadersCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		log.Errorf("注入 Span 上下文失败: %s", err.Error())
		return nil, err
	}
	log.Infof("成功注入 Span 上下文到 HTTP headers")
	return carrier, nil
}

// GetParentSpan 从 HTTP 请求头中提取追踪信息，并基于该信息启动一个父 Span
// 该函数用于在服务端接收到请求后，恢复跨服务追踪链路，继续记录调用链
func GetParentSpan(spanName string, traceId string, header http.Header) (opentracing.Span, error) {
	// 构造一个 carrier，用于注入已有的追踪信息
	carrier := opentracing.HTTPHeadersCarrier{}
	carrier.Set("uber-trace-id", traceId)

	tracer := opentracing.GlobalTracer()
	wireContext, err := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header),
	)
	if err != nil {
		log.Errorf("从 HTTP headers 提取追踪上下文失败: %s", err.Error())
		return nil, err
	}

	parentSpan := opentracing.StartSpan(
		spanName,
		ext.RPCServerOption(wireContext), // 将提取到的上下文作为父上下文
	)
	log.Infof("启动父 Span: %s", spanName)
	return parentSpan, nil
}
