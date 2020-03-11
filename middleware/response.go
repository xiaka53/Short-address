package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xiaka53/DeployAndLog/lib"
)

type ResponseCode int

//默认状态码
const (
	SuccessCode    ResponseCode = 200
	ErrorCode      ResponseCode = 500
	ParameterError ResponseCode = 10001 //参数错误
	InternalErrorCode
)

//返回信息格式
type Response struct {
	ErrorCode ResponseCode `json:"errno"`
	ErrorMsg  string       `json:"errmsg"`
	Data      interface{}  `json:"data"`
	TraceId   interface{}  `json:"trace_id"`
}

//返回前端错误信息
func ResponseError(c *gin.Context, code ResponseCode, err error) {
	var (
		trace        interface{}
		traceContext *lib.TraceContext
		traceId      string
		resp         Response
		respone      []byte
	)
	trace, _ = c.Get("trace")
	traceContext, _ = trace.(*lib.TraceContext)
	if traceContext != nil {
		traceId = traceContext.TraceId
	}

	resp = Response{
		ErrorCode: code,
		ErrorMsg:  err.Error(),
		Data:      "",
		TraceId:   traceId,
	}
	c.JSON(200, resp)
	respone, _ = json.Marshal(resp)
	c.Set("response", string(respone))
	_ = c.AbortWithError(200, err)
}

//返回前端信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	var (
		trace        interface{}
		traceContext *lib.TraceContext
		traceId      string
		response     []byte
		resp         Response
	)
	trace, _ = c.Get("trace")
	traceContext, _ = trace.(*lib.TraceContext)
	if traceContext != nil {
		traceId = traceContext.TraceId
	}

	resp = Response{
		ErrorCode: SuccessCode,
		ErrorMsg:  "",
		Data:      data,
		TraceId:   traceId,
	}
	c.JSON(200, resp)
	response, _ = json.Marshal(resp)
	c.Set("response", string(response))
}
