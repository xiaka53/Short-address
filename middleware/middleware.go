package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaka53/DeployAndLog/lib"
)

//type Middleware struct {
//}

//输出请求时间
//func (m Middleware) LoggingHandler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		t1 := time.Now()
//		next.ServeHTTP(w, r)
//		t2 := time.Now()
//		log.Printf("[%s] %q %v", r.Method, r.URL.String(), t2.Sub(t1))
//	})
//}

//输出系统错误请求
//func (m Middleware) RecoverHandler(next http.Handler) http.Handler {
//	fn := func(w http.ResponseWriter, r *http.Request) {
//		defer func() {
//			if err := recover(); err != nil {
//				log.Printf("Recover from panic:%+v", err)
//				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//			}
//		}()
//		next.ServeHTTP(w, r)
//	}
//	return http.HandlerFunc(fn)
//}

//白名单过滤
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _, host := range lib.GetStringSliceConf("base.http.allow_ip") {
			if c.ClientIP() == host || host == "*" {
				isMatched = true
			}
		}
		if !isMatched {
			ResponseError(c, InternalErrorCode, errors.New(fmt.Sprintf("%v, not in iplist", c.ClientIP())))
			c.Abort()
			return
		}
		c.Next()
	}
}
