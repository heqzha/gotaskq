package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/heqzha/gotaskq/conf"
	"github.com/heqzha/gotaskq/process"
	ccc "github.com/heqzha/goutils/concurrency"
	"github.com/heqzha/goutils/logger"
)

func init() {
	logger.Config(conf.Config.LogDir, logger.LOG_LEVEL_DEBUG)
}

func createPID(name string) int {
	wd, _ := os.Getwd()
	pidFile, err := os.OpenFile(filepath.Join(wd, fmt.Sprintf("%s.pid", name)), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal("main", fmt.Sprintf("failed to create pid file: %s", err.Error()))
		os.Exit(1)
	}
	defer pidFile.Close()

	pid := os.Getpid()
	pidFile.WriteString(strconv.Itoa(pid))
	return pid
}

func main() {
	pid := createPID("gotaskq")
	workers := &ccc.WorkersPool{}
	workers.Start(conf.Config.TaskQueue.MaxWorkers, conf.Config.TaskQueue.MaxBuffer)
	defer workers.Stop()

	if len(conf.Config.Tasks) == 0 {
		fmt.Println("No task to run")
		return
	}
	if err := process.RunTaskExecutor(workers, conf.Config.Tasks); err != nil {
		e := fmt.Sprintf("failed to initialize task executor: %s", err.Error())
		logger.Fatal("main", e)
		fmt.Println(e)
		return
	}
	defer process.StopAll()

	fmt.Printf("Start to run executor(PID:%d)\n", pid)
	process.Serve()
}
