package utils

import (
	"fmt"
	"os"
)

// ExitOnFailure prints a fatal error message and exits the process with status 1.
func ExitOnFailure(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "[CRIT] %s. ", err.Error())
	os.Exit(1)
}
