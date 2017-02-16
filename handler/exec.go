package handler

import (
	"os/exec"
)

func ExeSync(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
