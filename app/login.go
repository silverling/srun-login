package main

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func login(config Config) bool {
	log.Print("Login......")

	resp, _ := RequestChallenge(config.Username)
	challenge := resp["challenge"].(string)
	ip := resp["online_ip"].(string)

	resp, _ = RequestLogin(config.Username, config.Password, ip, challenge)
	result := resp["error"].(string)

	if result == "ok" {
		log.Println("Login success, ip:", ip)
		return true
	} else {
		log.Println("Login failed, error:", result)
		return false
	}
}

func testConnection() bool {
	resp, err := http.Get("http://www.baidu.com")
	return err == nil && resp.StatusCode == 200
}

func loadConfig(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	if config.Host != "" {
		log.Println("Use given host:", config.Host)
		Host = config.Host
	} else {
		log.Println("Host is empty, use default host:", Host)
	}
	ChallengeURL = Host + "/cgi-bin/get_challenge"
	UserInfoURL = Host + "/cgi-bin/rad_user_info"
	LoginURL = Host + "/cgi-bin/srun_portal"

	if config.Username == "" || config.Password == "" {
		log.Println("Error: username or password is empty")
		return Config{}, err
	}

	return config, nil
}
