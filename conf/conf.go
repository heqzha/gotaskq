package conf

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

var (
	reload       = flag.Bool("reload", false, "reload process")
	configFile   = flag.String("config", "__unset__", "service config file")
	maxThreadNum = flag.Int("max-thread", 0, "max threads of service")
	debugMode    = flag.Bool("debug", false, "debug mode")

	Config    = &ConfigT{}
	DebugMode bool
)

type TaskQueueT struct {
	MaxWorkers int `yaml:"max_workers"`
	MaxBuffer  int `yaml:"max_buffer"`
}

type TaskT struct {
	Type     string   `yaml:"type"`
	Name     string   `yaml:"name"`
	Args     []string `yaml:"args"`
	Output   string   `yaml:"output"`
	Duration int64    `yaml:"duration"`
}

type ConfigT struct {
	LogDir    string      `yaml:"log_dir"`
	TaskQueue *TaskQueueT `yaml:"task_queue"`
	Tasks     []*TaskT    `yaml:"tasks"`
}

func init() {
	flag.Parse()
	DebugMode = *debugMode

	if *reload {
		wd, _ := os.Getwd()
		pidFile, err := os.Open(filepath.Join(wd, "gotaskq.pid"))
		if err != nil {
			log.Printf("Failed to open pid file: %s", err.Error())
			os.Exit(1)
		}
		pids := make([]byte, 10)
		n, err := pidFile.Read(pids)
		if err != nil {
			log.Printf("Failed to read pid file: %s", err.Error())
			os.Exit(1)
		}
		if n == 0 {
			log.Printf("No pid in pid file: %s", err.Error())
			os.Exit(1)
		}
		_, err = exec.Command("kill", "-USR2", string(pids[:n])).Output()
		if err != nil {
			log.Printf("Failed to restart service: %s", err.Error())
			os.Exit(1)
		}
		pidFile.Close()
		os.Exit(0)
	}
	if *maxThreadNum == 0 {
		*maxThreadNum = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*maxThreadNum)

	if *configFile == "__unset__" {
		p, _ := os.Getwd()
		*configFile = filepath.Join(p, "config.yml")
	}

	confFile, err := filepath.Abs(*configFile)
	if err != nil {
		log.Printf("No correct config file: %s - %s", *configFile, err.Error())
		os.Exit(1)
	}

	confBs, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Printf("Failed to read config fliel <%s> : %s", confFile, err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal(confBs, Config)
	if err != nil {
		log.Printf("Failed to parse config fliel <%s> : %s", confFile, err.Error())
		os.Exit(1)
	}
}
