// Authored By @SixExtreme at 2018.12.11
// 早期的路由注册实现为每个 Handler 承担多种 HTTP Method 处理
// 扩展的 Resource in Group 机制使得每种HTTP Method都有一个Hndler负责处理
// 参考了开源社区常见的路由注册实现方式，用路由器添加路由, 不是路由注册到路由器

package service

import (
	"encoding/json"
	"github.com/xybstone/fasthttp-routing"
)

type (
	// RESTful 资源
	// 每种HTTP Method对应一个Handler
	Resource map[string]routing.Handler

	// 实现此接口的类可以被注册路由
	// 返回 path -> resource 注册表
	IResourceGrouper interface {
		GetResourceMap() map[string]Resource
	}
)

// 调用此接口将一组API注册到总路由器
func RegisterGroup(router *routing.Router, group IResourceGrouper) {
	for path, resoure := range group.GetResourceMap() {
		for method, handler := range resoure {
			router.To(method, path, handler)
		}
	}
}


// 示例Group
type BaseGroup struct {}

// 实现IResourceGrouper接口
func (bg BaseGroup) GetResourceMap() map[string]Resource {
	return map[string]Resource{
		"/author": map[string]routing.Handler{
			"GET": bg.GetAuthor,
			"PUT": bg.PutAuthor,
		},
	}
}

// 示例API, 返回作者信息
func (bg BaseGroup) GetAuthor(ctx *routing.Context) (err error) {
	data := map[string]string{
		"master":  "bob.xu",
	}
	jdata, _ := json.Marshal(data)
	_, err = ctx.Write(jdata)
	return
}

func (bg BaseGroup) PutAuthor(ctx *routing.Context) (err error) {
	data := map[string]string{
		"master":  "bob.xu",
		"slave": "six extreme",
	}
	jdata, _ := json.Marshal(data)
	_, err = ctx.Write(jdata)
	return
}