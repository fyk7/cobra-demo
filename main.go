package main

import (
	"fmt"
	"os"

	"github.com/fyk7/cobra-demo/cmd"
)

func main() {
	// cobraエントリーポイント
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
