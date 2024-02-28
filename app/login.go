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

	resp, err := RequestChallenge(config.Username)
	if err != nil {
		log.Println("Request challenge failed, error:", err)
		return false
	}
	challenge, ok := resp["challenge"].(string)
	if !ok {
		log.Println("Error: parse challenge failed, resp:", resp)
		return false
	}
	ip, ok := resp["online_ip"].(string)
	if !ok {
		log.Println("Error: parse ip failed, resp:", resp)
		return false
	}

	resp, err = RequestLogin(config.Username, config.Password, ip, challenge)
	if err != nil {
		log.Println("Request login failed, error:", err)
		return false
	}
	result, ok := resp["error"].(string)
	if !ok {
		log.Println("Error: parse login result failed, resp:", resp)
		return false
	}

	if result == "ok" {
		log.Println("Login success, ip:", ip)
		return true
	} else {
		log.Println("Login failed, result:", result)
		return false
	}
}

func testConnection() bool {
	// If you were offline, when you visit any website using http (not https),
	// you will be redirected to the login page and the status code will be 200
	// If you visit any website using https, you will get a timeout error
	// So we use http generate_204 to test connection
	url := "http://wifi.vivo.com.cn/generate_204"

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true

	client := &http.Client{Transport: transport}
	resp, err := client.Get(url)

	return err == nil && resp.StatusCode == 204
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
		ChallengeURL = Host + "/cgi-bin/get_challenge"
		UserInfoURL = Host + "/cgi-bin/rad_user_info"
		LoginURL = Host + "/cgi-bin/srun_portal"
	} else {
		log.Println("Host is empty, use default host:", Host)
	}

	if config.Username == "" || config.Password == "" {
		log.Println("Error: username or password is empty")
		return Config{}, err
	}

	return config, nil
}
