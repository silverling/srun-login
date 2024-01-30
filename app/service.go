package main

import (
	_ "embed"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

const AUTOSTART_REGISTRY_PATH = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

var lockFile *os.File

//go:embed assets/favicon.ico
var iconData []byte

func enableAutoStart() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, AUTOSTART_REGISTRY_PATH, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("Srun Login", GetProgramPath())
	if err != nil {
		return err
	}
	return nil
}

func disableAutoStart() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, AUTOSTART_REGISTRY_PATH, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.DeleteValue("Srun Login")
	if err != nil {
		return err
	}
	return nil
}

func checkAutoStart() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, AUTOSTART_REGISTRY_PATH, registry.ALL_ACCESS)
	if err != nil {
		log.Output(2, err.Error())
	}
	defer key.Close()
	value, valueType, err := key.GetStringValue("Srun Login")
	if err != nil {
		return false
	}
	if value != GetProgramPath() || valueType != registry.SZ {
		return false
	}
	return true
}

func singleInstance() {
	// check lock file in temp folder
	// if lock file exist, exit
	lockFilePath := filepath.Join(os.TempDir(), "srun-login.lock")
	_, err := os.Stat(lockFilePath)
	if err == nil {
		log.Output(2, "lock file exist, trying to delete")
		err = os.Remove(lockFilePath) // you can not delete a file which is in use
		if err != nil {
			log.Output(2, err.Error())
			os.Exit(0)
		}
	}

	// if lock file not exist, create lock file
	lockFile, err = os.Create(lockFilePath)
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(0)
	}
}

func cleanSingleInstance() {
	// when exit, delete lock file
	lockFile.Close()
	lockFilePath := filepath.Join(os.TempDir(), "srun-login.lock")
	err := os.Remove(lockFilePath)
	if err != nil {
		log.Output(2, err.Error())
	}
}

var _run = func() {}

func RunService(_runFunc func()) {
	_run = _runFunc

	// setup logging
	logPath := filepath.Join(GetProgramFolder(), "log.txt")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Output(2, err.Error())
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Start")

	singleInstance()
	// if run on windows
	if runtime.GOOS == "windows" {
		// use systray
		systray.Run(onReady, onQuit)
	} else {
		// directly run
		_run()
	}
	cleanSingleInstance()
}

func onReady() {
	systray.SetTitle("Srun Login")
	systray.SetTooltip("Srun Login")
	systray.SetIcon(iconData)

	mOpenFolder := systray.AddMenuItem("Open Folder", "Open Folder")
	mAutoStart := systray.AddMenuItemCheckbox("Autostart", "Autostart", false)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	if checkAutoStart() {
		mAutoStart.Check()
	} else {
		mAutoStart.Uncheck()
	}

	// create goroutine to run
	go _run()

	for {
		select {
		case <-mOpenFolder.ClickedCh:
			OpenProgramFoler()
		case <-mAutoStart.ClickedCh:
			if mAutoStart.Checked() {
				err := disableAutoStart()
				if err != nil {
					log.Output(2, err.Error())
				} else {
					mAutoStart.Uncheck()
				}
			} else {
				err := enableAutoStart()
				if err != nil {
					log.Output(2, err.Error())
				} else {
					mAutoStart.Check()
				}
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		case <-sigc:
			systray.Quit()
			return
		}
	}
}

func onQuit() {
	log.Println("Quit")
}
