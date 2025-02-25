package cerinit

import (
	"fmt"
	"os"
	"testing"
)

func TestEmptyCertDirInit(T *testing.T) {
	if err := TargetDirInit("test", false); err != nil {
		fmt.Println(err)
		T.Error(err)
	}
}

func TestDirExistsCertDirInit(T *testing.T) {
	if err := TargetDirInit("test", false); err != nil {
		fmt.Println(err)
		T.Error(err)
	}
}

func TestFileExistsCertDirInit(T *testing.T) {
	if err := os.RemoveAll("test"); err != nil {
		fmt.Println(err)
		T.Error(err)
	}
	{
		err := os.WriteFile("test", []byte("This file is a place holder for test.\n"), 0644)
		if err != nil {
			fmt.Println(err)
			T.Error(err)
		}
	}
	if err := TargetDirInit("test", false); err != nil {
		fmt.Println(err)
		T.Error(err)
	}

}

func TestRemoveTestDir(T *testing.T) {
	err := os.RemoveAll("test")
	if err != nil {
		T.Error(err)
	}
}
