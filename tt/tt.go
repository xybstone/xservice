package tt

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

//Convey 执行
func Convey(info string, f func(), at func() bool) {
	fmt.Println("Test:", info)
	f()
	if at != nil {
		if !at() {
			panic("fail")
		}
	}
}

//SkipConvey 跳过测试
func SkipConvey(info string, f func(), at func() bool) {
	fmt.Println("Skip Test:", info)
}

func StartServerOnPort(t *testing.T, port int, h fasthttp.RequestHandler) io.Closer {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("cannot start tcp server on port %d: %s", port, err)
	}
	go fasthttp.Serve(ln, h)
	return ln
}

func DoHTTPTestURI(method string, port int, api string, params map[string]string) []byte {
	if strings.HasPrefix(api, "/") {
		api = strings.TrimPrefix(api, "/")
	}

	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	var req *http.Request
	var err error
	var host string
	if !strings.HasPrefix("http", api) {
		host = fmt.Sprintf("http://localhost:%d", port)
	}

	if len(params) > 0 {
		req, err = http.NewRequest(method, fmt.Sprintf("%s/%s?%s", host, api, vals.Encode()), nil)
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", host, api), nil)
	}
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	return body
}
