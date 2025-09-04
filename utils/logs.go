package utils

import (
	"fmt"
	"strings"
)

func Logs(name string, numLines int) {
	output, err := Execute("journalctl -u " + "main_test --reverse")
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(output, "\n")
	for i := range numLines {
		fmt.Println(lines[i])
	}
}
