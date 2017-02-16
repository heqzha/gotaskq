package handler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ExeSync(outputer func(io.Reader, interface{}) error, outputerArgs interface{}, name string, arg ...string) error {
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
		n, err := buf.ReadFrom(out)
		if err != nil {
			panic(err.Error())
		}
		buf.WriteTo(os.Stdout)
		fmt.Printf("number of bytes: %d\n", n)
	}()
	// if err := outputer(out, outputerArgs); err != nil {
	// 	return err
	// }
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
