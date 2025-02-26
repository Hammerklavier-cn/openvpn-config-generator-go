package main

import (
	"fmt"
	"os"
	"path"

	cli "gitcode.com/hammerklavier/openvpn-config-generator-go/CLI"
	cerinit "gitcode.com/hammerklavier/openvpn-config-generator-go/cer_init"
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
		initArgs := subArgs.(cli.InitArguments)
		// 0. init working directory
		err := cerinit.TargetDirInit(initArgs.Dir, rootArgs.Verbose)
		if err != nil {
			fmt.Printf("Failed to initialise target directory: %v\n", err)
			os.Exit(1)
		}
		// 1. create certificate authority
		err = cerinit.CertificateAuthorityInit(path.Join(initArgs.Dir, "CA"), initArgs.Algorithm, initArgs.Digest, rootArgs.Verbose)
		if err != nil {
			fmt.Printf("Failed to create certificate authority: %v\n", err)
			os.Exit(1)
		}
		// 2. create server certificate
		// 3. create client certificate
	case cli.ClientArguments:
		if rootArgs.Verbose {
			fmt.Println("Task: Create client configurations")
		}
	case nil:
		if rootArgs.Verbose {
			fmt.Println("Note: You did not pass a subcommand.")
			fmt.Println("Exit")
		}
		os.Exit(0)
	default:
		fmt.Println("Error: Failed to infer type of subArgs!")
		fmt.Printf("%T\n", subArgs)
	}
}
