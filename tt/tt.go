package tt

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/xybstone/xservice/service"
)

func init() {
	service.NewUniqueIDAsync()
}

//Convey 执行
func Convey(info string, t *testing.T, f func() bool) {
	fmt.Println("Test:", info)
	if !f() {
		t.Fatalf("fail:%s", info)
	}
}

//SkipConvey 跳过测试
func SkipConvey(info string, t *testing.T, f func() bool) {
	fmt.Println("Skip Test:", info)
}

// StartServerOnPort 开始测试
func StartServerOnPort(t *testing.T, port int, h fasthttp.RequestHandler) io.Closer {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("cannot start tcp server on port %d: %s", port, err)
	}
	go fasthttp.Serve(ln, h)
	return ln
}

var header http.Header
var Header = func() http.Header {
	if header == nil {
		header = make(map[string][]string)
	}
	return header
}()

//DoHTTPTestURI http
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

	for k, v := range Header {
		for _, vv := range v {
			if _, has := req.Header[k]; has {
				req.Header.Set(k, vv)

			} else {
				req.Header.Add(k, vv)
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return body
}

//DoHTTPTestBody body
func DoHTTPTestBody(method string, port int, api string, body []byte) []byte {
	if strings.HasPrefix(api, "/") {
		api = strings.TrimPrefix(api, "/")
	}

	var req *http.Request
	var err error
	var host string
	if !strings.HasPrefix("http", api) {
		host = fmt.Sprintf("http://localhost:%d", port)
	}

	req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", host, api), bytes.NewReader(body))

	if err != nil {
		panic(err)
	}

	for k, v := range Header {
		for _, vv := range v {
			if _, has := req.Header[k]; has {
				req.Header.Set(k, vv)

			} else {
				req.Header.Add(k, vv)
			}
		}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	outBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(outBody))
	return outBody
}