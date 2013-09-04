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
	CopyDirectory("destination_test", orig_dir)
	if plus, minus, err := DiffDirectories("destination_test", orig_dir); false ||
		err != nil || plus != nil || minus != nil {
		fmt.Println(len(plus), len(minus))
		t.Error("Directories do not match!", plus, minus)
	}
	DeleteDirectory("destination_test")
	DeleteDirectory(orig_dir)

}
