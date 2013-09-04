package fileutil

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func RecursiveDirectoryList(basedir string) ([]string, error) {
	dirs := []string{}
	return dirs, filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		dirs = append(dirs, path)
		return nil
	})
}

// CopyDirectory copies all files and folders underneath src/ and
// copies them under the dst/ directory, recursively.
//
// TODO: The error handling in this function is a bit wonky. We log
// errors and continue on others. I need to do a couple of tests to
// see if legitimate cases need to `continue' after hitting one of
// these errors.
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

	for _, origfile := range subfiles {
		// TODO: It probably doesn't need the EvalSymlinks
		// here since we check to see if the filepath is a
		// symlink after and do something different if it is.
		source_file, err := filepath.EvalSymlinks(filepath.Join(src, origfile.Name()))
		if err != nil {
			return err
		}
		file, err := os.Stat(source_file)
		if err != nil {
			log.Println(err)
			return err
		}
		destination_file := filepath.Join(dst, file.Name())
		// If the original file is a symlink then we can
		// simply create a new symlink with the appropriate
		// link name.
		if issym, err := IsSymLink(filepath.Join(src, origfile.Name())); issym && err == nil {
			linkname, err := os.Readlink(filepath.Join(src, origfile.Name()))
			if err != nil {
				return err
			}
			err = os.Symlink(linkname, filepath.Join(dst, origfile.Name()))
			if err != nil {
				log.Println(err)
			}
			continue
		}
		// If it's a directory we're looking at, recurse to
		// copy its contents.
		if file.IsDir() {
			err = CopyDirectory(destination_file, source_file)
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			// Otherwise we just copy the original file.
			err = CopyFile(destination_file, source_file)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
	return nil
}

// CopyFile takes two pathnames and copies src into dst.
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

// Detects whether the path is a symbolic link or not.
func IsSymLink(fpath string) (bool, error) {
	fi, err := os.Lstat(fpath)
	if err != nil {
		return false, err
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

// MakeAbs detects if the path is an absolute path and if it is not an
// absolute path then it joins together the currect working directory
// path along with the original relative pathname.
func MakeAbs(fpath string) string {
	if !filepath.IsAbs(fpath) {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, fpath)
	}
	return fpath
}
