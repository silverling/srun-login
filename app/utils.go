package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func CurrentMilisecond() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func GetClientIP() (string, error) {
	res, err := RequestUserInfo()
	if err != nil {
		return "", err
	}
	ip := res["online_ip"].(string)
	return ip, nil
}

func GetProgramPath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return exePath
}

func GetProgramFolder() string {
	return filepath.Dir(GetProgramPath())
}

func OpenProgramFoler() {
	err := exec.Command("explorer", GetProgramFolder()).Start()
	if err != nil {
		log.Output(2, err.Error())
	}
}
