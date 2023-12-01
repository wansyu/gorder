package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	// "time"
)

func main() {
	r := gin.Default()
	// 定义路由
	r.GET("/", func(c *gin.Context) {
		// time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	// 定义一个接口用于触发本地程序的执行
	r.POST("/run-local-program", func(c *gin.Context) {
		// 执行本地程序的命令
		cmd := exec.Command("curl", "ip.sb")
		output, err := cmd.CombinedOutput()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回本地程序的输出结果
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSON(http.StatusOK, gin.H{"output": string(output)})
	})

	// 启动 HTTP 服务器
	r.Run(":8080")
}
