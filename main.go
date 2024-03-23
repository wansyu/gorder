package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
)

type ProgramConfig struct {
	Name string   `json:"name"`
	Path string   `json:"path"`
	Args []string `json:"args"`
}

type Config struct {
	Key          string          `json:"key"`
	Salt         string          `json:"salt"`
	ProgramPaths []ProgramConfig `json:"program_paths"`
}

// this function is used to initialize the config in the first run
func InitializeConfig(configPath string) {
	var ramConfig Config
	//generate a random ASCII chars for ramConfig.Salt
	ramConfig.Salt = generateRandomASCIIChars(32)
	//generate a random ASCII chars for ramConfig.Key
	ramkey := generateRandomASCIIChars(32)
	log.Printf("apikey is: %v", ramkey)
	ramConfig.Key, _ = EncryptPassword(ramkey, []byte(ramConfig.Salt))
	//generate a defaut Program for ramConfig.ProgramPaths
	ramConfig.ProgramPaths = []ProgramConfig{
		{
			Name: "getip0",
			Path: "curl",
			Args: []string{"ip.sb"},
		},
		{
			Name: "getip1",
			Path: "curl",
			Args: []string{"-x", "socks5://127.0.0.1:1080", "ip.sb"},
		},
	}
	// write the ramConfig to the file configPath
	saveConfig(ramConfig, configPath)
}
func saveConfig(config Config, configPath string) {
	// Convert config to a formatted JSON string
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatalf("Error marshalling config to JSON: %v", err)
	}

	// Write the JSON string to configPath
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON data to %s: %v", configPath, err)
	}
}
func generateRandomASCIIChars(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating random bytes: %v", err)
	}

	var result string
	for _, v := range b {
		// Convert each byte to its corresponding ASCII character
		vv := v%94 + 33 // ASCII printable characters range from 33 to 126
		if vv == 32 || vv == 34 || vv == 39 || vv == 92 {
			vv += 1
		}
		result += string(vv)
	}

	return result
}

// this function is used to encrypt the password
func EncryptPassword(password string, salt []byte) (string, error) {
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}

// this function is used to load the configuration file
func loadConfig(configPath string) (Config, error) {
	// check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found: %s", configPath)
		InitializeConfig(configPath)
	}
	// read the config file
	file, err := os.ReadFile(configPath)
	var ramConfig Config
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return ramConfig, err
	}

	// parse the config file
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
		return ramConfig, err
	}

	return config, nil
}
func main() {
	// parse command line arguments
	configPath := flag.String("c", "config.json", "Path to config file")
	ip := flag.String("ip", "", "IP address")
	port := flag.String("p", "8080", "Port number")
	route := flag.String("r", "/call-program", "Route path")
	flag.Parse()

	// check if config file path is provided
	if *configPath == "" {
		log.Fatal("Config file path is not provided, use default path ./config.json")
		*configPath = "./config.json"
		// return
	}
	// load the config file
	config, cerr := loadConfig(*configPath)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// define a simple route
	router.GET("/", func(c *gin.Context) {
		if cerr != nil {
			c.String(http.StatusOK, "Error Gorder Server")
			return
		}
		c.String(http.StatusOK, "Welcome Gorder Server")
	})
	if cerr != nil {
		log.Fatalf("error: Loading JSON data: %v", cerr)
		return
	}
	router.POST(*route, func(c *gin.Context) {
		// parse the request data
		var requestData struct {
			ProgramName string `json:"program_name"`
			Key         string `json:"key"`
		}
		err := c.BindJSON(&requestData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		// authenticate the request
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

		// check if the program name is valid
		var programPath string
		var programArgs []string
		for _, program := range config.ProgramPaths {
			if program.Name == requestData.ProgramName {
				programPath = program.Path
				programArgs = program.Args
				break
			}
		}

		if programPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program name"})
			return
		}

		// execute the program
		cmd := exec.Command(programPath, programArgs...)
		// err = cmd.Run()
		output, err := cmd.Output()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute program", "output": string(output)})
			return
		}

		// return the output
		c.JSON(http.StatusOK, gin.H{"message": "Program called successfully ", "output": string(output)})
	})

	router.Run(*ip + ":" + *port)
}
