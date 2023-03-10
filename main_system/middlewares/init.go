package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/main_system/platform"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/controllers"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	usePluginFlag = "plugin"

	// 所有路由转发到插件中心处理
	pluginPrefix = "/plugin_party"

	// 标记为external的路由
	pluginExternal = "/plugin_external"
)

const (
	headerInstanceID  = "instance_id"
	headerRequestType = "request_type"
)

const (
	defaultTimeoutSec = 30
)

var (
	timeout = time.Duration(config.Int("main_system.timeout_sec", defaultTimeoutSec)) * time.Second
	addr    = fmt.Sprintf("http://%s:%d",
		config.String("platform.host", "127.0.0.1"),
		config.Int("platform.http_port", 9005),
	)
)

func AdditionProcessor() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("AdditionProcessor")
		context.Next()
	}
}

func ReplaceProcessor() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("ReplaceProcessor")
		context.Next()
	}
}

func PrefixProcessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !(c.GetHeader("role") == usePluginFlag) {
			c.Next()
			return
		}

		route := platform.MatchRouter(common.RouterTypePrefix, c.Request.Method, c.Request.RequestURI)
		if route == nil {
			c.Next()
			return
		}

		if err := requestPrefix(c, route); err != nil {
			log.ErrorDetails(err)
			err = errors.PluginCallError(errors.CallPluginFailure, "")
			controllers.RenderJSONAndStop(c, err, nil)
			return
		}
		c.Next()
	}
}

func SuffixProcessor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果没有标志,说明是普通请求,不需要处理,直接next
		if !(c.GetHeader("role") == usePluginFlag) {
			c.Next()
			return
		}

		w := &toolBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
			status:         Origin,
		}
		c.Writer = w

		c.Next()

		w.status = Replace
		originBytes := w.body
		w.body = &bytes.Buffer{} // clear Origin Buffer

		if code := c.Writer.Status(); code != http.StatusOK {
			log.Trace("resp.StatusCode: %d", code)
			whenErr(c, w, originBytes)
			return
		}

		route := platform.MatchRouter(common.RouterTypeSuffix, c.Request.Method, c.Request.RequestURI)
		if route == nil {
			whenErr(c, w, originBytes)
			return
		}

		err := requestSuffix(c, w, originBytes, route)
		if err != nil {
			log.ErrorDetails(err)
			whenErr(c, w, originBytes)
			return
		}
	}
}

func requestSuffix(c *gin.Context, w *toolBodyWriter, origin *bytes.Buffer, route *http_router.RouterInfo) error {
	log.Trace("suffix router: %s", c.Request.RequestURI)

	body := origin.Bytes()
	resp, err := request(c, c.Writer.Header(), body, route.InstanceUUID, common.RouterTypeSuffix)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)
	}

	for key, val := range c.Request.Header {
		for _, v := range val {
			c.Writer.Header().Add(key, v)
		}
	}

	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	if _, err = w.Write(result); err != nil {
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
	resp, err := request(c, nil, body, route.InstanceUUID, common.RouterTypePrefix)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)
	}

	for key, val := range resp.Header {
		for _, v := range val {
			c.Request.Header.Add(key, v)
		}
	}

	c.Request.Body = resp.Body // 插件的responseBody作为主程序接口的requestBody
	return nil
}

func whenErr(c *gin.Context, w *toolBodyWriter, origin *bytes.Buffer) {
	if _, err := w.Write(origin.Bytes()); err != nil {
		log.ErrorDetails(err)
	}
	c.Next()
}

func request(
	c *gin.Context,
	headers http.Header,
	body []byte,
	instance string,
	routeType common.RouterType,
) (*http.Response, error) {
	url := addr + pluginPrefix + c.Request.RequestURI
	r, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
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
