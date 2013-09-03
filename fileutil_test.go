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
	dirpaths := make([]string, 11)
	temp_dir, err := ioutil.TempDir(dir, "testing")
	orig_dir := temp_dir
	dirpaths = append(dirpaths, temp_dir)

	filepaths := make([]string, 10)
	for x := 0; x < 10; x++ {
		temp_dir, err = ioutil.TempDir(temp_dir, "testing_dest")
		if err != nil {
			t.Error(err.Error())
		}
		temp_file, err := ioutil.TempFile(temp_dir, "testing")
		if err != nil {
			t.Error(err.Error())
		}
		dirpaths = append(dirpaths, temp_dir)
		filepaths = append(filepaths, temp_file.Name())
	}
	CopyDirectory(orig_dir, "destination_test")
	for _, path := range filepaths {
		err = os.Remove(path)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
