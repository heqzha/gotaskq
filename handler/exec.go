package handler

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func ExeSync(outputer func([]byte, interface{}), outputerArgs interface{}, name string, arg ...string) error {
	fmt.Printf("Exec: %s %v %s\n", name, strings.Join(arg, " "), time.Now())
	defer func() {
		fmt.Printf("Finish: %s %v %s\n\n", name, strings.Join(arg, " "), time.Now())
	}()
	cmd := exec.Command(name, arg...)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	multi := io.MultiReader(stdOut, stdErr)
	sc := bufio.NewScanner(multi)

	if err := cmd.Start(); err != nil {
		return err
	}
	// read command's stdout line by line
	go func() {
		lines := []byte{}
		for sc.Scan() {
			l := append(sc.Bytes(), '\n')
			lines = append(lines, l...)
		}
		outputer(lines, outputerArgs)
	}()
	if err := cmd.Wait(); err != nil {
		panic(err.Error())
	}
	return nil
}
