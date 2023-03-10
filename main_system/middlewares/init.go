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
	return prefix
}

func SuffixProcessor() gin.HandlerFunc {
	return suffix
}

func suffix(c *gin.Context) {
	c.Next()

	if !(c.GetHeader("role") == usePluginFlag) {
		c.Next()
		return
	}

	route := platform.MatchRouter(common.RouterTypePrefix, c.Request.Method, c.Request.RequestURI)
	if route == nil {
		c.Next()
		return
	}
}

func prefix(c *gin.Context) {
	if !(c.GetHeader("role") == usePluginFlag) {
		c.Next()
		return
	}

	route := platform.MatchRouter(common.RouterTypePrefix, c.Request.Method, c.Request.RequestURI)
	if route == nil {
		c.Next()
		return
	}

	if err := requestForwardingPrefix(c, route); err != nil {
		log.ErrorDetails(err)
		err = errors.PluginCallError(errors.CallPluginFailure, "")
		controllers.RenderJSONAndStop(c, err, nil)
		return
	}
	c.Next()
}

func requestForwardingPrefix(c *gin.Context, route *http_router.RouterInfo) error {
	log.Trace("prefix router: %s", c.Request.RequestURI)

	url := addr + pluginPrefix + c.Request.RequestURI
	method := c.Request.Method

	body, err := c.GetRawData()
	if err != nil {
		return errors.Trace(err)
	}

	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return errors.Trace(err)
	}
	r.Header.Add(headerInstanceID, route.InstanceUUID)
	r.Header.Add(headerRequestType, common.RouterTypePrefix)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(r)
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
