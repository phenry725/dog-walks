package cmd

import (
	"fmt"
	"os"
)

func osExitErr(errMessage string) {
	fmt.Fprint(os.Stderr, errMessage)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		osExitErr(fmt.Sprintf("Whoops. There was an error while executing your CLI '%s'", err))
	}
}
