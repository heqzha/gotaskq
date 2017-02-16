package process

import (
	"fmt"
	"time"

	"github.com/heqzha/gotaskq/conf"
	"github.com/heqzha/gotaskq/handler"
	ccc "github.com/heqzha/goutils/concurrency"
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
	l, err := fh.NewLine(switcher, running, sleep)
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

func running(c *flow.Context) {
	defer c.Next()
	w := c.MustGet("workers").(*ccc.WorkersPool)
	t := c.MustGet("task").(*conf.TaskT)
	f := c.MustGet("func").(func(interface{}) interface{})
	output := make(chan interface{}, 16)
	w.CollectWithOutput(f, map[string]interface{}{
		"name": t.Name,
		"args": t.Args,
	}, 0, output)

	select {
	case o := <-output:
		m := o.(map[string]interface{})
		if m["err"] != nil {
			err := m["err"].(error)
			logger.Error("RunTaskExecutor.running", err.Error())
		} else {
			switch t.Output {
			case "stdout":
				fmt.Println(string(m["output"].([]byte)))
			}
		}
	}
}

func cmd(params interface{}) interface{} {
	mp := params.(map[string]interface{})
	output, err := handler.ExeSync(mp["name"].(string), mp["args"].([]string)...)
	return map[string]interface{}{
		"output": output,
		"err":    err,
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
