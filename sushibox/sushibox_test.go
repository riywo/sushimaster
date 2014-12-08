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
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SushiboxTestSuite))
}

var tempDir string

func (suite *SushiboxTestSuite) SetupTest() {
	tempDir, _ = ioutil.TempDir("", "sushibox_test_")
	HomeDir = tempDir
	err := initDirs()
	suite.Nil(err)
}

func (suite *SushiboxTestSuite) TearDownTest() {
	err := os.RemoveAll(tempDir)
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
