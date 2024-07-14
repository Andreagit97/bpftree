package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"bpftree": Main,
	}))
}

func TestInfoCmd(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:         "testdata/script/cmd_info",
		WorkdirRoot: "testdata/output",
	})
}

func TestLineageCmd(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:         "testdata/script/cmd_lineage",
		WorkdirRoot: "testdata/output",
	})
}

func TestTreeCmd(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:         "testdata/script/cmd_tree",
		WorkdirRoot: "testdata/output",
	})
}

func TestFieldsCmd(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:         "testdata/script/cmd_fields",
		WorkdirRoot: "testdata/output",
	})
}

func TestDumpCmd(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:         "testdata/script/cmd_dump",
		WorkdirRoot: "testdata/output",
	})
}

func TestFields(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:         "testdata/script/supported_fields",
		WorkdirRoot: "testdata/output",
	})
}
