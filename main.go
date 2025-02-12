package main

import (
	"fmt"
	"os"

	cli "gitcode.com/hammerklavier/openvpn-config-generator-go/CLI"
)

func main() {
	rootArgs, subArgs, err := cli.ParseCli()
	if err != nil {
		fmt.Printf("Failed to parse arguments: %v\n", err)
		os.Exit(1)
	}
	if rootArgs.Verbose {
		fmt.Printf("Root args: %v; Sub args: %v\n", rootArgs, subArgs)
	}

	switch subArgs.(type) {
	case cli.InitArguments:
		if rootArgs.Verbose {
			fmt.Println("Task: Initialise configuration, creates server configuration")
		}
	case cli.ClientArguments:
		if rootArgs.Verbose {
			fmt.Println("Task: Create client configurations")
		}
	case nil:
		if rootArgs.Verbose {
			fmt.Println("Exit")
		}
		os.Exit(0)
	default:
		fmt.Println("Error: Failed to infer type of subArgs!")
		fmt.Printf("%T\n", subArgs)
	}
}
