package gopkg2files

import (
	"flag"
	"io"
	"log"
	"os"
)

type Option struct {
	Debug            bool
	Goroot           bool
	Cgo              bool
	Test             bool
	XTest            bool
	WorkingDirectory string
	Packages         []string
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) FlagSet(name string) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)

	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	flags.BoolVar(&o.Debug, "debug", false, "debug mode")
	flags.BoolVar(&o.Goroot, "goroot", false, "list undered $GOROOT files")
	flags.BoolVar(&o.Cgo, "cgo", false, "list cgo files")
	flags.BoolVar(&o.Test, "test", false, "list test files")
	flags.BoolVar(&o.XTest, "xtest", false, "list xtest files")
	flags.StringVar(&o.WorkingDirectory, "w", workingDirectory, "Working directory")

	return flags, nil
}

func (o *Option) Parse(name string, arguments []string) error {
	fs, err := o.FlagSet(name)
	if err != nil {
		return err
	}

	if err := fs.Parse(arguments); err != nil {
		return err
	}

	o.Packages = fs.Args()
	return nil
}

func (o *Option) Logger() *log.Logger {
	if o.Debug {
		return log.New(os.Stderr, "", log.LstdFlags)
	} else {
		return log.New(io.Discard, "", 0)
	}
}
