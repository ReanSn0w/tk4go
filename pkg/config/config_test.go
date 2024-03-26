package config_test

import (
	"os"
	"testing"

	"github.com/ReanSn0w/tk4go/pkg/config"
)

var opts = struct {
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Help    bool `short:"h" long:"help" description:"Show this help message"`

	Struct *struct {
		Password string `short:"t" long:"test" description:"Test string"`
	}
}{}

func Test_ParseConfig(t *testing.T) {
	os.Setenv("VERBOSE", "true")
	os.Setenv("HELP", "true")

	err := config.Parse(&opts)
	if err != nil {
		t.Error("error parsing config: ", err)
	}
}

func Test_ParseConfigError(t *testing.T) {
	os.Setenv("VERBOSE", "")
	os.Setenv("HELP", "")

	err := config.Parse(&opts)
	if err != nil {
		t.Error("error parsing config: ", err)
	}
}

// func Test_Print(t *testing.T) {
// 	opts.Verbose = true
// 	opts.Struct = &struct {
// 		Password string `short:"t" long:"test" description:"Test string"`
// 	}{
// 		Password: "test",
// 	}

// 	config.Print(t, "test", "1.0.0", opts)
// 	t.Error("error printing config: ", t)
// }
