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

var HomeDir = homeDir()
var SushiBoxDir string
var VersionsDir string
var BaseDir string
var BinDir string

func main() {
	os.Exit(realMain())
}

func realMain() int {
	if err := initDirs(); err != nil {
		return errorExit("initDirs failed by %+v", err)
	}

	cmd, args, err := parseArgs()
	if *version {
		return 0
	}
	if err != nil {
		return errorExit("parseArgs failed by %+v", err)
	}

	if err := checkFilesInfo(); err != nil {
		if err = restoreFiles(); err != nil {
			return errorExit("restoreFiles failed by %+v", err)
		}
	}

	_, _, err = execCmd(cmd, args)
	if err != nil {
		return errorExit("execCmd failed by %+v", err)
	}
	return 0
}

func homeDir() string {
	homedir, err := homedir.Dir()
	if err != nil {
		errorExit("homeDir failed by %+v", err)
	}
	return homedir
}

func initDirs() error {
	SushiBoxDir = filepath.Join(HomeDir, ".sushibox")
	VersionsDir = filepath.Join(SushiBoxDir, "versions")

	err := os.MkdirAll(SushiBoxDir, os.FileMode(0755))
	if err != nil {
		return err
	}
	err = os.MkdirAll(VersionsDir, os.FileMode(0755))
	if err != nil {
		return err
	}

	BaseDir = filepath.Join(VersionsDir, Version)
	BinDir = filepath.Join(BaseDir, "bin")
	return nil
}

var version = flag.Bool("version", false, "show version")

func parseArgs() (cmd string, args []string, err error) {
	cmd, args = filepath.Base(os.Args[0]), os.Args[1:]
	*version = false
	if cmd == "sushibox" {
		flag.Usage = func() {
			fmt.Printf("Usage: %s [options] command args...\n\n", cmd)
			flag.PrintDefaults()
		}

		flag.Parse()

		if *version {
			fmt.Printf("%s\n", Version)
			return
		}

		if len(args) == 0 {
			flag.Usage()
			err = fmt.Errorf("missing args")
		} else {
			cmd, args = flag.Args()[0], flag.Args()[1:]
		}
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

var execFunc = func(arg0 string, argv, envv []string) (stdout, stderr []byte, err error) {
	syscall.Exec(arg0, argv, envv)
	return // never called
}

func execCmd(cmd string, args []string) (stdout, stderr []byte, err error) {
	cmdPath := filepath.Join(BinDir, cmd)
	argv := append([]string{cmdPath}, args...)
	envv := os.Environ()
	return execFunc(cmdPath, argv, envv)
}

func errorExit(format string, a ...interface{}) int {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", a...)
	return 1
}
