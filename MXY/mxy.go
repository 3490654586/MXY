package MXY

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string  //前缀
	middlewares []HandlerFunc  //支持中间件
	parent      *RouterGroup   //支持嵌套
	engine      *Engine        //所有引擎
}

//Engine实现了ServerHttp接口 router参数为路由表
type  Engine struct {
	Post string
	Host string
	*RouterGroup
	router *router
	groups []*RouterGroup
}


//获取mxy.Engine
func New() * Engine{
	//使用分组的情况下将路由函数交给RouterGroup实现
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default()* Engine {
	engine := New()
	engine.Use(Logger())
	return engine
}

//创建一个新的RouterGroup,所有groups恭喜一个engine引擎
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//注册GET请求
func (engine *Engine)GET(pattern string,handler HandlerFunc){
         engine.router.addRouter("GET",pattern,handler)
}

//注册POST请求
func (engine *Engine)POST(pattern string,handler HandlerFunc){
	engine.router.addRouter("POST",pattern,handler)
}

//封装ListenAndServe方法
func (engince *Engine)Run()error{
	if engince.Post == ""{
		engince.Post="localhost"
		engince.Host = "699"
	}

	debugPrint("Listening and serving HTTP on %s:%s\n", engince.Post,engince.Host)
	return http.ListenAndServe(engince.Post+":"+engince.Host,engince)
}

func debugPrint(format string, values ...interface{}) {
		fmt.Fprintf(os.Stdout, format, values...)
}
/**
让Engine实现ServeHTTP方法才能注册服务函数
当我们接收到一个具体请求时，要判断该请求适用于哪些中间件
 */
func (engine *Engine)ServeHTTP(w http.ResponseWriter,req *http.Request)  {

	var middlewares []HandlerFunc
	//循环分组数组
	for _, group := range engine.groups {
		//根据路由前缀判断是否走的分组
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			//添加中间件函数
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c:= NewContext(w,req)
	//将中件间添加至上下文
	c.handlers = middlewares
	engine.router.handle(c)
}

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		fmt.Println(11)
		c.Next()
		log.Printf("| %d |  %s |  %v  | %v     %v ", c.StatusCode, time.Since(t),ClientIP(c.Req),c.Req.Method,c.Path)
	}
}

//获取客户端ip
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
