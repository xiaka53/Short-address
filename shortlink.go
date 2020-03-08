package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

//路由
type App struct {
	Router      *mux.Router
	Middlewares *Middleware
}

//接受参数
type shortenReq struct {
	URL                  string `json:"url" validate:"nonzero"`
	ExpirationInMinuxtes int64  `json:"expiration_in_minuxtes" validate:"min=0"`
}

//返回参数
type shortLintResp struct {
	ShortLink string `json:"short_link"`
}

//初始化服务
func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.Middlewares = &Middleware{}
	app.initializeRouter()
}

//路由绑定
func (app *App) initializeRouter() {
	m := alice.New(app.Middlewares.LoggingHandler, app.Middlewares.RecoverHandler)
	app.Router.Handle("/api/shorter", m.ThenFunc(app.createShortlink)).Methods("POST")
	app.Router.Handle("/api/info", m.ThenFunc(app.getShortlinkInfo)).Methods("GET")
	app.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(app.rediect)).Methods("GET")
}

//创建短链接
func (app *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, NewBadReqErr(fmt.Errorf("%s", "happen err when parsing json from body")), nil)
		return
	}

	defer r.Body.Close()

	if err := validator.Validate(req); err != nil {
		responseWithError(w, NewBadReqErr(fmt.Errorf("validate parameters failed : %+v", req)), nil)
		return
	}
	//if err := validator.Validate{}
	fmt.Println(req)
}

//端地址解析
func (app *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	fmt.Println("getShortlinkInfo:", s)
}

//重定向
func (app *App) rediect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("rediect:", vars["shortlink"])
}

//启动服务
func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

//浏览器输出错误信息
func responseWithError(w http.ResponseWriter, err error, payload interface{}) {
	switch e := err.(type) {
	case MiErr:
		resp, _ := json.Marshal(Response{Code: e.Status(), Message: e.Error(), Content: payload})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	default:
		responseWithJson(w, http.StatusInternalServerError, payload)
	}
}

//默认输出系统错误
func responseWithJson(w http.ResponseWriter, status int, payload interface{}) {
	resp, _ := json.Marshal(Response{Code: status, Message: http.StatusText(status), Content: payload})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

//返回浏览器信息格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}
