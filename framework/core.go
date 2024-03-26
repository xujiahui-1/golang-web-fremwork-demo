package framework

import (
	"log"
	"net/http"
)

// Core 核心结构体
type Core struct {
	router map[string]ControllerHandler
}

func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

// ServerHTTP 创建服务的用户侧方法
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	//TODO
	log.Println("core.serverHTTP")
	ctx := NewContext(request, response)
	//路由选择器
	router := c.router["foo"]
	if router == nil {
		return
	}
	log.Println("core.router")

	router(ctx)
}
