package main

import (
	"flag"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"
	// "os"
	"os/exec"
	"log"
	"github.com/gin-gonic/gin"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

type ProgramConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	Key           string          `json:"key"`
	Salt          string          `json:"salt"`
	ProgramPaths  []ProgramConfig `json:"program_paths"`
}
func EncryptPassword(password string, salt []byte) (string, error) {
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}

func main() {
	// 解析命令行参数
	configPath := flag.String("c", "config.json", "Path to config file")
	ip := flag.String("ip", "", "IP address")
	port := flag.String("p", "8080", "Port number")
	route := flag.String("r", "/call-program", "Route path")
	flag.Parse()

	// 检查是否提供了配置文件路径
	if *configPath == "" {
		log.Fatal("Config file path is required")
		return
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 定义路由
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gorder Server")
	})
	router.POST(*route, func(c *gin.Context) {
		// 读取config.json文件
		configData, err := ioutil.ReadFile(*configPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read config file"})
			return
		}

		// 解析config.json文件
		var config Config
		err = json.Unmarshal(configData, &config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse config file"})
			return
		}

		// 解析请求中的JSON数据
		var requestData struct {
			ProgramName string `json:"program_name"`
			Key         string `json:"key"`
		}
		err = c.BindJSON(&requestData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		// 验证密钥
		encryptedDk, err := EncryptPassword(requestData.Key, []byte(config.Salt))
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		if config.Key != encryptedDk {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid key"})
			return
		}

		// 查找程序路径
		var programPath string
		for _, program := range config.ProgramPaths {
			if program.Name == requestData.ProgramName {
				programPath = program.Path
				break
			}
		}

		if programPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program name"})
			return
		}

		// 执行本地程序
		cmd := exec.Command(programPath)
		// err = cmd.Run()
		output, err := cmd.Output()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute program","output": string(output)})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "Program called successfully ","output": string(output)})
	})

	router.Run(*ip+":"+*port)
}
