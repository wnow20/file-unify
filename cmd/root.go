package cmd

import (
	"fmt"
	"os"
)

var rootCmd = unifyCmd

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
