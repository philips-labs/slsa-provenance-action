package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func main() {
	rootFlagSet := flag.NewFlagSet("slsa-provenance", flag.ExitOnError)

	app := &ffcli.Command{
		Name:    "slsa-provenance [flags] <subcommand>",
		FlagSet: rootFlagSet,
		Subcommands: []*ffcli.Command{
			cli.Generate(),
		},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println("slsa-provenance")
			fmt.Println()
		},
	}

	if err := app.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
