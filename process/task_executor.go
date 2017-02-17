package process

import (
	"fmt"
	"os"
	"time"

	"github.com/heqzha/gotaskq/conf"
	"github.com/heqzha/gotaskq/handler"
	ccc "github.com/heqzha/goutils/concurrency"
	gufile "github.com/heqzha/goutils/file"
	"github.com/heqzha/goutils/flow"
	"github.com/heqzha/goutils/logger"
)

var (
	fh = flow.FlowNewHandler()
)

func Serve() (err error) {
	var stopped bool
	for stopped, err = fh.AreAllStopped(); !stopped && err == nil; {
		ticker := time.NewTicker(time.Second)
	slp:
		select {
		case <-ticker.C:
			break slp
		}
	}
	return
}

func StopAll() error {
	fh.AreAllStopped()
	return fh.Destory()
}

func RunTaskExecutor(workers *ccc.WorkersPool, tasks []*conf.TaskT) error {
	l, err := fh.NewLine(switcher, outputer, running, sleep)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fh.Start(l, flow.Params{
			"workers": workers,
			"task":    task,
		})
	}
	return nil
}

func switcher(c *flow.Context) {
	t := c.MustGet("task").(*conf.TaskT)
	switch t.Type {
	case "cmd":
		c.Set("func", cmd)
		c.Next()
		return
	default:
		logger.Error("RunTaskExecutor.switcher", fmt.Sprintf("Unknown task type: %s", t.Type))
		c.Abort()
		return
	}
}

func outputer(c *flow.Context) {
	t := c.MustGet("task").(*conf.TaskT)
	switch t.Output {
	case "stdout":
		c.Set("outputer", stdout)
		c.Next()
		return
	case "file":
		c.Set("outputer", file)
		c.Set("outputerArgs", t.OutputFile)
		c.Next()
		return
	}
}

func running(c *flow.Context) {
	defer c.Next()
	w := c.MustGet("workers").(*ccc.WorkersPool)
	t := c.MustGet("task").(*conf.TaskT)
	f := c.MustGet("func").(func(interface{}) interface{})
	o := c.MustGet("outputer").(func([]byte, interface{}))
	oa, _ := c.Get("outputerArgs")
	w.Collect(f, map[string]interface{}{
		"outputer":     o,
		"outputerArgs": oa,
		"name":         t.Name,
		"args":         t.Args,
	}, 0)
}

func cmd(params interface{}) interface{} {
	mp := params.(map[string]interface{})
	if err := handler.ExeSync(mp["outputer"].(func([]byte, interface{})), mp["outputerArgs"], mp["name"].(string), mp["args"].([]string)...); err != nil {
		logger.Error("RunTaskExecutor.cmd", err.Error())
	}
	return nil
}

func stdout(buf []byte, outputerArgs interface{}) {
	// buf.WriteTo(os.Stdout)
	fmt.Print(string(buf))
}

func file(buf []byte, outputerArgs interface{}) {
	fileName := outputerArgs.(string)
	if !gufile.Exists(fileName) {
		path, err := gufile.GetPath(fileName)
		if err != nil {
			logger.Error("RunTaskExecutor.file", err.Error())
			return
		}
		if err := gufile.MkPath(path, 0777); err != nil {
			logger.Error("RunTaskExecutor.file", err.Error())
			return
		}
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		logger.Error("RunTaskExecutor.file", err.Error())
		return
	}
	defer f.Close()
	if _, err := f.Write(buf); err != nil {
		logger.Error("RunTaskExecutor.file", err.Error())
		return
	}
}

func sleep(c *flow.Context) {
	defer c.Repeat()
	t := c.MustGet("task").(*conf.TaskT)
	ticker := time.NewTicker(time.Duration(t.Duration) * time.Second)

	select {
	case <-ticker.C:
		return
	}
}
