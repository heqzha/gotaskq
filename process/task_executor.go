package process

import (
	"time"

	"github.com/heqzha/gotaskq/conf"
	ccc "github.com/heqzha/goutils/concurrency"
	"github.com/heqzha/goutils/flow"
)

var (
	fh = flow.FlowNewHandler()
)

func StopAll() error {
	return fh.Destory()
}

func sleepAndRestart(c *flow.Context) {
	defer c.Repeat()
	interval := c.MustGet("duration").(time.Duration)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			return
		}
	}
}

func RunTaskExecutor(workers *ccc.WorkersPool, tasks *conf.TaskT) {

}
