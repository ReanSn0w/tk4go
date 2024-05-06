package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/ReanSn0w/tk4go/pkg/config"
	"github.com/ReanSn0w/tk4go/pkg/tools"
	"github.com/ReanSn0w/tk4go/pkg/web"
	"github.com/go-pkgz/lgr"
)

var (
	revision = "unknown"
	log      = lgr.Default()

	opts = struct {
		Debug bool   `long:"debug" description:"debug logs level"`
		Port  int    `short:"p" long:"port" default:"8080" description:"server port"`
		Dir   string `short:"d" long:"dir" default:"." description:"serve directory"`
	}{}
)

func main() {
	err := config.Parse(&opts)
	if err != nil {
		log.Logf("%v", err)
		os.Exit(2)
	}

	if opts.Debug {
		lgr.Setup(lgr.Debug, lgr.CallerFunc)
	}

	config.Print(log, "Simple http server", revision, opts)

	ctx, cancel := context.WithCancelCause(context.TODO())
	defer cancel(nil)

	gs := tools.NewShutdownStack(log)

	srv := web.NewServer(log)
	srv.Run(cancel, opts.Port, http.FileServer(http.Dir(opts.Dir)))
	gs.Add(func(ctx context.Context) {
		srv.Shutdown(ctx)
	})

	go tools.AnyKeyToExit(log, func() {
		cancel(nil)
	})

	gs.Wait(ctx, time.Second*10)
}
