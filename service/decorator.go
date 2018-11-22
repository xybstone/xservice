package service

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strings"

	routing "github.com/xybstone/fasthttp-routing"
)

var (
	CorsAllowOrigin      = []string{"gaodun.com","gaodunwangxiao.com"}
	CorsAllowMethods     = "HEAD,GET,POST,PUT,DELETE"
	CorsAllowCredentials = "true"
	CorsAllowHead        = "Origin, X-Requested-With, Content-Type, Accept, Authtoken, Authentication, X_Requested_With, X-Request-ID"
	ExposeHeaders        = "AccessToken, RefreshToken"
)

func init() {
	Dct.SetVerifyList("/status")
	Dct.SetVerifyList("/cors")
	Dct.SetVerifyList("/")
}

var (
	OPTIONS = []byte("OPTIONS")
)

//Decorator 控制拦截
type Decorator struct {
	RunFuc  func(ctx *routing.Context) error
	PathStr string
}

var verifyList map[string]bool

func (d Decorator) GetVerifyList() map[string]bool {
	return verifyList
}

func (d Decorator) SetVerifyList(key string) {
	if verifyList == nil {
		verifyList = make(map[string]bool)
	}

	verifyList[key] = true
}

func (d Decorator) SetCorsHeader(ctx *routing.Context) {
	origin := string(ctx.Request.Header.PeekBytes([]byte("origin")))
	if origin == "" {
		origin = string(ctx.Host())
	}
	allow := false
	for _, v := range CorsAllowOrigin {
		if strings.Contains(origin, v) || verifyList[string(ctx.Path())] {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", CorsAllowCredentials)
			ctx.Response.Header.Set("Access-Control-Allow-Methods", CorsAllowMethods)
			ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
			ctx.Response.Header.Set("Access-Control-Allow-Headers", CorsAllowHead)
			ctx.Response.Header.Set("Access-Control-Expose-Headers", ExposeHeaders)
			reqID := string(ctx.Request.Header.Peek("X-Request-ID"))
			if reqID == "" {
				reqID = GetUID()
				ctx.Request.Header.Set("X-Request-ID", reqID)
			}
			allow = true
			break
		}
	}
	if !allow {
		ctx.Redirect("/cors", 302)
	}
}

func (d Decorator) recovery() {
	if r := recover(); r != nil {
		fmt.Println("XLogger:", XLogger)
		XLogger.Printf("panic:%v", map[string]interface{}{"recover": r, "stack": string(debug.Stack())})
	}
}

func (d Decorator) Decorator(ctx *routing.Context) (err error) {
	defer d.recovery()
	d.SetCorsHeader(ctx)

	fmt.Println("im here1")

	if bytes.Compare(OPTIONS, ctx.Method()) == 0 {
		return
	}

	fmt.Println("im here")

	return d.RunFuc(ctx)
}
