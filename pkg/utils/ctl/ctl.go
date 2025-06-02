package ctl

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"

	"douyin/pkg/error_code" // 本地错误码包，存放错误码和提示信息
	"douyin/pkg/utils/log"  // 本地日志工具包，用于日志记录
)

// Response 基础响应结构体，用于封装 HTTP 响应数据
type Response struct {
	Status  int         `json:"status"`   // 状态码
	Data    interface{} `json:"data"`     // 返回数据
	Msg     string      `json:"msg"`      // 提示信息
	Error   string      `json:"error"`    // 错误信息
	TrackId string      `json:"track_id"` // 请求追踪 ID
}

// TrackedErrorResponse 带追踪信息的错误响应结构体
type TrackedErrorResponse struct {
	Response
	TrackId string `json:"track_id"` // 请求追踪 ID（重复字段，为了兼容前端数据结构）
}

// RespSuccess 返回成功响应，并输出中文日志
func RespSuccess(ctx *gin.Context, data interface{}, code ...int) *Response {
	// 从上下文中获取追踪 ID
	trackId, _ := getTrackIdFromCtx(ctx)
	status := error_code.SUCCESS // 默认成功状态码
	if len(code) > 0 {
		status = code[0]
	}
	if data == nil {
		data = "操作成功"
	}
	// 通过 error_code 包获取状态对应的提示信息
	msg := error_code.GetMsg(status)
	r := &Response{
		Status:  status,
		Data:    data,
		Msg:     msg,
		TrackId: trackId,
	}
	// 控制台输出中文日志
	log.Infof("响应成功: 状态码=%d, 追踪ID=%s", status, trackId)
	return r
}

// RespError 返回错误响应，并输出中文日志
func RespError(ctx *gin.Context, err error, data string, code ...int) *TrackedErrorResponse {
	// 从上下文中获取追踪 ID
	trackId, _ := getTrackIdFromCtx(ctx)
	status := error_code.ERROR // 默认错误状态码
	if len(code) > 0 {
		status = code[0]
	}
	// 通过 error_code 包获取状态对应的提示信息
	msg := error_code.GetMsg(status)
	r := &TrackedErrorResponse{
		Response: Response{
			Status: status,
			Msg:    msg,
			Data:   data,
			Error:  err.Error(),
		},
		TrackId: trackId,
	}
	// 控制台输出中文错误日志
	log.Errorf("响应错误: 状态码=%d, 错误信息=%s, 追踪ID=%s", status, err.Error(), trackId)
	return r
}

// getTrackIdFromCtx 从 gin.Context 中提取追踪 ID
func getTrackIdFromCtx(ctx *gin.Context) (string, error) {
	spanCtxInterface, exists := ctx.Get(error_code.SpanCTX) // 使用 error_code 包中定义的追踪上下文 key
	if !exists {
		return "", errors.New("未找到追踪上下文")
	}
	str := fmt.Sprintf("%v", spanCtxInterface)
	// 使用正则表达式匹配 16 位的追踪 ID
	re := regexp.MustCompile(`([0-9a-fA-F]{16})`)
	match := re.FindStringSubmatch(str)
	if len(match) > 0 {
		return match[1], nil
	}
	return "", errors.New("获取追踪ID失败")
}

// CodedError 定义一个带错误码的接口，便于未来扩展错误处理（例如业务错误、参数校验错误等）
type CodedError interface {
	Code() int
	Error() string
}

// ErrorResponse 统一错误响应处理函数
// 该函数根据传入错误构造标准错误响应并返回
// 调用方可以通过返回值直接传递给 ctx.JSON，同时调用 ctx.Abort() 终止后续处理
func ErrorResponse(ctx *gin.Context, err error) *TrackedErrorResponse {
	if err == nil {
		return nil
	}

	// 默认使用错误码 ERROR
	status := error_code.ERROR
	// 如果错误实现了 CodedError 接口，则可以获取具体的错误码
	if ce, ok := err.(CodedError); ok {
		status = ce.Code()
	}

	// 构造错误响应，此处 data 字段默认为 "操作失败"
	resp := RespError(ctx, err, "操作失败", status)
	// 输出统一错误日志
	log.Errorf("统一错误响应已构造: %s", err.Error())
	// 终止后续处理，避免重复写入响应
	ctx.Abort()
	return resp
}
