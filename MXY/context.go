package MXY

import (
	"MXY-WEB/MXY/binding"
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path string
	Method string
	Params map[string]string
	StatusCode int
	//中间件
	handlers []HandlerFunc
	index    int
}


func NewContext(w http.ResponseWriter,req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		index:-1,//默认没有
	}
}
func (c *Context) Next() {
	//记录执行第几个中间件
	c.index++
	//判断存储中间件的函数数据
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		//执行中间件
		c.handlers[c.index](c)
	}
}

//获取路由参数接下到Params中
func (c *Context)Param(key string)string  {
	  value,_:= c.Params[key]
	  return value
}
//获取post参数
func (c *Context)PostForm(key string)string{
	return c.Req.FormValue(key)
}

//获取GET参数
func (c *Context)Query(key string)string{
	return c.Req.URL.Query().Get(key)
}

//设置Status值
func (c *Context)Status(code int){
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//设置请求头编码格式
func (c *Context)SetHeader(key string,value string)  {
	c.Writer.Header().Set(key,value)
}

//设置String返回编码
func (c *Context)String(code int,format string,values ...interface{}){
	 c.SetHeader("Content-Type","text/plain")
	 c.Status(code)
	 //返回主体
	 c.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

//设置JSON返回格式
func (c *Context)JSON(code int,obj interface{})  {
	  c.SetHeader("Content-Type","application/json")
	  c.Status(code)
	  encoder := json.NewEncoder(c.Writer)
	 err :=  encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer,err.Error(),500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//设置返回HTML格式
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

//将请求参数解析到结构体
func (c *Context)ShouldBind(obj interface{})error{
 b :=	binding.Default(c.Req.Header.Get("Content-Type"))
 return c.ShouldBindWith(obj,b)
}


func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Req, obj)
}