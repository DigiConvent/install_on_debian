package utils

import (
	"errors"
	"os/exec"
	"strings"
)

func Execute(c string) (string, error) {
	segments := strings.Split(c, " ")
	cmdSegments := []string{"systemctl"}
	cmdSegments = append(cmdSegments, segments...)
	s, e := exec.Command("sudo", cmdSegments...).CombinedOutput()
	if e != nil {
		return "", errors.New(e.Error() + ": " + string(s))
	}
	return string(s), e
}
