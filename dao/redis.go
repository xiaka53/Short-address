package dao

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/mattheath/base62"
	"github.com/speps/go-hashids"
	"github.com/xiaka53/DeployAndLog/lib"
	"time"
)

const (
	UrlIdKey           = "next.url.id"         //自增ID，唯一
	ShortLinkKey       = "shortLink:%s:url"    //短域名->hash信息
	UrlHashKey         = "urlHash:%s:url"      //hash->短域名
	ShortLinkDetailKey = "shortLink:%s:detail" //短域名->域名详细信息
)

type RedisCli struct {
	PoolName string
	Db       int //数据库
}

//url生成短域名
func (r *RedisCli) Shorten(trace *lib.TraceContext, url string, exp int64) (string, error) {
	urlHash := toHash(url)
	rs, err := lib.RedisConnFactory(r.PoolName)
	if err != nil {
		return "", err
	}
	defer lib.RedisConnClose(trace, rs)
	if _, err := lib.RedisLogDo(trace, rs, "select", r.Db); err != nil {
		return "", err
	}
	d, err := redis.String(lib.RedisLogDo(trace, rs, "get", fmt.Sprintf(UrlHashKey, urlHash)))
	if err != nil && err != redis.ErrNil {
		return "", err
	} else {
		if d == "" {
			//空数据
		} else {
			return d, nil
		}
	}
	if _, err := lib.RedisLogDo(trace, rs, "incr", UrlIdKey); err != nil {
		return "", err
	}
	id, err := redis.Int64(lib.RedisLogDo(trace, rs, "get", UrlIdKey))
	if err != nil {
		return "", err
	}
	shortLink := base62.EncodeInt64(id)
	exp *= 60 * 60 * 24
	if _, err := lib.RedisLogDo(trace, rs, "setex", fmt.Sprintf(ShortLinkKey, shortLink), exp, url); err != nil {
		return "", err
	}

	if _, err := lib.RedisLogDo(trace, rs, "setex", fmt.Sprintf(UrlHashKey, urlHash), exp, shortLink); err != nil {
		return "", nil
	}
	detail, err := json.Marshal(&UrlDetail{
		Url:                 url,
		CreateAt:            time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})
	if err != nil {
		return "", err
	}
	if _, err := lib.RedisLogDo(trace, rs, "setex", fmt.Sprintf(ShortLinkDetailKey, shortLink), exp, detail); err != nil {
		return "", err
	}
	return shortLink, nil
}

//url转换hash
func toHash(url string) string {
	hd := hashids.NewData()
	hd.Salt = url
	hd.MinLength = 0
	h, _ := hashids.NewWithData(hd)
	r, _ := h.Encode([]int{45, 434, 1313, 99})
	return r
}

//定义url详细信息
type UrlDetail struct {
	Url                 string        `json:"url"`
	CreateAt            string        `json:"create_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

//短链接获取域名信息
func (r *RedisCli) ShortLinkInfo(trace *lib.TraceContext, shortLink string) (interface{}, error) {
	rs, err := lib.RedisConnFactory(r.PoolName)
	if err != nil {
		return "", err
	}
	defer lib.RedisConnClose(trace, rs)
	if _, err := lib.RedisLogDo(trace, rs, "select", r.Db); err != nil {
		return "", err
	}
	detail, err := lib.RedisLogDo(trace, rs, "get", fmt.Sprintf(ShortLinkDetailKey, shortLink))
	if err != nil {
		return "", err
	} else {
		return detail, nil
	}
}

//短链接获取原始链接
func (r *RedisCli) UnShorten(trace *lib.TraceContext, shortLink string) (string, error) {
	rs, err := lib.RedisConnFactory(r.PoolName)
	if err != nil {
		return "", err
	}
	defer lib.RedisConnClose(trace, rs)
	if _, err := lib.RedisLogDo(trace, rs, "select", r.Db); err != nil {
		return "", err
	}
	url, err := redis.String(lib.RedisLogDo(trace, rs, "get", fmt.Sprintf(ShortLinkKey, shortLink)))
	if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

//初始化redis配置
func InitRedis() *RedisCli {
	return &RedisCli{
		Db:       0,
		PoolName: "default",
	}
}
