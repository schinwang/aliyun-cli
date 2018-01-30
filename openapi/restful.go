package openapi

import (
	"github.com/aliyun/aliyun-cli/meta"
	"strings"
	"github.com/aliyun/aliyun-cli/cli"
	"fmt"
	"time"
	"io/ioutil"
)

func (c *Caller) InvokeRestful(ctx *cli.Context, product *meta.Product, method string, path string) {
	client, request, err := c.InitClient(ctx, product, false)

	request.Headers["Date"] = time.Now().Format(time.RFC1123Z)
	request.PathPattern = path
	request.Method = method

	if v, ok := ctx.Flags().GetValue("body"); ok {
		request.SetContent([]byte(v))
	}
	if v, ok := ctx.Flags().GetValue("body-file"); ok {
		buf, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Errorf("failed read file: %s %v", v, err)
		}
		request.SetContent(buf)
	}

	resp, err := client.ProcessCommonRequest(request)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unmarshal") {
			// fmt.Printf("%v\n", err)
		} else {
			ctx.Command().PrintFailed(err, "---")
		}
	}
	fmt.Println(resp.GetHttpContentString())
}

func CheckRestfulMethod(ctx *cli.Context, methodOrPath string, pathPattern string) (ok bool, method string, path string, err error) {
	if method, ok = CheckHttpMethod(methodOrPath); ok {
		if strings.HasPrefix(pathPattern, "/") {
			path = pathPattern
			return
		} else {
			err = fmt.Errorf("bad restful path %s", pathPattern)
			return
		}
	} else if method, ok = ctx.Flags().GetValue("roa"); ok {
		if strings.HasPrefix(methodOrPath, "/") && pathPattern == "" {
			path = methodOrPath
			return
		} else {
			err = fmt.Errorf("bad restful path %s", methodOrPath)
			return
		}
	} else {
		ok = false
		return
	}
}

func CheckHttpMethod(s string) (string, bool) {
	m := strings.ToUpper(s)
	if m == "GET" || m == "POST" || m == "PUT" || m == "DELETE" {
		return m, true
	}
	return "", false
}
