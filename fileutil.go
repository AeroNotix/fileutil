package fileutil

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyDirectory(src, dst string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	_, err = os.Open(dst)
	if !os.IsNotExist(err) {
		return errors.New("Destination directory already exists.")
	}

	err = os.MkdirAll(dst, fi.Mode())
	if err != nil {
		return err
	}

	subfiles, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range subfiles {
		source_file := filepath.Join(src, file.Name())
		destination_file := filepath.Join(dst, file.Name())
		if file.IsDir() {
			err = CopyDirectory(source_file, destination_file)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(source_file, destination_file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(src, dst string) error {
	source_file, err := os.Open(src)
	if err != nil {
		return err
	}
	source_stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	destination_file, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, source_stat.Mode())
	if err != nil {
		return err
	}
	_, err = io.Copy(source_file, destination_file)
	return err
}
