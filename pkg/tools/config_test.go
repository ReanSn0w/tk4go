package tools_test

import (
	"os"
	"testing"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

var opts = struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Help    bool `short:"h" long:"help" description:"Show this help message"`
}{}

func Test_ParseConfig(t *testing.T) {
	os.Setenv("VERBOSE", "true")
	os.Setenv("HELP", "true")

	err := tools.ParseConfig(&opts)
	if err != nil {
		t.Error("error parsing config: ", err)
	}
}

func Test_ParseConfigError(t *testing.T) {
	os.Setenv("VERBOSE", "")
	os.Setenv("HELP", "")

	err := tools.ParseConfig(&opts)
	if err != nil {
		t.Error("error parsing config: ", err)
	}
}
