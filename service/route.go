package service

import routing "github.com/xybstone/fasthttp-routing"

// Route 路由
type Route struct {
	method []string
	handle func(*routing.Context) error
}

// IRegistRouter 注册路由接口
type IRegistRouter interface {
	GetRouterMap() map[string]Route
}

var Dct Decorator

func RegsitRouter(router *routing.Router, regRouter IRegistRouter) {
	routers := regRouter.GetRouterMap()
	for k, v := range routers {
		v.method = append(v.method, "OPTIONS")
		for _, value := range v.method {
			Dct.PathStr = k
			Dct.RunFuc = v.handle
			router.To(value, k, Dct.Decorator)
		}
	}
}
