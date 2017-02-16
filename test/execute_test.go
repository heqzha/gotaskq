package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/heqzha/gotaskq/handler"
)

func TestExecCmd(t *testing.T) {
	output, err := handler.ExeSync("echo", "hello", "world")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}
