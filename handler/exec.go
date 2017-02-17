package handler

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func ExeSync(outputer func(*bytes.Buffer, interface{}), outputerArgs interface{}, name string, arg ...string) error {
	fmt.Printf("Exec: %s %v\n", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	go func() {
		_, err := buf.ReadFrom(out)
		if err != nil {
			panic(err.Error())
		}
		outputer(buf, outputerArgs)
	}()
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
