package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

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
	log.Logf("[INFO] Application: %v (rev: %v)", title, revision)

	buf := new(bytes.Buffer)
	structPrinter(buf, 0, "", opts)
	log.Logf("[DEBUG] \n%s", buf.String())
}

func structPrinter(b io.Writer, lvl int, name string, v any) {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
		structPrinter(b, lvl, "", val.Interface())
		return
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			val := val.Type().Field(i).Name

			if field.Kind() == reflect.Ptr {
				if field.IsZero() {
					printSingleValue(b, lvl+1, val, "nil")
					continue
				}

				field = field.Elem()
			}

			if field.Kind() == reflect.Struct {
				printSingleValue(b, lvl+1, val, "")
			}

			structPrinter(b, lvl+1, val, field.Interface())
		}
	default:
		printSingleValue(b, lvl, name, val.Interface())
	}
}

func printSingleValue(b io.Writer, tab int, name string, val any) {
	for i := 0; i < tab; i++ {
		b.Write([]byte("  "))
	}

	switch strings.ToLower(name) {
	case "password", "pass", "token", "secret":
		val = "********"
	}

	b.Write([]byte(fmt.Sprintf("%s:  %v\n", name, val)))
}
