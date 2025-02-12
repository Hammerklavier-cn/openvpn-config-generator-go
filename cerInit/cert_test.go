package cerinit

import (
	"fmt"
	"os"
	"testing"
)

func TestCertDirInit(T *testing.T) {
	err := targetDirInit("test", false)
	if err != nil {
		fmt.Println(err)
		T.Error(err)
	}
}

func TestRemoveTestDir(T *testing.T) {
	err := os.Remove("test")
	if err != nil {
		T.Log(err)
	}
	err = os.RemoveAll("test")
	if err != nil {
		T.Error(err)
	}
}
