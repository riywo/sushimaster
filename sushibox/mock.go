package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const mockDirEnv = "SUSHIBOX_MOCK"

var mockDir = os.Getenv(mockDirEnv)

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

func RestoreAssets(dir, name string) error { // name is ignored
	if mockDir == "" {
		return fmt.Errorf("Please specify src directory by %s", mockDirEnv)
	}

	f := func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsDir() {
			return nil
		}
		path := strings.TrimPrefix(src, mockDir+string(os.PathSeparator))
		dst := filepath.Join(dir, path)
		return _copyFile(src, info, dst)
	}
	return filepath.Walk(mockDir, f)
}

func _copyFile(src string, info os.FileInfo, dst string) error {
	dir, _ := filepath.Split(dst)
	err := os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	err = os.Chmod(dst, info.Mode())
	if err != nil {
		return err
	}

	err = os.Chtimes(dst, info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}

	return nil
}
