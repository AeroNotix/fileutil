package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyDirs(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error(err.Error())
	}
	temp_dir, err := ioutil.TempDir(dir, "testing")
	orig_dir := temp_dir
	for x := 0; x < 3; x++ {
		temp_dir, err = ioutil.TempDir(temp_dir, "testing_dest")
		if err != nil {
			t.Error(err.Error())
		}
		temp_file, err := ioutil.TempFile(temp_dir, "testing")
		if err != nil {
			t.Error(err.Error())
		}
		temp_file.Close()
	}
	}
	DeleteDirectory("destination_test")
	DeleteDirectory(orig_dir)
}
