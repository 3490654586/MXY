package MXY

import (
	"fmt"
	"net/http"
	"strings"
)
//路由结构体
type router struct {
	rootTrie map[string]*Node
	handlers map[string]HandlerFunc
}


//返回router实例
func NewRouter()*router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

//进行路由分割
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router)addRouter(method string , pattern string,handler HandlerFunc){
	parts := parsePattern(pattern)
	key := method +"-"+pattern;
	fmt.Println(key)
	_,ok := r.rootTrie[method]
	//如果查询不到此方法的树,就为次方法新增一棵树
   if !ok{
   	r.rootTrie = make(map[string]*Node)
   	r.rootTrie[method] = &Node{}
   }
   //为次方法添加路由
   r.rootTrie[method].insert(pattern,parts,0)
	r.handlers[key] = handler
}

//获取路由函数
func (r *router) getRoute(method string, path string) (*Node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	//通过请求方法查询路由树
	Node, ok := r.rootTrie[method]

	if !ok {
		return nil, nil
	}
    //查询路由树
	n := Node.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		fmt.Println("parts",parts)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		fmt.Println("params=",params)
		return n, params
	}

	return nil, nil
}

//注册路由函数
func (r *router)handle(c *Context)  {
	//通过路由树查找路由函数
	n,params :=r.getRoute(c.Method,c.Path)
	if n != nil{

		//根据请求方法和请求路径组成key
		key :=c.Method +"-"+c.Path
		c.Params = params
		//从路由函数根据key找出函数添加到中间件数组
		c.handlers = append(c.handlers,r.handlers[key])
	}else {
		//没找到返回404,并将这个错误路由处理函数添加到中间件数组
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404", c.Path)
		})
	}
	//接触处理下一个请求
  c.Next()
}