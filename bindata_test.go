package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestAsset(t *testing.T) {
	for _, name := range []string{"sushibox.go", "version.go"} {
		expected, _ := ioutil.ReadFile(filepath.Join("sushibox", name))
		actual, _ := Asset(name)
		assert.Equal(t, string(expected), string(actual))
	}
}
