package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router"
	"github.com/obgnail/plugin-platform/platform/controllers"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"strconv"
	"strings"
)

const (
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
	defaultContentType = "application/json"
)

func PluginInvoke() gin.HandlerFunc {
	return invoke
}

func invoke(c *gin.Context) {
	uri := c.Request.RequestURI
	instanceUUID := c.GetHeader(headerInstanceID)
	requestType := strings.ToLower(c.GetHeader(headerRequestType))

	// NOTE: 必须先判断pluginExternal,因为header的优先级高于url
	if requestType == common.RouterTypeExternal || strings.Contains(uri, pluginExternal) {
		req, err := convert2Request(c, uri, c.Request.Method)
		if err != nil {
			handlerError(c, err.Error())
			return
		}
		resp := <-handler.CallPluginExternalHTTP(instanceUUID, req)
		convertResponse(c, resp)
	} else if strings.Contains(uri, pluginPrefix) {
		url, err := getUrl(uri, pluginPrefix)
		if err != nil {
			handlerError(c, err.Error())
			return
		}

		routerInfo := router.MatchRouter(requestType, c.Request.Method, url)
		if routerInfo == nil {
			log.Trace("dismatch: %s", url)
			return
		}

		req, err := convert2Request(c, url, routerInfo.Method)
		if err != nil {
			handlerError(c, err.Error())
			return
		}
		resp := <-handler.CallPluginInternalHTTP(instanceUUID, req, routerInfo.FunctionName)
		convertResponse(c, resp)
	}
}

func getUrl(uri string, splitString string) (string, error) {
	parts := strings.Split(uri, splitString)
	if len(parts) < 2 {
		return "", fmt.Errorf("uri error: %s", uri)
	}
	url := parts[1]
	return url, nil
}

func convert2Request(c *gin.Context, url, method string) (*common_type.HttpRequest, error) {
	body, err := c.GetRawData()
	if err != nil {
		return nil, errors.Trace(err)
	}

	req := &common_type.HttpRequest{
		Method:  method,
		Url:     url,
		Body:    body,
		Headers: make(map[string][]string),
	}
	for k, v := range c.Request.Header {
		req.Headers[k] = v
	}

	return req, nil
}

func convertResponse(c *gin.Context, resp *common_type.HttpResponse) {
	if resp == nil {
		return
	}
	if resp.Err != nil {
		handlerError(c, resp.Err.Msg())
		return
	}

	log.Trace("resp: %+v", resp)

	contentType := defaultContentType
	for k, val := range resp.Headers {
		if strings.ToUpper(k) == "Content-Type" {
			contentType = k
		}
		for _, v := range val {
			c.Writer.Header().Add(k, v)
		}
	}
	c.Writer.Header().Add("Content-Length", strconv.Itoa(len(resp.Body)))

	c.Render(resp.StatusCode, render.Data{
		ContentType: contentType,
		Data:        resp.Body,
	})
}

func handlerError(c *gin.Context, reason string) {
	err := errors.PluginCallError(errors.CallPluginFailed, reason)
	controllers.RenderJSONAndStop(c, err, nil)
}
