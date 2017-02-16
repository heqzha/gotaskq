package handler

import (
	"fmt"
	"io"
	"os/exec"
)

func ExeSync(outputer func(io.Reader, interface{}) error, outputerArgs interface{}, name string, arg ...string) error {
	fmt.Printf("Exec: %s %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := outputer(out, outputerArgs); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
