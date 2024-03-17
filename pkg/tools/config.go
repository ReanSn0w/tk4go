package tools

import (
	"log"

	"github.com/umputun/go-flags"
)

func ParseConfig(opts any) error {
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag|flags.IgnoreUnknown)
	p.SubcommandsOptional = true

	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		return err
	}

	return nil
}
