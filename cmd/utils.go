package cmd

import (
	"fmt"
	"os"
)

func osExitErr(errMessage string) {
	fmt.Fprintln(os.Stderr, errMessage)
	os.Exit(1)
}

func Execute() error {
	return rootCmd.Execute()
}
