package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/utils/errors"
	"net/http"
)

func RenderJSON(c *gin.Context, err error, obj interface{}) {
	errp := buildErrPayloadAndLog(err, true)
	if errp.HttpStatus == http.StatusOK {
		if obj == nil {
			c.JSON(errp.HttpStatus, errp)
		} else {
			c.JSON(errp.HttpStatus, obj)
		}
	} else {
		c.JSON(errp.HttpStatus, errp)
	}
	c.Next()
}

func RenderError(c *gin.Context, result error) {
	errp := buildErrPayloadAndLog(result, true)
	c.JSON(errp.HttpStatus, errp)
}

func RenderJSONAndStop(c *gin.Context, result error, obj interface{}) {
	errp := buildErrPayloadAndLog(result, true)
	if errp.HttpStatus == http.StatusOK {
		if obj == nil {
			c.JSON(errp.HttpStatus, errp)
		} else {
			c.JSON(errp.HttpStatus, obj)
		}
	} else {
		c.JSON(errp.HttpStatus, errp)
	}
}

func buildErrPayloadAndLog(err error, shouldLog bool) (errp *errors.ErrPayload) {
	errp = errors.NewErrPayload(err)
	if shouldLog {
		// 根据状态码打印日志
		if errp.HttpStatus < 400 {
			// 不需要打印日志
		} else if errp.HttpStatus >= 500 && errp.HttpStatus < 600 {
			// 服务端错误

			log.ErrorDetails(err)
			// 对客户端隐藏详细信息
			errp.Code = errors.ServerError
			errp.HttpStatus = http.StatusInternalServerError
			errp.Desc = ""
			errp.Values = nil
		} else {
			// 客户端错误 & 自定义错误，Warn
			log.WarnDetails(err)
		}
	}
	return
}
