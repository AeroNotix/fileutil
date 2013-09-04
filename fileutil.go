package fileutil

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RecursiveDirectoryList(basedir string) ([]string, error) {
	dirs := []string{}
	return dirs, filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		dirs = append(dirs, MakeAbs(path))
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

func DeleteDirectory(base string) error {
	destination_list, err := RecursiveDirectoryList(base)
	if err != nil {
		return err
	}
	for x := len(destination_list) - 1; x != -1; x-- {
		err = os.Remove(destination_list[x])
		if err != nil {
			return err
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

func makeSetFromStrSlice(strslice []string) map[string]struct{} {
	m := make(map[string]struct{}, len(strslice))
	for _, str := range strslice {
		m[str] = struct{}{}
	}
	return m
}

func RootPath(path string) string {
	paths := strings.Split(path, string(filepath.Separator))
	return paths[0]
}

func DiffDirectories(dir_a, dir_b string) (plus, minus []string, err error) {
	list_a, err := RecursiveDirectoryList(dir_a)
	if err != nil {
		return nil, nil, err
	}
	list_b, err := RecursiveDirectoryList(dir_b)
	if err != nil {
		return nil, nil, err
	}

	// the first element is the pathname itself, so we can ignore
	// that.
	set_a := makeSetFromStrSlice(list_a[1:])
	set_b := makeSetFromStrSlice(list_b[1:])

	// We need to get the absolute path of the directory so we can
	// swap the other directory's leading path section with ours
	// so the directory listing will match properly.
	full_path_a := MakeAbs(dir_a)
	full_path_b := MakeAbs(dir_b)

	for path, _ := range set_a {
		rel_path := filepath.Join(full_path_b, path[len(full_path_a):])
		if _, ok := set_b[rel_path]; !ok {
			plus = append(plus, rel_path)
		}
	}

	for path, _ := range set_b {
		rel_path := filepath.Join(full_path_a, path[len(full_path_b):])
		if _, ok := set_a[rel_path]; !ok {
			minus = append(minus, rel_path)
		}
	}

	return
}
