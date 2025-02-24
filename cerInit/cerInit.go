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

func targetDirInit(dir string, verbose bool) error {
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

		// drop CA files, etc
		// 1. Drop easy-rsa directory
		if _, err := os.Stat(path.Join(dir, "easy-rsa")); err == nil {
			err := os.RemoveAll(path.Join(dir, "easy-rsa"))
			if err != nil {
				return fmt.Errorf("Failed to remove easy-rsa directory: %w", err)
			}
			if verbose {
				fmt.Printf("Removed easy-rsa directory in %s\n", dir)
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
	err := os.Mkdir(path.Join(dir, "easy-rsa"), 0755)
	if err != nil {
		return fmt.Errorf("Failed to create easy-rsa directory: %w", err)
	}
	if verbose {
		fmt.Printf("Created easy-rsa directory in %s\n", dir)
	}

	// Copy files from /usr/share/easy-rsa to the target directory
	utils.CopyDir("/usr/share/easy-rsa", path.Join(dir, "easy-rsa"))

	return nil
}

func CAGeneration(dir string, algorithm string, verbose bool) error {
	if verbose {
		fmt.Println("Generating CA...")
	}

	// Initialise target dir in case it exists.
	err := targetDirInit(dir, verbose)
	if err != nil {
		return err
	}

	os.Mkdir("easy-rsa", 0755)
	os.Symlink("/usr/share/easy-rsa/*", "./easy-rsa/")

	return nil
}
