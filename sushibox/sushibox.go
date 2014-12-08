package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

var SushiBoxDir = ".sushibox"
var VersionsDir = "versions"

var BaseDir = Version
var BinDir = "bin"

func main() {
	if err := initDirs(); err != nil {
		errorExit("initDirs failed by %+v", err)
	}

	cmd, args := parseArgs()

	if err := checkFilesInfo(); err != nil {
		if err = restoreFiles(); err != nil {
			errorExit("restoreFiles failed by %+v", err)
		}
	}

	if err := execCmd(cmd, args); err != nil {
		errorExit("execCmd %s failed by %+v", cmd, err)
	}
}

func initDirs() error {
	homedir, err := homedir.Dir()
	if err != nil {
		return err
	}
	SushiBoxDir = filepath.Join(homedir, SushiBoxDir)
	err = os.MkdirAll(SushiBoxDir, os.FileMode(0755))
	if err != nil {
		return err
	}
	VersionsDir = filepath.Join(SushiBoxDir, VersionsDir)
	err = os.MkdirAll(VersionsDir, os.FileMode(0755))
	if err != nil {
		return err
	}

	BaseDir = filepath.Join(VersionsDir, BaseDir)
	BinDir = filepath.Join(BaseDir, BinDir)
	return nil
}

func parseArgs() (cmd string, args []string) {
	cmd = filepath.Base(os.Args[0])

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] command args...\n\n", cmd)
		flag.PrintDefaults()
	}
	var version = flag.Bool("version", false, "show version")

	flag.Parse()
	args = flag.Args()

	if *version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	if cmd == "sushibox" {
		if len(args) == 0 {
			flag.Usage()
			os.Exit(1)
		}
		cmd, args = args[0], args[1:]
	}

	return
}

func checkFilesInfo() error {
	for _, name := range AssetNames() {
		assetinfo, err := AssetInfo(name)
		if err != nil {
			return err
		}

		path := filepath.Join(BaseDir, name)
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		if assetinfo.Size() != fileInfo.Size() {
			return fmt.Errorf("check info error %s: size is different", path)
		}
		if assetinfo.Mode() != fileInfo.Mode() {
			return fmt.Errorf("check info error %s: mode is different", path)
		}
		if assetinfo.ModTime() != fileInfo.ModTime() {
			return fmt.Errorf("check info error %s: mtime is different", path)
		}
	}
	return nil
}

func restoreFiles() error {
	tempDir, err := ioutil.TempDir("", fmt.Sprintf("sushibox_%s_", Version))
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	err = RestoreAssets(tempDir, "")
	if err != nil {
		return err
	}

	err = os.Rename(tempDir, BaseDir)
	if os.IsNotExist(err) {
		return err
	}

	return checkFilesInfo()
}

func execCmd(cmd string, args []string) error {
	cmdPath := filepath.Join(BinDir, cmd)
	argv := append([]string{cmdPath}, args...)
	envv := os.Environ()
	return syscall.Exec(cmdPath, argv, envv)
}

func errorExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", a...)
	os.Exit(1)
}
