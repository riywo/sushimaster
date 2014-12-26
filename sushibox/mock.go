package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const mockDirEnv = "SUSHIBOX_MOCK"

var mockDir = os.Getenv(mockDirEnv)

var execMockFunc = func(arg0 string, argv, envv []string) (stdout, stderr []byte, err error) {
	cmd := exec.Command(arg0, argv[1:]...)
	cmd.Env = envv
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outbuf, &errbuf
	err = cmd.Run()
	stdout, stderr = outbuf.Bytes(), errbuf.Bytes()
	return
}

func Asset(name string) ([]byte, error) {
	if mockDir == "" {
		return nil, fmt.Errorf("Please specify src directory by %s", mockDirEnv)
	}

	path := filepath.Join(mockDir, name)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

func AssetInfo(name string) (os.FileInfo, error) {
	if mockDir == "" {
		return nil, fmt.Errorf("Please specify src directory by %s", mockDirEnv)
	}

	path := filepath.Join(mockDir, name)
	fi, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset info %s at %s: %v", name, path, err)
	}
	return fi, err
}

func AssetNames() (names []string) {
	if mockDir == "" {
		errorExit("Please specify src directory by %s", mockDirEnv)
	}

	f := func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsDir() {
			return nil
		}
		path := strings.TrimPrefix(src, mockDir+string(os.PathSeparator))
		names = append(names, path)
		return nil
	}
	err := filepath.Walk(mockDir, f)
	if err != nil {
		errorExit("AssetNames failed by %+v", err)
	}
	return
}

func AssetDir(name string) (names []string, err error) {
	path := filepath.Join(mockDir, name)
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = fmt.Errorf("file %s", path)
		return
	}
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, fi := range list {
		n := strings.TrimPrefix(fi.Name(), mockDir+string(os.PathSeparator))
		names = append(names, n)
	}
	return
}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	if err != nil { // File
		return RestoreAsset(dir, name)
	} else { // Dir
		for _, child := range children {
			err = RestoreAssets(dir, path.Join(name, child))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
