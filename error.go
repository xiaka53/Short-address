package main

import "net/http"

//定义错误接口
type MiErr interface {
	error
	Status() int
}

type ReqError struct {
	Code int
	Err  error
}

//错误内容
func (re ReqError) Error() string {
	return re.Err.Error()
}

//错误码
func (re ReqError) Status() int {
	return re.Code
}

//创建新的404错误信息
func NewNotFindErr(err error) ReqError {
	return ReqError{
		Code: http.StatusNotFound,
		Err:  err,
	}
}

//创建400错误信息
func NewBadReqErr(err error) ReqError {
	return ReqError{
		Code: http.StatusBadRequest,
		Err:  err,
	}
}
