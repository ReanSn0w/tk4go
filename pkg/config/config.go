package config

import (
	"bytes"
	"io"
	"log"
	"reflect"

	"github.com/ReanSn0w/tk4go/pkg/tools"
	"github.com/umputun/go-flags"
)

func Parse(opts any) error {
	p := flags.NewParser(opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag|flags.IgnoreUnknown)
	p.SubcommandsOptional = true

	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		return err
	}

	return nil
}

func Print(log tools.Logger, title string, revision string, opts any) {
	buf := new(bytes.Buffer)
	structPrinter(buf, 0, opts)
	log.Logf("\nApplication: %v (rev: %v) \n%s", title, revision, buf.String())
}

func structPrinter(b io.Writer, lvl int, v any) {
	val := reflect.ValueOf(v)
	if val.IsNil() {
		return
	}

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
		structPrinter(b, lvl, val.Interface())
		return
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			structPrinter(b, lvl+1, field.Interface())
		}
	default:
		// print value
		name := val.Type().Name()
		val = reflect.Indirect(val)

		for i := 0; i < lvl; i++ {
			b.Write([]byte("  "))
		}

		b.Write([]byte(name + ": " + val.String() + "\n"))
	}
}
