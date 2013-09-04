package fileutil

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func CopyDirectory(dst, src string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
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
		source_file, err := filepath.EvalSymlinks(filepath.Join(src, file.Name()))
		if err != nil {
			continue
		}
		destination_file := filepath.Join(dst, file.Name())

		if file.IsDir() {
			err = CopyDirectory(destination_file, source_file)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(destination_file, source_file)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
	return nil
}

func CopyFile(dst, src string) error {
	source_file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source_file.Close()
	source_stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	destination_file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, source_stat.Mode())
	if err != nil {
		return err
	}
	defer destination_file.Close()
	_, err = io.Copy(destination_file, source_file)
	return err
}
