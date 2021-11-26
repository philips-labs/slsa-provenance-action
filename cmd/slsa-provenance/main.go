package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func main() {
	for i, arg := range os.Args {
		if (strings.HasPrefix(arg, "-") && len(arg) == 2) || (strings.HasPrefix(arg, "--") && len(arg) >= 4) {
			continue
		} else if strings.HasPrefix(arg, "-") {
			newArg := fmt.Sprintf("-%s", arg)
			fmt.Fprintf(os.Stderr, "WARNING: the flag %s is deprecated and will be removed in a future release. Please use the flag %s.\n", arg, newArg)
			os.Args[i] = newArg
		}
	}

	if err := cli.New().ExecuteContext(context.Background()); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
