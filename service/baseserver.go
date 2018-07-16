package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	routing "github.com/xybstone/fasthttp-routing"
	"gitlab.gaodun.com/golib/filetool"
)

//BaseResponse  返回参数
type BaseResponse struct {
	HTTPCode int         `json:"http_code,omitempty"`
	Status   int         `json:"status"`
	Info     string      `json:"info"`
	Result   interface{} `json:"result"`
}

// BaseServer 服务基类
type BaseServer struct {
	Logger func(title string, msg map[string]interface{})
}

// GetRouterMap 路由表
func (bs *BaseServer) GetRouterMap() map[string]Route {
	AddRote("/", Route{[]string{"GET", "POST"}, bs.HandleRoot})
	AddRote("/status", Route{[]string{"GET"}, bs.GetStatus})
	AddRote("/cors", Route{[]string{"GET"}, bs.Cors})
	return routes
}

// Cors 跨域
func (bs *BaseServer) Cors(ctx *routing.Context) (err error) {
	ctx.Write([]byte(`{"http_code":200,"status":0,"info":"NOT_ALLOW_CORS","result":null}`))
	return
}

//HandleRoot 测试页
func (bs *BaseServer) HandleRoot(ctx *routing.Context) (err error) {
	v := fmt.Sprintf("{\"version\":\"%s\"}", ServiceVersion)
	ctx.Write([]byte(v))
	return
}

//GetStatus 状态页
func (bs *BaseServer) GetStatus(ctx *routing.Context) (err error) {
	dir, _ := os.Getwd()
	fileOperate := filetool.FileOperate{Filename: dir + string(os.PathSeparator) + "DEPLOY"}
	if fileOperate.CheckFileExist(fileOperate.Filename) {
		text, _ := fileOperate.ReadFile(fileOperate.Filename)
		lineList := strings.Split(string(text), "\n")
		begin := "{"
		tmpStr := ""
		for _, line := range lineList {
			if len(line) == 0 {
				break
			}
			n := strings.Split(line, "|")
			jsonKey := "\"" + n[0] + "\""
			jsonValue := ":\"" + n[1] + "\""
			tmpStr += jsonKey + jsonValue + ","
		}
		mStr := strings.TrimRight(tmpStr, ",")
		end := "}"
		info := begin + mStr + end
		ctx.Write([]byte("{\"status\": \"1\",\"data\":" + info + "}"))
		return
	}

	ctx.Write([]byte("{\"status\": \"1\"}"))
	return
}

var routes map[string]Route

// BindRote 注册绑定路由
func BindRote(path string, methods []string, handle func(*routing.Context) error) {
	AddRote(path, Route{method: methods, handle: handle})
}

// AddRote 添加路由
func AddRote(path string, r Route) {
	if routes == nil {
		routes = make(map[string]Route)
	}
	routes[path] = r
}

//ServerJSON 服务器返回
func (bs *BaseServer) ServerJSON(ctx *routing.Context, v interface{}, status int) {
	if b, err := json.Marshal(v); err == nil {
		if status == 0 {
			XLogger.Printf("ServerJSON:%v", map[string]interface{}{"Respose OK": string(b)})
		} else {
			XLogger.Printf("ServerJSON ERR:%v", map[string]interface{}{"Respose ERR": string(b)})
		}
		ctx.Write(b)
	}
}

// FormString 获取字符串
func (bs *BaseServer) FormString(ctx *routing.Context, key string) string {
	return strings.TrimSpace(string(ctx.FormValue(key)))
}

// FormInt64 获取int
func (bs *BaseServer) FormInt64(ctx *routing.Context, key string) int64 {
	i, _ := strconv.ParseInt(string(ctx.FormValue(key)), 10, 64)
	return i
}
