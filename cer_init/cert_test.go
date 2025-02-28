package cerinit

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

var temp_dir, _ = os.MkdirTemp("", "pattern")

func TestEmptyCertDirInit(T *testing.T) {
	// Skip if windows system
	if runtime.GOOS == "windows" {
		T.Skip("Skipping test on Windows system")
	}
	if err := TargetDirInit(temp_dir, false); err != nil {
		fmt.Println(err)
		T.Error(err)
	}
}

func TestDirExistsCertDirInit(T *testing.T) {
	if runtime.GOOS == "windows" {
		T.Skip("Skipping test on Windows system")
	}
	if err := TargetDirInit(temp_dir, false); err != nil {
		fmt.Println(err)
		T.Error(err)
	}
}

func TestFileExistsCertDirInit(T *testing.T) {
	if runtime.GOOS == "windows" {
		T.Skip("Skipping test on Windows system")
	}
	if err := os.RemoveAll(temp_dir); err != nil {
		fmt.Println(err)
		T.Error(err)
	}
	{
		err := os.WriteFile(temp_dir, []byte("This file is a place holder for test.\n"), 0644)
		if err != nil {
			fmt.Println(err)
			T.Error(err)
		}
	}
	if err := TargetDirInit(temp_dir, false); err != nil {
		fmt.Println(err)
		T.Error(err)
	}

}

func TestRemoveTestDir(T *testing.T) {
	if _, err := os.Stat(temp_dir); os.IsNotExist(err) {
		T.Skip("No test directory detected. Skip.")
	}
	err := os.RemoveAll(temp_dir)
	if err != nil {
		T.Error(err)
	}
	T.Logf("Removed test directory: %s", temp_dir)
}
