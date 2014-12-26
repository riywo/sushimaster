package main

import (
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type SushiboxTestSuite struct {
	suite.Suite
	tempDir string
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SushiboxTestSuite))
}

func (suite *SushiboxTestSuite) SetupTest() {
	mockDir = "test"
	execFunc = execMockFunc
	suite.tempDir, _ = ioutil.TempDir("", "sushibox_test_")
	HomeDir = suite.tempDir
	err := initDirs()
	suite.Nil(err)
}

func (suite *SushiboxTestSuite) TearDownTest() {
	err := os.RemoveAll(suite.tempDir)
	suite.Nil(err)
}

func (suite *SushiboxTestSuite) TestInitSushiBoxDir() {
	suite.Equal(filepath.Join(HomeDir, ".sushibox"), SushiBoxDir)
	info, _ := os.Stat(SushiBoxDir)
	suite.True(info.IsDir())
}
func (suite *SushiboxTestSuite) TestInitVersionsDir() {
	suite.Equal(filepath.Join(SushiBoxDir, "versions"), VersionsDir)
	info, _ := os.Stat(VersionsDir)
	suite.True(info.IsDir())
}

func (suite *SushiboxTestSuite) TestInitBaseDir() {
	suite.Equal(filepath.Join(VersionsDir, "developing"), BaseDir)
	_, err := os.Stat(BaseDir)
	suite.True(os.IsNotExist(err))
}

func (suite *SushiboxTestSuite) TestInitBinDir() {
	suite.Equal(filepath.Join(BaseDir, "bin"), BinDir)
	_, err := os.Stat(BinDir)
	suite.True(os.IsNotExist(err))
}

func (suite *SushiboxTestSuite) TestParseArgs() {
	os.Args = []string{"sushibox", "foo", "bar", "baz"}
	cmd, args, err := parseArgs()
	suite.Equal("foo", cmd)
	suite.Equal([]string{"bar", "baz"}, args)
	suite.Nil(err)
	suite.False(*version)
}

func (suite *SushiboxTestSuite) TestParseArgsAsOtherCommand() {
	os.Args = []string{"a", "b", "c"}
	cmd, args, err := parseArgs()
	suite.Equal("a", cmd)
	suite.Equal([]string{"b", "c"}, args)
	suite.Nil(err)
	suite.False(*version)
}

func (suite *SushiboxTestSuite) TestParseArgsMissing() {
	os.Args = []string{"sushibox"}
	_, _, err := parseArgs()
	suite.NotNil(err)
	suite.False(*version)
}

func (suite *SushiboxTestSuite) TestParseArgsVersion() {
	os.Args = []string{"sushibox", "-version"}
	_, _, err := parseArgs()
	suite.Nil(err)
	suite.True(*version)
}

func (suite *SushiboxTestSuite) TestParseArgsVersionAsOtherCommand() {
	os.Args = []string{"foo", "-version"}
	cmd, args, err := parseArgs()
	suite.Equal("foo", cmd)
	suite.Equal([]string{"-version"}, args)
	suite.Nil(err)
	suite.False(*version)
}

func (suite *SushiboxTestSuite) TestCheckFilesInfo() {
	suite.Nil(RestoreAssets(BaseDir, ""))
	suite.Nil(checkFilesInfo())
}

func (suite *SushiboxTestSuite) TestCheckFilesInfoNoFile() {
	suite.True(os.IsNotExist(checkFilesInfo()))
}

func (suite *SushiboxTestSuite) TestCheckFilesInfoModified() {
	suite.Nil(RestoreAssets(BaseDir, ""))
	os.Truncate(filepath.Join(BaseDir, AssetNames()[0]), 1)
	suite.NotNil(checkFilesInfo())
}

func (suite *SushiboxTestSuite) TestRestoreFiles() {
	suite.Nil(restoreFiles())
	fi1, _ := os.Stat(filepath.Join(BinDir, "foo"))
	suite.Nil(restoreFiles())
	fi2, _ := os.Stat(filepath.Join(BinDir, "foo"))
	suite.True(os.SameFile(fi1, fi2))
}

func (suite *SushiboxTestSuite) TestExecCmd() {
	suite.Nil(restoreFiles())
	stdout, stderr, err := execCmd("foo", []string{})
	suite.Equal("foo\n", string(stdout))
	suite.Equal("", string(stderr))
	suite.Nil(err)
}

func (suite *SushiboxTestSuite) TestExecCmdNotExecutable() {
	suite.Nil(restoreFiles())
	stdout, stderr, err := execCmd("bar", []string{})
	suite.Equal("", string(stdout))
	suite.Equal("", string(stderr))
	suite.True(os.IsPermission(err))
}

func (suite *SushiboxTestSuite) TestRealMain() {
	os.Args = []string{"sushibox", "foo"}
	suite.Equal(0, realMain())
}

func (suite *SushiboxTestSuite) TestRealMainAsOtherCommand() {
	os.Args = []string{"foo"}
	suite.Equal(0, realMain())
}

func (suite *SushiboxTestSuite) TestRealMainNotExecutable() {
	os.Args = []string{"sushibox", "bar"}
	suite.Equal(1, realMain())
}

func (suite *SushiboxTestSuite) TestRealMainNotFound() {
	os.Args = []string{"sushibox", "baz"}
	suite.Equal(1, realMain())
}
