package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/main_system/platform/hub"
	"io/ioutil"
)

func registerRouter(app *gin.Engine) {
	app.GET("/prefix", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		c.String(200, "Hello Wold main system"+string(data))
	})

	app.GET("/suffix", func(c *gin.Context) {
		fmt.Println("this is before message")
		c.String(200, "Hello Wold main system before message")
	})
	app.GET("/suffix_error", func(c *gin.Context) {
		fmt.Println("this is before message")
		c.String(400, "Hello Wold main system before message")
	})
	app.POST("/replace", func(c *gin.Context) {
		c.String(400, "replace message")
	})
	app.GET("/ability_test", func(c *gin.Context) {
		fmt.Println("ability test")

		instanceID := "TNcoTKHS"
		args1 := "args1"
		result1, err := hub.ExecuteAbility(instanceID, "send_short_message-QWERASDF",
			"send_short_message", "getEmail", []byte(args1))
		if err != nil {
			panic(err)
		}

		args2 := "args2"
		result2, err := hub.ExecuteAbility(instanceID, "send_short_message-QWERASDF",
			"send_short_message", "sendShortMessage", []byte(args2))
		if err != nil {
			panic(err)
		}

		c.String(400, string(result1)+string(result2))
	})
}
