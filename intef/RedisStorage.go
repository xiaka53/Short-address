package intef

import "github.com/xiaka53/DeployAndLog/lib"

type RedisStorage interface {
	//将长地址转为短地址
	Shorten(trace *lib.TraceContext, url string, exp int64) (string, error)
	ShortLinkInfo(trace *lib.TraceContext, eid string) (interface{}, error)
	UnShorten(trace *lib.TraceContext, eid string) (string, error)
}
