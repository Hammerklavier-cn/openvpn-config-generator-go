package cli

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1-alpha1"

var activated_subcommand string

type wrongValue struct {
	message string
}

func (w wrongValue) Error() string {
	return w.message
}

var rootCmd = &cobra.Command{
	Use:   "opvn-gen",
	Short: "opvn-gen is an highly automatic openvpn configuration file (.opvn) generator",
	Long: `A highly automatic openvpn configuration file (.opvn) generator built with
            love by hammerklavier in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
	Version: VERSION,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise configuration and creates server configuration.",
	Long: `Initialise configuration and creates server configuration.
			Please execute this before ‘opvn-gen'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generate certificate for OpenVPN server.")
		activated_subcommand = "init"
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate client certificate (.opvn).",
	Long:  "Generate client certificate (.opvn).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generate certificate for clients.")
		activated_subcommand = "client"
	},
}

var AlgorithmChoices = []string{
	"RSA", "ECDSA", "EDDSA",
}

var DigestChoices = []string{
	"MD5", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512",
}

type SubcommandArguments interface {
	isSubcommandArgumentStruct() bool
}

type rootcommandArgument interface {
	isRootcommandArgumentStruct() bool
}

type RootArguments struct {
	// For containing results of root arguments
	Version bool
	Verbose bool
}

type InitArguments struct {
	Dir       string
	KeySize   int
	Algorithm string
	Digest    string
	days      int
}

type ClientArguments struct {
	Dir  string
	name string
}

func (r RootArguments) isRootcommandArgumentStruct() bool  { return true }
func (i InitArguments) isSubcommandArgumentStruct() bool   { return true }
func (c ClientArguments) isSubcommandArgumentStruct() bool { return true }

func ParseCli() (RootArguments, SubcommandArguments, error) {
	rootArgs := RootArguments{}
	initArgs := InitArguments{}
	clientArgs := ClientArguments{}

	rootCmd.SetVersionTemplate("{{.Use}} Version {{.Version}}\n")
	rootCmd.PersistentFlags().BoolVarP(
		&rootArgs.Verbose,
		"verbose", "V", false,
		"Show verbose process.",
	)

	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(
		&initArgs.Dir,
		"path", "p", ".",
		"Set the directory where configuration files are stored.",
	)
	initCmd.Flags().IntVar(
		&initArgs.KeySize,
		"keysize", 2048,
		"Set key size.",
	)
	initCmd.Flags().StringVar(
		&initArgs.Algorithm,
		"algorithm", "RSA",
		"Set algorithm for certificate.",
	)
	initCmd.Flags().StringVar(
		&initArgs.Digest,
		"digest", "SHA256",
		"Set digest algorithm for certificate.",
	)
	initCmd.Flags().IntVarP(
		&initArgs.days,
		"days", "d", 180,
		"For how long the certificate remains valid.")

	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(
		&clientArgs.Dir,
		"path", "p", ".",
		"Set the directory where client configuration files are stored.",
	)
	clientCmd.Flags().StringVarP(
		&clientArgs.name,
		"name", "n", "client_cert",
		"Set the client name.",
	)

	var err error = rootCmd.Execute()
	if err != nil {
		return rootArgs, nil, err
	}

	if !slices.Contains(AlgorithmChoices, initArgs.Algorithm) {
		return rootArgs, nil, wrongValue{
			message: fmt.Sprintf(
				"--algorithm expect one of %v; got %v",
				AlgorithmChoices, initArgs.Algorithm),
		}
	}

	switch activated_subcommand {
	case "init":
		if !slices.Contains(AlgorithmChoices, initArgs.Algorithm) {
			return rootArgs, nil, wrongValue{
				message: fmt.Sprintf(
					"--algorithm expect one of %v; got %v",
					AlgorithmChoices, initArgs.Algorithm),
			}
		}
		return rootArgs, initArgs, nil
	case "client":
		return rootArgs, clientArgs, nil
	default:
		return rootArgs, nil, nil
	}

}
