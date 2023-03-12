package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/main_system/platform/hub"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/controllers"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// NOTE:
	// Q: 主系统、platform、插件都有自己的路由系统,但是三者的入口都是主系统,如何分辨?
	// A: 通过route这个header明确【最终】要路由到哪里
	//     1. 无      : 路由到主系统
	//     2. plugin  : 路由到插件的自定义路由函数
	//     3. external: 路由到插件的OnExternalHttpRequest函数
	//     4. platform: 路由到platform的路由(一般是插件的生命周期函数)
	routeKey         = "route"
	routeValPlugin   = "plugin"
	routeValExternal = "external"
	routeValPlatform = "platform"

	// 传给platform所必须的url
	urlPluginParty    = "/plugin_party"    // 所有路由转发到插件中心处理
	urlPluginExternal = "/plugin_external" // 标记为external的路由
	// 传给platform所必须的header
	headerInstanceID  = "instance_id"
	headerRequestType = "request_type"

	defaultTimeoutSec = 30
)

var (
	timeout      = time.Duration(config.Int("main_system.timeout_sec", defaultTimeoutSec)) * time.Second
	platformAddr = fmt.Sprintf("http://%s:%d",
		config.String("platform.host", "127.0.0.1"),
		config.Int("platform.http_port", 9005),
	)
)

func AdditionProcessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch {
		case isPlatformRoute(c):
			if err := requestPlatform(c); err != nil {
				handleError(c, err)
				return
			}
		case isExternalRoute(c):
			route := hub.MatchRouter(common.RouterTypeExternal, c.Request.Method, c.Request.RequestURI)
			if route == nil {
				c.Next()
				return
			}
			if err := requestExternal(c, route); err != nil {
				handleError(c, err)
				return
			}
		case isPluginRoute(c):
			route := hub.MatchRouter(common.RouterTypeAddition, c.Request.Method, c.Request.RequestURI)
			if route == nil {
				c.Next()
				return
			}
			if err := requestAddition(c, route); err != nil {
				handleError(c, err)
				return
			}
		}

		return
	}
}

func ReplaceProcessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isPluginRoute(c) {
			c.Next()
			return
		}

		route := hub.MatchRouter(common.RouterTypeReplace, c.Request.Method, c.Request.RequestURI)
		if route == nil {
			c.Next()
			return
		}

		if err := requestReplace(c, route); err != nil {
			handleError(c, err)
			return
		}

		return
	}
}

func PrefixProcessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isPluginRoute(c) {
			c.Next()
			return
		}

		route := hub.MatchRouter(common.RouterTypePrefix, c.Request.Method, c.Request.RequestURI)
		if route == nil {
			c.Next()
			return
		}

		if err := requestPrefix(c, route); err != nil {
			handleError(c, err)
			return
		}

		c.Next()
	}
}

func SuffixProcessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isPluginRoute(c) {
			c.Next()
			return
		}

		w := &toolBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
			status:         Origin,
		}
		c.Writer = w // 劫持writer, 将handleFunc的数据先暂时写到w.body

		c.Next()

		w.status = Replace
		originBytes := w.body
		w.body = &bytes.Buffer{} // clear Origin Buffer

		code := c.Writer.Status()
		if code != http.StatusOK {
			log.Trace("resp.StatusCode: %d", code)
			// 状态不为200,就不走插件了,直接返回
			// NOTE: 注意要把原来的数据写回去
			if _, err := w.Write(originBytes.Bytes()); err != nil {
				log.ErrorDetails(err)
			}
			return
		}

		route := hub.MatchRouter(common.RouterTypeSuffix, c.Request.Method, c.Request.RequestURI)
		if route == nil {
			handleError(c, nil)
			return
		}

		err := requestSuffix(c, w, originBytes, route)
		if err != nil {
			handleError(c, err)
			return
		}
	}
}

func requestPlatform(c *gin.Context) error {
	log.Trace("prefix platform: %s", c.Request.RequestURI)
	body, err := c.GetRawData()
	if err != nil {
		return errors.Trace(err)
	}
	url := platformAddr + c.Request.RequestURI
	resp, err := request(url, c.Request.Method, c.Request.Header, body, "", "")
	if err != nil {
		return errors.Trace(err)
	}
	for key, val := range resp.Header {
		for _, v := range val {
			c.Writer.Header().Add(key, v)
		}
	}

	if err = pipe(resp, c.Writer); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func requestAddition(c *gin.Context, route *http_router.RouterInfo) error {
	log.Trace("prefix addition: %s", c.Request.RequestURI)

	body, err := c.GetRawData()
	if err != nil {
		return errors.Trace(err)
	}

	url := platformAddr + urlPluginParty + c.Request.RequestURI
	resp, err := request(url, route.Method, c.Request.Header, body, route.InstanceUUID, route.Type)
	if err != nil {
		return errors.Trace(err)
	}
	for key, val := range resp.Header {
		for _, v := range val {
			c.Writer.Header().Add(key, v)
		}
	}

	if err = pipe(resp, c.Writer); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func requestSuffix(c *gin.Context, w *toolBodyWriter, origin *bytes.Buffer, route *http_router.RouterInfo) error {
	log.Trace("suffix router: %s", c.Request.RequestURI)

	body := origin.Bytes()
	header := c.Writer.Header()
	url := platformAddr + urlPluginParty + c.Request.RequestURI

	resp, err := request(url, route.Method, header, body, route.InstanceUUID, route.Type)
	if err != nil {
		return errors.Trace(err)
	}

	for key, val := range resp.Header {
		for _, v := range val {
			c.Writer.Header().Add(key, v)
		}
	}

	if err := pipe(resp, w); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func requestPrefix(c *gin.Context, route *http_router.RouterInfo) error {
	log.Trace("prefix router: %s", c.Request.RequestURI)

	body, err := c.GetRawData()
	if err != nil {
		return errors.Trace(err)
	}

	url := platformAddr + urlPluginParty + c.Request.RequestURI

	resp, err := request(url, route.Method, c.Request.Header, body, route.InstanceUUID, route.Type)
	if err != nil {
		return errors.Trace(err)
	}

	c.Request.Header = resp.Header.Clone()
	c.Request.Body = resp.Body // 插件的responseBody作为主程序接口的requestBody
	return nil
}

func requestReplace(c *gin.Context, route *http_router.RouterInfo) error {
	log.Trace("replace router: %s", c.Request.RequestURI)

	body, err := c.GetRawData()
	if err != nil {
		return errors.Trace(err)
	}

	url := platformAddr + urlPluginParty + c.Request.RequestURI
	resp, err := request(url, route.Method, c.Request.Header, body, route.InstanceUUID, route.Type)
	if err != nil {
		return errors.Trace(err)
	}

	for key, val := range resp.Header {
		for _, v := range val {
			c.Writer.Header().Add(key, v)
		}
	}

	if err := pipe(resp, c.Writer); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func requestExternal(c *gin.Context, route *http_router.RouterInfo) error {
	log.Trace("external router: %s", c.Request.RequestURI)

	body, err := c.GetRawData()
	if err != nil {
		return errors.Trace(err)
	}
	url := platformAddr + urlPluginExternal + c.Request.RequestURI
	resp, err := request(url, route.Method, c.Request.Header, body, route.InstanceUUID, route.Type)
	if err != nil {
		return errors.Trace(err)
	}
	for key, val := range resp.Header {
		for _, v := range val {
			c.Writer.Header().Add(key, v)
		}
	}
	if err := pipe(resp, c.Writer); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func isPluginRoute(c *gin.Context) bool {
	return c.GetHeader(routeKey) == routeValPlugin
}

func isPlatformRoute(c *gin.Context) bool {
	return c.GetHeader(routeKey) == routeValPlatform
}

func isExternalRoute(c *gin.Context) bool {
	return c.GetHeader(routeKey) == routeValExternal
}

func request(url string, method string, headers http.Header, body []byte,
	instance string, routeType common.RouterType,
) (*http.Response, error) {
	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Trace(err)
	}

	for key, val := range headers {
		for _, v := range val {
			r.Header.Add(key, v)
		}
	}
	r.Header.Add(headerInstanceID, instance)
	r.Header.Add(headerRequestType, routeType)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resp, nil
}

func pipe(resp *http.Response, writer gin.ResponseWriter) error {
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	if _, err = writer.Write(result); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func handleError(c *gin.Context, err error) {
	if err != nil {
		log.ErrorDetails(err)
	}
	err = errors.PluginCallError(errors.CallPluginFailure, "")
	controllers.RenderJSONAndStop(c, err, nil)
}

type toolBodyWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status byte
}

const (
	Origin  byte = 0x0
	Replace      = 0x1
)

func (r toolBodyWriter) Write(b []byte) (int, error) {
	if r.status == Replace {
		r.body.Write(b)
		return r.ResponseWriter.Write(b)
	} else {
		return r.body.Write(b) //r.ResponseWriter.Write(b)
	}
}
