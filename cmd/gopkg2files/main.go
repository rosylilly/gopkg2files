package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/rosylilly/gopkg2files"
)

func main() {
	option := gopkg2files.NewOption()
	if err := option.Parse("gopkg2files", os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		os.Exit(1)
		return
	}

	resolver := gopkg2files.NewResolver(option)
	if err := resolver.ResolveAll(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
		return
	}

	if err := resolver.PrintFiles(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
		return
	}
}
