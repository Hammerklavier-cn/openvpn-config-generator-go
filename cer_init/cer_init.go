package cerinit

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"gitcode.com/hammerklavier/openvpn-config-generator-go/utils"
)

// Error because of user abort
type UserAbort struct {
	message string
}

func (u UserAbort) Error() string { return u.message }

func TargetDirInit(dir string, verbose bool) error {
	if file_info, err := os.Stat(dir); os.IsNotExist(err) {
		// In this case, target dir does not exists.
		// A new directory named `dir` will be created.
		if verbose {
			fmt.Printf("Assigned dir %s does not exist. Create one.\n", dir)
		}
		os.Mkdir(dir, 0755)
	} else if file_info.IsDir() == true {
		// In this case, a directory with the same name exists.
		// After confirmation, some contents of `dir` will be removed.
		fmt.Printf("Target dir `%s` already exist. All changes made in this directory will be purged.\n", dir)

		fmt.Printf("Sure to proceed? [Y/n]\t")
		var doProceed bool
		{
			var inputString string
			fmt.Scanln(&inputString)
			if inputString == "" {
				doProceed = true
			} else if slices.Contains([]string{"yes", "y"}, strings.ToLower(inputString)) {
				doProceed = true
			} else {
				doProceed = false
			}
		}
		if doProceed == false {
			return UserAbort{message: "Abort"}
		} else {
			fmt.Println("Proceed confirmed.")
		}

		// 1. Drop easy-rsa related files and directories
		{
			var files_and_folders = []string{
				"x509-types", "easyrsa", "openssl-easyrsa.cnf", "vars.example", "pki",
			}

			for _, file_or_folder := range files_and_folders {
				if _, err := os.Stat(path.Join(dir, file_or_folder)); err == nil {
					err := os.RemoveAll(path.Join(dir, file_or_folder))
					if err != nil {
						return fmt.Errorf("Failed to remove %s: %w", file_or_folder, err)
					}
					if verbose {
						fmt.Printf("Removed %s in %s\n", file_or_folder, dir)
					}
				}
			}
		}
		// 2. _PLACE Holder_

	} else if file_info.IsDir() == false {
		// In this case, a file with the same name exists.
		// After confirmation, the file will be removed and replaced with a directory named `dir`.
		fmt.Printf("Target path `%s` already exists and is not a directory. This folder will be purged.\n", dir)

		fmt.Printf("Sure to proceed? [Y/n]\t")
		var doProceed bool
		{
			var inputString string
			fmt.Scanln(&inputString)
			if inputString == "" {
				doProceed = true
			} else if slices.Contains([]string{"yes", "y"}, strings.ToLower(inputString)) {
				doProceed = true
			} else {
				doProceed = false
			}
		}
		if doProceed == false {
			return UserAbort{message: "Abort"}
		} else {
			fmt.Println("Proceed confirmed.")
		}
		// remove the file and replace it with a folder
		err := os.Remove(dir)
		if err != nil {
			return fmt.Errorf("Failed to remove existing file: %w", err)
		}
		if verbose {
			fmt.Printf("Removed existing file at %s\n", dir)
		}
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return fmt.Errorf("Failed to create directory: %w", err)
		}
		if verbose {
			fmt.Printf("Created new directory at %s\n", dir)
		}
	}

	// Copy /usr/share/easy-rsa to the target dir.
	// err := os.Mkdir(path.Join(dir, "easy-rsa"), 0755)
	// if err != nil {
	// 	return fmt.Errorf("Failed to create easy-rsa directory: %w", err)
	// }
	if verbose {
		fmt.Printf("Created easy-rsa directory in %s\n", dir)
	}

	// Copy files from /usr/share/easy-rsa to the target directory
	if err := utils.CopyDir("/usr/share/easy-rsa", path.Join(dir)); err != nil {
		return fmt.Errorf("Failed to copy easy-rsa files: %w", err)
	}

	return nil
}

func CAGeneration(dir string, algorithm string, verbose bool) error {
	if verbose {
		fmt.Println("Generating CA...")
	}

	// Initialise target dir in case it exists.
	err := TargetDirInit(dir, verbose)
	if err != nil {
		return err
	}

	os.Mkdir("easy-rsa", 0755)
	os.Symlink("/usr/share/easy-rsa/*", "./easy-rsa/")

	return nil
}

func CertificateAuthorityInit(dir string, algorithm string, digest string, verbose bool) error {
	if verbose {
		fmt.Printf("Initialising Certificate Authority in %s...\n", dir)
	}

	// Initialise target dir in case it exists.
	if err := TargetDirInit(dir, verbose); err != nil {
		return err
	}

	// create `vars` file
	if err := CreateVarsFile(dir, algorithm, digest); err != nil {
		return err
	}

	// initialise PKI
	if err := initPKI(dir); err != nil {
		return err
	}

	return nil
}

func CreateVarsFile(dir string, algorithm string, digest string) error {
	file, err := os.Create(path.Join(dir, "vars"))
	if err != nil {
		return fmt.Errorf("Failed to create vars file: %w", err)
	}
	defer file.Close()

	// file.WriteString("set_var ")
	file.WriteString("set_var EASYRSA_REQ_COUNTRY    \"CN\"\n")
	file.WriteString("set_var EASYRSA_REQ_PROVINCE   \"Guangdong\"\n")
	file.WriteString("set_var EASYRSA_REQ_CITY       \"Guangzhou\"\n")
	file.WriteString("set_var EASYRSA_REQ_ORG        \"OpenVPN Config Generator\"\n")
	file.WriteString("set_var EASYRSA_REQ_EMAIL      \"admin@ovpngen.com\"\n")
	file.WriteString("set_var EASYRSA_REQ_OU         \"Community\"\n")
	file.WriteString(fmt.Sprintf("set_var EASYRSA_ALGO           \"%s\"\n", strings.ToLower(algorithm)))
	file.WriteString(fmt.Sprintf("set_var EASYRSA_DIGEST         \"%s\"\n", strings.ToLower(digest)))

	return nil
}
