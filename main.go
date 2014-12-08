package main

import (
	_ "github.com/mitchellh/go-homedir"
	"github.com/riywo/go-bindata"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func makeVersion() string {
	t := time.Now()
	return t.Format("20060102-150405")
}

func main() {
	input := os.Args[1]

	sushibox_go, err := Asset("sushibox.go")
	if err != nil {
		panic(err)
	}
	version_go, err := Asset("version.go")
	if err != nil {
		panic(err)
	}

	version := makeVersion()
	version_go = []byte(strings.Replace(string(version_go), "developing", version, 1))

	workDir, err := ioutil.TempDir("", "sushimaster_")
	if err != nil {
		panic(err)
	}
	defer func() {
		os.RemoveAll(workDir)
	}()

	err = ioutil.WriteFile(filepath.Join(workDir, "sushibox.go"), sushibox_go, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(workDir, "version.go"), version_go, os.FileMode(0644))
	if err != nil {
		panic(err)
	}

	cfg := bindata.NewConfig()
	cfg.Input = []bindata.InputConfig{
		{Path: input, Recursive: true},
	}
	cfg.Prefix = input
	cfg.Output = filepath.Join(workDir, "bindata.go")
	err = bindata.Translate(cfg)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-o", "sushibox")
	cmd.Dir = workDir
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = os.Rename(filepath.Join(workDir, "sushibox"), "sushibox")
	if err != nil {
		panic(err)
	}
}
