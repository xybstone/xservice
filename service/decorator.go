package service

import (
	"bytes"
	"runtime/debug"
	"strings"

	routing "github.com/xybstone/fasthttp-routing"
)

var (
	CorsAllowOrigin      = []string{"gaodun.com"}
	CorsAllowMethods     = "HEAD,GET,POST,PUT,DELETE"
	CorsAllowCredentials = "true"
	CorsAllowHead        = "Origin, X-Requested-With, Content-Type, Accept, Authtoken, Authentication, X_Requested_With"
	ExposeHeaders        = "AccessToken, RefreshToken"
)

var (
	OPTIONS = []byte("OPTIONS")
)

//Decorator 控制拦截
type Decorator struct {
	RunFuc  func(ctx *routing.Context) error
	PathStr string
	Logger  func(title string, msg map[string]interface{})
}

var verifyList map[string]bool

func (d Decorator) GetVerifyList() map[string]bool {
	return verifyList
}

func (d Decorator) SetVerifyList(key string) {
	if verifyList == nil {
		verifyList = make(map[string]bool)
		verifyList["/"] = true
		verifyList["/status"] = true
		verifyList["/cors"] = true
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
		d.Logger("panic", map[string]interface{}{"recover": r, "stack": string(debug.Stack())})
	}
}

func (d Decorator) Decorator(ctx *routing.Context) (err error) {
	defer d.recovery()
	d.SetCorsHeader(ctx)

	if bytes.Compare(OPTIONS, ctx.Method()) == 0 {
		return
	}

	d.RunFuc(ctx)
	return
}
