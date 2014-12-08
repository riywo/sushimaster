package main

import (
	"flag"
	"fmt"
	_ "github.com/mitchellh/go-homedir"
	"github.com/riywo/go-bindata"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	workDir, err := ioutil.TempDir("", "sushimaster_")
	if err != nil {
		errorExit("%+v", err)
	}
	defer func() {
		os.RemoveAll(workDir)
	}()

	input, output := parseArgs()

	err = writeAssets(workDir)
	if err != nil {
		errorExit("writeAssets failed by %+v", err)
	}
	err = writeBindata(workDir, input)
	if err != nil {
		errorExit("writeBindata failed by %+v", err)
	}
	err = buildSushibox(workDir, output)
	if err != nil {
		errorExit("buildSushibox failed by %+v", err)
	}
}

func parseArgs() (input, output string) {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <input directory>\n\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Missing <input directory>\n\n")
		flag.Usage()
		os.Exit(1)
	}

	input = flag.Args()[0]
	pwd, _ := os.Getwd()
	output = filepath.Join(pwd, "sushibox")
	return
}

func makeVersion() string {
	t := time.Now()
	return t.Format("20060102-150405")
}

func writeAssets(workDir string) error {
	sushibox_go, err := Asset("sushibox.go")
	if err != nil {
		return err
	}

	version_go, err := Asset("version.go")
	if err != nil {
		return err
	}

	version := makeVersion()
	version_go = []byte(strings.Replace(string(version_go), "developing", version, 1))

	err = ioutil.WriteFile(filepath.Join(workDir, "sushibox.go"), sushibox_go, os.FileMode(0644))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(workDir, "version.go"), version_go, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}

func writeBindata(workDir, input string) error {
	if _, err := os.Stat(filepath.Join(input, "bin")); err != nil {
		return fmt.Errorf("bin directory does not exist under %s", input)
	}

	cfg := bindata.NewConfig()
	cfg.Input = []bindata.InputConfig{
		{Path: input, Recursive: true},
	}
	cfg.Prefix = input
	cfg.Output = filepath.Join(workDir, "bindata.go")
	return bindata.Translate(cfg)
}

func buildSushibox(workDir, output string) error {
	cmd := exec.Command("go", "build", "-o", output)
	cmd.Dir = workDir
	return cmd.Run()
}

func errorExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", a...)
	os.Exit(1)
}
